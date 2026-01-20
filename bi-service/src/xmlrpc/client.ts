import xmlrpc from "xmlrpc";

const XML_RPC_URL = process.env.XML_RPC_URL || "http://xml-service:8099/RPC2";

const client = xmlrpc.createClient({ url: XML_RPC_URL });

export function profitByRegionXmlRpc(region: string): Promise<number> {
  return new Promise((resolve, reject) => {
    client.methodCall(
      "XMLRPCServer.ProfitByRegion",
      [region],
      (err: any, value: any) => {
        if (err) {
          console.error("[xmlrpc] error:", err);
          return reject(err);
        }
        resolve(Number(value?.lucroTotal ?? 0));
      }
    );
  });
}
