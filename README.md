# Microsservicos (repo unico)

Este repo organiza multiplos servicos em pastas separadas.
Cada servico tem seu proprio `cmd/` e `internal/`.

Servicos:
- consent-service
- policy-service
- audit-service
- report-service
- integration-service

Padrao:
- `cmd/api`
- `internal/`
- `configs/`
- `http/`
# Microsservicos (opcional)

Objetivo: avaliar quando faz sentido dividir o dominio em servicos.

Servicos sugeridos:
- consent-service
- audit-service
- report-service
- policy-service
- integration-service

Recomendacao: iniciar somente apos Onda B, quando ha dor real de escala
ou autonomia de equipes.
