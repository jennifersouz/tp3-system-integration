export const typeDefs = /* GraphQL */ `
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
