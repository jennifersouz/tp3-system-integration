import pkg from "pg";
const { Client } = pkg;

const dbClient = new Client({
  host: process.env.DB_HOST || "db",
  port: parseInt(process.env.DB_PORT || "5432"),
  user: process.env.DB_USER || "tp3",
  password: process.env.DB_PASSWORD || "tp3",
  database: process.env.DB_NAME || "tp3db",
});

let connected = false;

export async function ensureConnected() {
  if (!connected) {
    try {
      await dbClient.connect();
      connected = true;
      console.log("[DB] Conectado ao PostgreSQL");
    } catch (err) {
      console.error("[DB] Erro de conexão:", err);
      throw err;
    }
  }
}

export async function queryLossOrders(limit: number) {
  await ensureConnected();

  const query = `
    WITH items_com_lucro AS (
      SELECT 
        (xpath('//Encomenda/@Id', xml_documento))[1]::text as order_id,
        (xpath('.//Lucro/text()', item))[1]::text::numeric as lucro
      FROM relatorio, LATERAL unnest(xpath('//Item', xml_documento)) as item
    )
    SELECT order_id, lucro as lucro_total
    FROM items_com_lucro 
    WHERE lucro < 0 
    ORDER BY lucro ASC 
    LIMIT $1;
  `;

  try {
    const result = await dbClient.query(query, [limit]);
    console.log(`[DB] Encontrados ${result.rows.length} pedidos com prejuízo`);
    // Mapear order_id -> orderId, lucro_total -> lucroTotal
    return result.rows.map((row: any) => ({
      orderId: row.order_id,
      lucroTotal: row.lucro_total,
    }));
  } catch (err) {
    console.error("[DB] Erro na query:", err);
    throw err;
  }
}
