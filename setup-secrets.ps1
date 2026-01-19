# Script para adicionar secrets ao GitHub Actions de forma segura
# Você precisa fazer isso MANUALMENTE via interface do GitHub

Write-Host "===================================" -ForegroundColor Green
Write-Host "GitHub Secrets Setup" -ForegroundColor Green
Write-Host "===================================" -ForegroundColor Green
Write-Host ""

Write-Host "INSTRUÇÕES PARA ADICIONAR SECRETS:" -ForegroundColor Yellow
Write-Host ""

Write-Host "1. Acesse seu repositório no GitHub:" -ForegroundColor Cyan
Write-Host "   https://github.com/jennifersouz/tp3-system-integration" -ForegroundColor White
Write-Host ""

Write-Host "2. Vá até: Settings > Secrets and variables > Actions" -ForegroundColor Cyan
Write-Host ""

Write-Host "3. Clique em 'New repository secret'" -ForegroundColor Cyan
Write-Host ""

Write-Host "4. ADICIONE O PRIMEIRO SECRET:" -ForegroundColor Yellow
Write-Host "   Name: DOCKER_HUB_USERNAME" -ForegroundColor Magenta
Write-Host "   Value: <seu-usuario-docker-hub>" -ForegroundColor Magenta
Write-Host "   (Clique 'Add secret')" -ForegroundColor Magenta
Write-Host ""

Write-Host "5. ADICIONE O SEGUNDO SECRET:" -ForegroundColor Yellow
Write-Host "   Name: DOCKER_HUB_TOKEN" -ForegroundColor Magenta
Write-Host "   Value: <seu-token-docker-hub>" -ForegroundColor Magenta
Write-Host "   (Clique 'Add secret')" -ForegroundColor Magenta
Write-Host ""

Write-Host "===================================" -ForegroundColor Green
Write-Host "PRONTO!" -ForegroundColor Green
Write-Host "===================================" -ForegroundColor Green
Write-Host ""
Write-Host "Após adicionar os secrets, o workflow será:" -ForegroundColor Cyan
Write-Host "- Acionado ao fazer push para 'main' ou 'develop'" -ForegroundColor White
Write-Host "- Acionado ao criar tags v* (ex: v1.0.0)" -ForegroundColor White
Write-Host "- As imagens Docker serão enviadas para seu Docker Hub" -ForegroundColor White
Write-Host ""
