import express, { Request, Response } from "express";
import { S3Client, ListObjectsV2Command, GetObjectCommand, DeleteObjectCommand } from "@aws-sdk/client-s3";

const app = express();
app.use(express.json());

// ---- ENV ----
const S3_ENDPOINT = process.env.S3_ENDPOINT!;
const S3_REGION = process.env.S3_REGION!;
const S3_ACCESS_KEY_ID = process.env.S3_ACCESS_KEY_ID!;
const S3_SECRET_ACCESS_KEY = process.env.S3_SECRET_ACCESS_KEY!;
const S3_BUCKET = process.env.S3_BUCKET!;

const XML_SERVICE_URL = process.env.XML_SERVICE_URL ?? "http://xml-service:8081/ingest";
const EXTERNAL_API_URL = process.env.EXTERNAL_API_URL ?? "http://external-api:8090/tax";
const WEBHOOK_URL = process.env.WEBHOOK_URL ?? "http://processor:8080/webhook/xml-status";

const POLL_SECONDS = Number(process.env.POLL_SECONDS ?? 15);

// ---- S3 client (Supabase S3) ----
const s3 = new S3Client({
  region: S3_REGION,
  endpoint: S3_ENDPOINT,
  credentials: {
    accessKeyId: S3_ACCESS_KEY_ID,
    secretAccessKey: S3_SECRET_ACCESS_KEY
  },
  forcePathStyle: true // importante para Supabase S3
});

async function streamToBuffer(body: any): Promise<Buffer> {
  // Body no Node SDK v3 é um stream async iterable
  const chunks: Buffer[] = [];
  for await (const chunk of body) {
    chunks.push(Buffer.isBuffer(chunk) ? chunk : Buffer.from(chunk));
  }
  return Buffer.concat(chunks);
}

async function fetchTaxRate(state: string): Promise<number> {
  const url = `${EXTERNAL_API_URL}?state=${encodeURIComponent(state)}`;
  const resp = await fetch(url);
  if (!resp.ok) return 0.05;
  const json: any = await resp.json();
  return typeof json.taxRate === "number" ? json.taxRate : 0.05;
}

function detectStateFromCsv(csvText: string): string {
  // Pega a coluna "State" da primeira linha de dados (simples e suficiente)
  const lines = csvText.split(/\r?\n/).filter(Boolean);
  if (lines.length < 2) return "Kentucky";

  const header = lines[0].split(",");
  const stateIdx = header.findIndex(h => h.trim() === "State");
  if (stateIdx < 0) return "Kentucky";

  const firstData = lines[1].split(",");
  return (firstData[stateIdx] ?? "Kentucky").split('"').join("").trim() || "Kentucky";
}

async function sendToXmlService(csvBytes: Buffer) {
  const csvText = csvBytes.toString("utf-8");
  const state = detectStateFromCsv(csvText);
  const taxRate = await fetchTaxRate(state);

  const idRequisicao = `REQ-${Date.now()}`;
  const mapperVersion = "1.0";

  const form = new FormData();
  form.append("idRequisicao", idRequisicao);
  form.append("mapperVersion", mapperVersion);
  form.append("webhookUrl", WEBHOOK_URL);
  form.append("taxRate", String(taxRate));

  const csvBlob = new Blob([new Uint8Array(csvBytes)], { type: "text/csv" });
  form.append("file", csvBlob, "input.csv");

  const resp = await fetch(XML_SERVICE_URL, {
    method: "POST",
    body: form
  });

  const text = await resp.text();
  return { idRequisicao, status: resp.status, body: text, state, taxRate };
}

async function pollOnce() {
  const list = await s3.send(new ListObjectsV2Command({
    Bucket: S3_BUCKET,
    Prefix: "input/"
  }));

  const objects = (list.Contents ?? [])
    .filter(o => o.Key && o.Key.endsWith(".csv"))
    .sort((a, b) => (a.LastModified?.getTime() ?? 0) - (b.LastModified?.getTime() ?? 0));

  if (objects.length === 0) {
    console.log("[processor] no new CSVs");
    return;
  }

  const key = objects[0].Key!;
  console.log(`[processor] processing: ${key}`);

  const get = await s3.send(new GetObjectCommand({ Bucket: S3_BUCKET, Key: key }));
  const csvBytes = await streamToBuffer(get.Body);

  const result = await sendToXmlService(csvBytes);
  console.log("[processor] xml-service result:", result);

  if (result.status >= 200 && result.status < 300) {
    await s3.send(new DeleteObjectCommand({ Bucket: S3_BUCKET, Key: key }));
    console.log(`[processor] deleted source CSV: ${key}`);
  } else {
    console.log(`[processor] NOT deleting ${key} (xml-service status ${result.status})`);
  }
}

// ---- webhook endpoint ----
app.post("/webhook/xml-status", (req: Request, res: Response) => {
  console.log("Webhook recebido:", req.body);
  res.status(200).json({ ok: true });
});

// Endpoint manual (opcional) para forçar polling
app.post("/poll-now", async (_req: Request, res: Response) => {
  try {
    await pollOnce();
    res.json({ ok: true });
  } catch (e: any) {
    res.status(500).json({ ok: false, error: e?.message ?? String(e) });
  }
});

const port = Number(process.env.PORT ?? 8080);
app.listen(port, () => {
  console.log(`processor up on :${port}`);

  // Loop de polling
  setInterval(() => {
    pollOnce().catch(err => console.error("[processor] poll error:", err));
  }, POLL_SECONDS * 1000);
});
