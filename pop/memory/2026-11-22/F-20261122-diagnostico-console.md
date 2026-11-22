---
task: F-20261122-diagnostico-console
project: qmk-userspace
started: 2026-11-22
finished: 2026-11-22
commit: null
pr: null
authorization: pedido explícito de firmware de diagnóstico
---

# Console silencioso do Dilemma

- **Entrega:** UF2 instrumentado e consulta manual independente da TUI.
- **Verificação:** build local `./flash.sh --compile-only`; bancada exercitou os nove estados, dois heartbeats e a recusa de alteração. O VIA do teclado conectado respondeu que a nova consulta ainda não é suportada. Sem gravação física nem conclusão sobre a causa.
- **Contrato:** consulta avulsa e guia sincronizados; o protocolo do monitor permanece inalterado.

## Entradas

- [[F-20261122-diagnostico-console.01-sonda]] — consultar para separar o emissor do transporte.
