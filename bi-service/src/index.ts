import { createServer } from "node:http";
import { createYoga } from "graphql-yoga";
import { schema } from "./schema"; 

const yoga = createYoga({ schema });
const server = createServer(yoga);

server.listen(8082, () => {
  console.log("bi-service (GraphQL) up on :8082/graphql");
});
