import { salesByCategoryGrpc } from "./grpc/client";
import { profitByRegionXmlRpc } from "./xmlrpc/client";

export const resolvers = {
  Query: {
    salesByCategory: async (_: any, args: { category: string }) => {
      console.log("[BI] calling gRPC SalesByCategory", args.category);
      return await salesByCategoryGrpc(args.category);
    },

    profitByRegion: async (_: any, args: { region: string }) => {
      console.log("[BI] calling XML-RPC ProfitByRegion", args.region);
      return await profitByRegionXmlRpc(args.region);
    },
  },
};
