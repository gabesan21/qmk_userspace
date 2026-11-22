---
task: F-20261121-flash-dilemma
project: qmk-userspace
started: 2026-11-21
finished: 2026-11-21
commit: null
pr: null
authorization: pedido explícito de guia executável na raiz com instalação de dependências
---

# Guia de flash do Dilemma V3

- **Entrega:** `./flash.sh` verifica dependências e orienta build e flash das duas metades por padrão. O motor valida o volume UF2 antes de escrever.
- **Verificação:** `./flash.sh --compile-only` em ambiente QMK limpo gerou o UF2. CLI simulada cobriu uma e duas metades, alvo inválido/ambíguo e instalação cancelada. Flash físico não executado.
- **Impacto:** contrato de build/flash e comando de uso atualizados.

## Entradas

- [[F-20261121-flash-dilemma.01-guia]] — consultar para preparação e sequência de gravação.
