import { salesByCategoryGrpc } from "./grpc/client.js";
import { profitByRegionXmlRpc } from "./xmlrpc/client.js";
import { queryLossOrders } from "./db/client.js";

export const resolvers = {
  Query: {
    salesByCategory: async (_: any, args: { category: string }) => {
      console.log("[BI] chamando gRPC SalesByCategory", args.category);
      return await salesByCategoryGrpc(args.category);
    },

    profitByRegion: async (_: any, args: { region: string }) => {
      console.log("[BI] chamando XML-RPC ProfitByRegion", args.region);
      return await profitByRegionXmlRpc(args.region);
    },

    lossOrders: async (_: any, args: { limit: number }) => {
      console.log("[BI] consultando pedidos com prejuizo", args.limit);
      return await queryLossOrders(args.limit);
    },
  },
};
