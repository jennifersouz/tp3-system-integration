import { createServer } from "node:http";
import { createYoga } from "graphql-yoga";
import { schema } from "./schema.js"; 

const yoga = createYoga({ schema });
const server = createServer(yoga);

server.listen(8082, () => {
  console.log("bi-service (GraphQL) iniciado em :8082/graphql");
});
