CREATE TABLE IF NOT EXISTS relatorio (
  id BIGSERIAL PRIMARY KEY,
  xml_documento XML NOT NULL,
  data_criacao TIMESTAMP NOT NULL DEFAULT now(),
  mapper_version TEXT NOT NULL
);
