# TP3 - Integração de Sistemas

Arquitetura distribuída para integração e interoperabilidade de sistemas com múltiplos protocolos, linguagens e formatos de dados.

## Arquitetura

**Serviços (4+)**:
- **Crawler**: Gera dados CSV periodicamente (Python)
- **Processor**: Enriquece CSV com dados externos, valida e envia para XML-Service (Node.js)
- **XML-Service**: Converte para XML, valida XSD, persiste em PostgreSQL (Go/gRPC)
- **BI-Service**: Expõe GraphQL/REST para consultas complexas com XPath (Node.js)
- **External-API**: Fornece dados complementares via XML-RPC/REST (Python)

**Protocolos (4+)**:
- REST (Processor → XML-Service)
- GraphQL (BI-Service consultas)
- XML-RPC (External-API)
- gRPC (XML-Service persistência)

**Dados**: XML com hierarquia + JSON

**Armazenamento**: Supabase S3 + PostgreSQL (XML com XPath)

## Quick Start

```bash
cp .env.example .env
docker compose up -d
```

**Endpoints**:
- Processor REST: http://localhost:8081/ingest
- BI-Service GraphQL: http://localhost:8082/graphql
- XML-Service gRPC: localhost:50051

## Estrutura

```
├── crawler/          # Python - Gerador CSV
├── processor/        # TypeScript - Orquestrador
├── xml-service/      # Go - Serviço XML/gRPC
├── bi-service/       # Node.js - GraphQL
├── external-api/     # Python - XML-RPC
├── db/               # PostgreSQL init
└── docker-compose.yml
```

As imagens são automaticamente construídas e publicadas no Docker Hub via GitHub Actions.

**Configuração necessária em Secrets do GitHub:**
- `DOCKER_HUB_USERNAME`: seu nome de usuário Docker Hub
- `DOCKER_HUB_TOKEN`: token de acesso pessoal Docker Hub

**Versionamento:**
- Branch `main` → tag `latest`
- Tags `v*` → tag da versão (ex: `v1.0.0`)
- Outras branches → tag `{branch}-{commit}`

Exemplo:
```bash
docker pull seuusuario/tp3-xml-service:latest
docker pull seuusuario/tp3-xml-service:v1.0.0
```

### 3. AWS ECS

```bash
# Configurar variáveis de ambiente
export DOCKER_REGISTRY=your-ecr-registry
export IMAGE_TAG=v1.0.0
export S3_ENDPOINT=https://your-s3.amazonaws.com
export S3_BUCKET=your-bucket

# Deploy com docker-compose
docker compose -f docker-compose.ecs.yml up

# Ou criar cluster ECS e task definition via AWS CLI
aws ecs register-task-definition --cli-input-json file://task-definition.json
aws ecs create-service --cluster tp3-cluster --service-name tp3 --task-definition tp3:1
```

**Recomendações AWS:**
- RDS para PostgreSQL (Aurora)
- S3 para bucket de armazenamento
- Application Load Balancer para rotas
- ECR para repositório de imagens privadas

### 4. Azure Container Instances

```bash
# Configurar variáveis
export DOCKER_REGISTRY=your-acr-registry
export IMAGE_TAG=v1.0.0

# Deploy com docker-compose
docker compose -f docker-compose.azure.yml up

# Ou criar resource group e container
az container create \
  --resource-group mygroup \
  --name tp3-system \
  --image your-acr-registry/tp3-xml-service:latest \
  --ports 8081 50051 50052

# Usar Azure Database for PostgreSQL para DB relacional
```

**Recomendações Azure:**
- Azure Database for PostgreSQL (com suporte XML)
- Azure Blob Storage para S3-compatible storage
- Azure Container Registry para imagens privadas
- Azure Application Gateway para routing

## 📡 Protocolos Implementados

### REST (Multipart-Form)
- **Processor → XML-Service**: `POST /ingest` com CSV em multipart/form-data
- **BI-Service → XML-Service**: `GET /query/*` para consultas XPath

### GraphQL
- **BI-Service**: `POST /graphql` com queries de vendas, lucro e encomendas
  ```graphql
  query {
    salesByCategory(category: "Electronics")
    profitByRegion(region: "South")
    lossOrders(limit: 100)
  }
  ```

### XML-RPC
- **External-API**: `POST /RPC2` método `getTaxRate(state)`
  ```xml
  <?xml version="1.0"?>
  <methodCall>
    <methodName>getTaxRate</methodName>
    <params>
      <param><value><string>California</string></value></param>
    </params>
  </methodCall>
  ```

### gRPC
- **Ingest Service** (porta 50051):
  - `rpc IngestCSV(IngestRequest) returns (IngestResponse)`
- **Query Service** (porta 50052):
  - `rpc GetSalesByCategory(CategoryQuery) returns (SalesResult)`
  - `rpc GetProfitByRegion(RegionQuery) returns (ProfitResult)`
  - `rpc GetLossOrders(LimitQuery) returns (LossOrdersResult)`

## 📊 Fluxo de Dados

1. **Crawler** gera CSV periodicamente → S3 (Supabase)
2. **Processor** monitora S3 → baixa CSV → chama External-API (XML-RPC) → envia para XML-Service (gRPC ou REST)
3. **XML-Service** valida CSV → cria XML desacoplado → persiste em PostgreSQL → envia webhook para Processor
4. **BI-Service** expõe GraphQL → chama XML-Service (gRPC ou REST) → retorna dados transformados
5. **Visualização** consome BI-Service (GraphQL)

## 🗄️ Schema XML

```xml
<RelatorioVendas DataGeracao="2024-01-15" VersaoMapper="1.0">
  <Encomendas>
    <Encomenda Id="123" Data="2024-01-10" DataEnvio="2024-01-12" ModoEnvio="Standard Class">
      <Cliente Id="C1" Segmento="Consumer" Nome="John Doe">
        <Localizacao Pais="USA" Cidade="New York" Estado="New York" CodigoPostal="10001" Regiao="East">
          <Imposto Taxa="0.088"/>
        </Localizacao>
      </Cliente>
      <Vendedor>Sales Person A</Vendedor>
      <Itens>
        <Item ProdutoId="P1" Categoria="Technology" SubCategoria="Copiers">
          <NomeProduto>Copier X200</NomeProduto>
          <Devolvido>false</Devolvido>
          <Quantidade>2</Quantidade>
          <ValorVenda>500.00</ValorVenda>
          <Desconto>0.10</Desconto>
          <Lucro>100.00</Lucro>
        </Item>
      </Itens>
    </Encomenda>
  </Encomendas>
</RelatorioVendas>
```

## 📝 Variáveis de Ambiente

Criar `.env` com:

```bash
# Supabase S3
S3_ENDPOINT=https://your-project.supabase.co/storage/v1/s3
S3_REGION=your-region
S3_ACCESS_KEY_ID=your-access-key
S3_SECRET_ACCESS_KEY=your-secret-key
S3_BUCKET=your-bucket

# Database
DB_PASSWORD=tp3

# URLs internas (docker-compose)
XML_SERVICE_URL=http://xml-service:8081/ingest
EXTERNAL_API_URL=http://external-api:8090/tax
WEBHOOK_URL=http://processor:8080/webhook/xml-status
```

## 🧪 Testes

### GraphQL Query
```bash
curl -X POST http://localhost:8082/graphql \
  -H "Content-Type: application/json" \
  -d '{"query": "{ salesByCategory(category: \"Technology\") }"}'
```

### XML-RPC Call
```bash
curl -X POST http://localhost:8091/RPC2 \
  -H "Content-Type: text/xml" \
  -d '<?xml version="1.0"?><methodCall><methodName>getTaxRate</methodName><params><param><value><string>California</string></value></param></params></methodCall>'
```

### gRPC Query
```bash
grpcurl -plaintext -d '{"category":"Technology"}' localhost:50052 xmlservice.QueryService.GetSalesByCategory
```

## 📚 Documentação dos Protocolos

- [REST](./docs/rest-api.md)
- [GraphQL](./docs/graphql-schema.md)
- [XML-RPC](./docs/xmlrpc-methods.md)
- [gRPC](./docs/grpc-services.md)
