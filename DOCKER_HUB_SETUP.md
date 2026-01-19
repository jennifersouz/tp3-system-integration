# Configuração de Secrets para Docker Hub

## ⚠️ IMPORTANTE: Nunca coloque credenciais no repositório!

Os secrets do GitHub devem ser adicionados **apenas** através da interface web, nunca diretamente no código.

## Como adicionar os secrets com segurança:

### Passo 1: Acesse seu repositório no GitHub
```
https://github.com/jennifersouz/tp3-system-integration
```

### Passo 2: Vá até Settings
- Clique em **Settings** (engrenagem no topo à direita)

### Passo 3: Acesse Secrets
- No menu lateral, vá para **Secrets and variables** > **Actions**

### Passo 4: Crie o primeiro secret - DOCKER_HUB_USERNAME
- Clique em **"New repository secret"**
- **Name:** `DOCKER_HUB_USERNAME`
- **Value:** `<seu-usuario-docker-hub>`
- Clique em **"Add secret"**

### Passo 5: Crie o segundo secret - DOCKER_HUB_TOKEN
- Clique em **"New repository secret"**
- **Name:** `DOCKER_HUB_TOKEN`
- **Value:** `<seu-token-docker-hub>`
- Clique em **"Add secret"**

## ✅ Resultado Final

Após adicionar os secrets:

1. **O workflow no `.github/workflows/docker-build.yml` usará automaticamente esses secrets**
2. **Quando você fizer push para `main` ou `develop`**, o GitHub Actions vai:
   - Fazer checkout do código
   - Fazer build de todas as 5 imagens Docker
   - Fazer login no Docker Hub usando suas credenciais
   - Fazer push das imagens com tags de versão

3. **Suas credenciais nunca aparecerão nos logs ou no repositório**

## 🔐 Segurança

- Os secrets são **criptografados** pelo GitHub
- Aparecem como `***` nos logs
- Não são clonados junto com o repositório
- Cada desenvolvedor pode ter seus próprios secrets

## 📋 Triggers do Workflow

O workflow será acionado automaticamente em:
- ✅ Push para `main`
- ✅ Push para `develop`
- ✅ Criação de tags `v*` (ex: `v1.0.0`)

