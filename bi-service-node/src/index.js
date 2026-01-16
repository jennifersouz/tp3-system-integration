import { createServer } from "http";
import { createYoga, createSchema } from "graphql-yoga";

const XML_BASE = process.env.XML_SERVICE_BASE ?? "http://xml-service:8081";

async function getJson(url) {
  const resp = await fetch(url);
  const text = await resp.text();
  try { return JSON.parse(text); } catch { return { raw: text }; }
}

const typeDefs = /* GraphQL */ `
  type LossOrder {
    orderId: String!
    lucroTotal: Float!
  }

  type Query {
    salesByCategory(category: String!): Float!
    profitByRegion(region: String!): Float!
    lossOrders(limit: Int = 100): [LossOrder!]!
  }
`;

const resolvers = {
  Query: {
    salesByCategory: async (_, { category }) => {
      const data = await getJson(`${XML_BASE}/query/vendas-por-categoria?categoria=${encodeURIComponent(category)}`);
      return Number(data.totalVendas ?? 0);
    },
    profitByRegion: async (_, { region }) => {
      const data = await getJson(`${XML_BASE}/query/lucro-por-regiao?regiao=${encodeURIComponent(region)}`);
      return Number(data.lucroTotal ?? 0);
    },
    lossOrders: async (_, { limit }) => {
      const data = await getJson(`${XML_BASE}/query/encomendas-prejuizo`);
      const list = Array.isArray(data.results) ? data.results : [];
      return list.slice(0, limit).map(x => ({
        orderId: String(x.orderId),
        lucroTotal: Number(x.lucroTotal)
      }));
    }
  }
};

const yoga = createYoga({
  schema: createSchema({ typeDefs, resolvers }),
  graphqlEndpoint: "/graphql"
});

const server = createServer(yoga);
server.listen(8082, () => {
  console.log("bi-service (GraphQL) up on :8082/graphql");
});
