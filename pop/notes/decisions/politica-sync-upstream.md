---
author: agent
created: 2026-08-27
---

# Política de sync com o upstream — moving targets e validação pelo build local

> Como lidamos com o fato de que este userspace depende de upstreams em movimento?

Origem: [[1.3.1-notas-decisoes]] · Contrato do build: [[pop/specs/build-flash-script|build-flash-script]] — *follow para o caminho único de build/flash*.

## Fato

O repo depende de **dois moving targets**, ambos na branch `main` do upstream, sem pinning de versão:

1. **Submódulo `modules/bastardkb`** → Bastardkb/qmk_modules, atualizado automaticamente pelo workflow `update-submodule.yml` (detalhes em [[auto-update-submodulo]] — *follow para o mecanismo exato*).
2. **Firmware** `bastardkb/bastardkb-qmk` (branch `main`), resolvido pelo `qmk setup -b main` no preflight de `scripts/build_flash.sh` quando `user.qmk_home` não está configurado.

Não há testes no projeto. Portanto, a política de sync é:

- **Aceitar o movimento:** updates do submódulo entram por commit automático em `main`; não se congela referência.
- **Validar com build local:** a confirmação de que qualquer sync (submódulo novo ou firmware novo) não quebrou o userspace é o compile local do target principal:
  `qmk compile -kb bastardkb/dilemma/3x5_3_procyon -km vendor` — na prática via `scripts/build_flash.sh`, que faz o preflight do submódulo e do `user.qmk_home` antes de compilar.
- **Nunca validar na CI:** compile remoto é proibido (Decision 2026-08-24 do [[pop/PROJECT|PROJECT]] e [[F-20260824-desativa-ci]] — *follow para o motivo e o registro, não repetidos aqui*).

## Consequência prática

- Um commit `chore: update submodule reference to latest` do bot **não é garantia de build** — ao puxar um desses, rode o build local antes de confiar.
- Se o build quebrar após um sync, as opções são fix manual no userspace ou revert do commit de update do submódulo; não há rollback automatizado.

## Links

- [[auto-update-submodulo]] — *follow para o que o workflow de update faz e deixa de fazer*.
- [[pop/specs/build-flash-script|build-flash-script]] — *follow para o contrato do script de build/flash (preflight, compile, flash opt-in)*.
- [[pop/PROJECT|PROJECT]] — *Context* ("Frágil/a observar") e *Decisions* — *follow para o registro canônico da fragilidade dos moving targets*.
