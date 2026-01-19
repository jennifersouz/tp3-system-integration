-- Create users and database if not exists
CREATE USER sistema_user WITH PASSWORD 'sistema_pass';
CREATE DATABASE sistema;

-- Grant privileges
GRANT ALL PRIVILEGES ON DATABASE sistema TO sistema_user;

-- Connect to the database
\c sistema

-- Create table
CREATE TABLE IF NOT EXISTS relatorio (
  id BIGSERIAL PRIMARY KEY,
  xml_documento XML NOT NULL,
  data_criacao TIMESTAMP NOT NULL DEFAULT now(),
  mapper_version TEXT NOT NULL
);

-- Grant table privileges
GRANT ALL PRIVILEGES ON TABLE relatorio TO sistema_user;
GRANT USAGE, SELECT ON SEQUENCE relatorio_id_seq TO sistema_user;
