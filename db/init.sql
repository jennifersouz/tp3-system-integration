-- usar a BD que o container já cria (POSTGRES_DB=tp3db)
-- e o user que o container já cria (POSTGRES_USER=tp3)
CREATE TABLE IF NOT EXISTS relatorio (
  id BIGSERIAL PRIMARY KEY,
  xml_documento XML NOT NULL,
  data_criacao TIMESTAMP NOT NULL DEFAULT now(),
  mapper_version TEXT NOT NULL
);
