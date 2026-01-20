import * as grpc from "@grpc/grpc-js";
import * as protoLoader from "@grpc/proto-loader";
import path from "path";

type BIServiceClient = grpc.Client & {
  SalesByCategory(
    req: { category: string },
    cb: (err: grpc.ServiceError | null, res?: { category: string; total: number }) => void
  ): void;
};

const PROTO_PATH = path.join(__dirname, "../../proto/bi_grpc.proto");

// carrega proto
const pkgDef = protoLoader.loadSync(PROTO_PATH, {
  keepCase: true,
  longs: String,
  enums: String,
  defaults: true,
  oneofs: true,
});

const loaded = grpc.loadPackageDefinition(pkgDef) as any;

// package = "xmlservice", service = "BIService"
const BIService = loaded.xmlservice.BIService;

const XML_GRPC_HOST = process.env.XML_GRPC_HOST || "xml-service:50051";

const client: BIServiceClient = new BIService(
  XML_GRPC_HOST,
  grpc.credentials.createInsecure()
);

export function salesByCategoryGrpc(category: string): Promise<number> {
  return new Promise((resolve, reject) => {
    client.SalesByCategory({ category }, (err, res) => {
      if (err) return reject(err);
      resolve(res?.total ?? 0);
    });
  });
}
