# TP3: Arquitetura de Integração de Sistemas com Múltiplos Protocolos e Formatos de Dados

## Sumário 

Este trabalho prático implementa uma arquitetura distribuída para integração e interoperabilidade de sistemas heterogéneos. O sistema demonstra a aplicação de linguagens de anotação (XML) para representação de dados e utiliza múltiplas linguagens de programação, protocolos de comunicação e formatos de dados para resolver um problema real de integração de sistemas.

**Objetivos alcançados:**
- Desenvolvimento de 5 serviços distribuídos com 3 linguagens distintas
- Implementação de 4 protocolos de comunicação (REST, GraphQL, XML-RPC, gRPC)
- Persistência de dados em formato XML com consultas XPath
- Comunicação assíncrona via webhooks
- Orquestração containerizada com Docker Compose

---

## 1. Introdução

### 1.1 Contexto

A integração de sistemas heterogéneos constitui um dos principais desafios da arquitetura de software moderna. Diferentes organizações utilizam tecnologias e formatos de dados distintos, dificultando a comunicação entre sistemas. Este projeto demonstra como utilizar linguagens de anotação (especificamente XML) como mecanismo de interoperabilidade.

### 1.2 Objectivos

1. Desenvolver uma solução de integração de sistemas utilizando múltiplas linguagens de programação
2. Implementar comunicação através de diversos protocolos e formatos de dados
3. Demonstrar a capacidade do XML para representação de dados desacoplada da origem
4. Validar a eficácia do modelo de comunicação assíncrona via webhooks
5. Implementar consultas complexas em dados XML persistidos em base de dados relacional

### 1.3 Escopo

O sistema compreende os seguintes componentes:
- Origem de dados (Crawler)
- Processamento assíncrono (Processor)
- Serviço de transformação (XML-Service)
- Interface de consumo de dados (BI-Service)
- API complementar (External-API)

---

## 2. Arquitetura do Sistema

### 2.1 Visão Geral
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

![Arquitetura do Sistema](docs/tp3.jpg)
```

### 2.2 Componentes

#### 2.2.1 Crawler (Python)
- **Função**: Geração periódica de dados de entrada
- **Implementação**: Script Python que gera CSV com dados de vendas simulados
- **Frequência**: A cada 30 segundos (configurável via `CRAWLER_PERIOD_SECONDS`)
- **Armazenamento**: Supabase S3 (protocolo: File Storage/S3 API)
- **Formato de saída**: CSV compatível com schema Superstore

#### 2.2.2 Processor (TypeScript/Node.js)
- **Função**: Orquestrador de processamento assíncrono
- **Responsabilidades**:
  - Monitorizar S3 por novos CSVs (polling a cada 15 segundos)
  - Baixar e ler dados em stream
  - Enriquecer dados com API externa
  - Gerar novo CSV com atributos desacoplados
  - Enviar para XML-Service e aguardar confirmação
  - Remover CSVs processados de S3
- **Comunicação**: REST (Multipart-Form), Webhook
- **Persistência**: Sem BD local (stateless)

#### 2.2.3 XML-Service (Go)
- **Função**: Transformação, validação e persistência de dados em XML
- **Responsabilidades**:
  - Receber CSV via multipart-form
  - Validar estrutura e dados
  - Gerar XML hierárquico a partir de CSV
  - Enriquecer XML com dados de API externa (taxa de imposto)
  - Persistir XML em coluna tipo XML do PostgreSQL
  - Executar consultas XPath
  - Notificar Processor via webhook
- **Comunicação**: 
  - REST (ingest, query endpoints)
  - gRPC (SalesByCategory)
  - XML-RPC (ProfitByRegion)
- **Base de dados**: PostgreSQL com tipo de dados XML

#### 2.2.4 BI-Service (TypeScript/Node.js)
- **Função**: Interface de consumo de dados com agregações e transformações
- **Tipo de interface**: GraphQL
- **Queries implementadas**:
  - `salesByCategory(category: String!): Float!` (via gRPC)
  - `profitByRegion(region: String!): Float!` (via XML-RPC)
  - `lossOrders(limit: Int): [LossOrder!]!` (via REST/XPath)
- **Comunicação multiprotocolo**:
  - gRPC para SalesByCategory
  - XML-RPC para ProfitByRegion
  - REST para LossOrders

#### 2.2.5 External-API (Python/Flask)
- **Função**: Fornecimento de dados complementares
- **Dados fornecidos**: Taxas de imposto por estado
- **Comunicação**: 
  - REST (GET /tax?state=State)
  - XML-RPC (getTaxRate method)
- **Dados**: Mock com 4 estados (Kentucky, California, Florida, New York)

### 2.3 Fluxo de Dados

1. **Origem**: Crawler gera CSV e armazena em S3 (`input/orders_*.csv`)
2. **Detecção**: Processor faz polling em S3 a cada 15 segundos
3. **Enriquecimento**: Processor consulta External-API para taxa de imposto
4. **Transformação**: Processor envia CSV + metadados para XML-Service via multipart-form
5. **Validação**: XML-Service valida CSV e gera XML hierárquico
6. **Persistência**: XML-Service armazena em PostgreSQL com `INSERT ... RETURNING id`
7. **Notificação**: XML-Service notifica Processor via webhook POST
8. **Limpeza**: Processor remove CSV de S3 após confirmação
9. **Consultas**: BI-Service consulta XML-Service via gRPC/XML-RPC/REST
10. **Agregação**: BI-Service formata dados em GraphQL e retorna ao cliente

---

## 3. Especificação Técnica

### 3.1 Linguagens de Programação

| Serviço | Linguagem | Versão | Justificação |
|---------|-----------|--------|--------------|
| Crawler | Python | 3.11 | Rapidez de desenvolvimento, bibliotecas S3 maduras |
| Processor | TypeScript | 5.0 | Type-safe, Node.js para I/O assíncrono |
| XML-Service | Go | 1.21 | Performance, goroutines, suporte nativo gRPC |
| BI-Service | TypeScript | 5.0 | Consistência com Processor, suporte GraphQL |
| External-API | Python | 3.11 | Flask lightweight, desenvolvimento rápido |

### 3.2 Protocolos de Comunicação

#### 3.2.1 REST (HTTP/JSON)
- **Utilização**: Processor → XML-Service (ingest), Queries simples
- **Endpoints**:
  - `POST /ingest` - Receber CSV multipart-form
  - `GET /query/vendas-por-categoria?categoria=X` - Consulta XPath
  - `GET /query/lucro-por-regiao?regiao=X` - Consulta XPath
  - `GET /query/encomendas-prejuizo` - Consulta XPath complexa
  - `GET /tax?state=X` - External-API

#### 3.2.2 GraphQL
- **Utilização**: Interface principal do BI-Service
- **Endpoint**: `POST /graphql`

#### 3.2.3 gRPC (Protocol Buffers)
- **Utilização**: BI-Service → XML-Service (SalesByCategory)
- **Porto**: 50051

#### 3.2.4 XML-RPC
- **Utilização**: BI-Service → XML-Service (ProfitByRegion)
- **Método**: `XMLRPCServer.ProfitByRegion`
- **Porto**: 8099

### 3.3 Base de Dados

#### 3.3.1 Schema PostgreSQL
```sql
CREATE TABLE relatorio (
  id BIGSERIAL PRIMARY KEY,
  xml_documento XML NOT NULL,
  data_criacao TIMESTAMP NOT NULL DEFAULT now(),
  mapper_version TEXT NOT NULL,
  id_requisicao VARCHAR(100),
  descricao TEXT
);
```

#### 3.3.2 Consultas XPath

**Query 1: Vendas por categoria**
```sql
SELECT COALESCE(SUM((v)::numeric), 0) AS total
FROM unnest(xpath('//Item[@Categoria=$1]/ValorVenda/text()', xml_documento))::text AS v;
```

**Query 2: Lucro por região**
```sql
SELECT COALESCE(SUM((p)::numeric), 0) AS lucro_total
FROM unnest(xpath('.//Item/Lucro/text()', e))::text AS p;
```

**Query 3: Encomendas com prejuízo**
```sql
SELECT order_id, lucro_total FROM (
  SELECT COALESCE(SUM((x::text)::numeric), 0) AS lucro_total
  FROM unnest(xpath('.//Item/Lucro/text()', e)) AS t(x)
) WHERE lucro_total < 0 ORDER BY lucro_total ASC;
```

---

## 4. Instalação e Configuração

### 4.1 Pré-requisitos

- Docker 20.10+
- Docker Compose 1.29+
- Git

### 4.2 Instalação Local

```bash
# Clone do repositório
git clone https://github.com/seu-usuario/tp3-system-integration.git
cd tp3-system-integration

# Configurar variáveis de ambiente
cp .env.example .env
# Editar .env com credenciais Supabase

# Iniciar serviços
docker-compose up -d

# Verificar status
docker-compose ps
```

### 4.3 Configuração de Variáveis

```bash
# .env
S3_ENDPOINT=https://[project].supabase.co/storage/v1/s3
S3_REGION=eu
S3_ACCESS_KEY_ID=***
S3_SECRET_ACCESS_KEY=***
S3_BUCKET=tp3-data

DB_HOST=localhost
DB_PORT=5432
DB_USER=tp3
DB_PASSWORD=tp3
DB_NAME=tp3db

XML_SERVICE_URL=http://xml-service:8081/ingest
EXTERNAL_API_URL=http://external-api:8090/tax
WEBHOOK_URL=http://processor:8080/webhook/xml-status

CRAWLER_PERIOD_SECONDS=30
POLL_SECONDS=15
```

---

## 6. Referências

- Docker Documentation: https://docs.docker.com/
- PostgreSQL XML Functions: https://www.postgresql.org/docs/current/functions-xml.html
- gRPC Protocol Buffers: https://grpc.io/
- GraphQL Specification: https://spec.graphql.org/
- XML-RPC Specification: http://xmlrpc.scripting.com/

---

## 7. Informação do Projeto

**Disciplina:** Integração de Sistemas  
**Instituto:** IPVC  
**Período:** Janeiro 2026  
**Por:** Jennifer Silva e Giulia Campos


