---
id: monitor-camadas
project: QMK Userspace
domain: monitor-camadas
kind: contract
status: draft
implementation: planned
origin: "M-1"
created: 2026-11-21
updated: 2026-11-21
supersedes: []
superseded_by:
---

# Spec — Monitor de camadas do Dilemma V3

## Contract

No Linux, o monitor mostra em tempo real o mapa da camada efetiva do Dilemma V3 (`bastardkb/dilemma/3x5_3_procyon`, keymap `vendor`), mesmo quando o foco e a digitação estão em outra aplicação. O firmware informa somente qual é a camada efetiva; o monitor apresenta o mapa já conhecido, sem ler a digitação nem interferir com o teclado ou o VIA.

## Expected behavior

- Ao iniciar com o teclado conectado, o monitor espera a primeira mensagem válida e mostra a camada já efetiva, sem exigir uma tecla ou troca de camada. O firmware anuncia o estado novamente a cada 500 ms, mesmo sem mudança.
- Quando a camada efetiva muda, inclusive ao soltar uma tecla de camada, o firmware anuncia o novo ID sem esperar o heartbeat. A tela passa a mostrar essa camada, seja a Base ou outra camada que continue efetiva segundo a prioridade do QMK.
- Se o teclado for desconectado, a tela sinaliza que não há uma camada atual confirmada. Ao reconectar, o monitor volta a mostrar a camada efetiva assim que recebe uma mensagem válida da nova conexão, sem reiniciar o aplicativo. Até receber essa mensagem, não apresenta a Base nem o último valor como se estivessem confirmados.
- A área superior da TUI mostra a camada atual e as 36 posições do teclado: três fileiras de cinco teclas e três polegares por metade, na ordem física. A área inferior mostra a legenda dos rótulos e abreviações, distinguindo toque, segurar, modificadores, posições transparentes e posições sem ação. Transparência indica herança; não é uma promessa de mostrar o resultado da combinação de camadas.
- Ao redimensionar a janela, o mapa e a legenda se recompõem para as células disponíveis. Nos dois tamanhos-alvo, todas as 36 posições, a camada e a legenda permanecem visíveis ao mesmo tempo, sem rolagem, corte ou sobreposição.

## Invariants

- A camada exibida é informada pelo QMK, não inferida de teclas, do foco da janela ou do VIA. O heartbeat é uma repetição do estado, não um evento de digitação; pressionamentos que não mudam a camada efetiva não geram mensagens de evento.
- O fluxo é unidirecional, do firmware ao host. O monitor não captura tecla, modificador, texto digitado, movimento do ponteiro nem eventos de outros dispositivos; não abre `/dev/input/event*` nem instala captura global. Não registra histórico de camadas, não envia telemetria e não precisa de rede em execução.
- Teclado, trackball e VIA mantêm o funcionamento normal com o monitor aberto. O monitor não abre nem escreve na interface Raw HID do VIA, não configura o teclado e não exige foco para receber a camada.
- A referência dos rótulos e posições é o keymap `vendor` vigente na árvore principal, inclusive suas personalizações locais; em divergência com o mapa visual ou notas, vale o `keymap.c`. Na Base, segurar `Z` **ou** `/` ativa Ponteiro; toque preserva a tecla original. Na camada Ponteiro, ambas as posições são transparentes. Símbolos preserva a ordem, as posições e as ações de toque/segurar do keymap atual, inclusive os pares `-/_`, `;/:`, `[/ {` (colchete abre/chave abre), `]/}` (colchete fecha/chave fecha), `=/+` e `\` com `|`.

## Interfaces

- **Dispositivo:** console USB do QMK com VID `a8f8`, PID `1836`, usage page `0xff31` e usage `0x74`, em hexadecimal. O host usa leitura passiva de `/dev/hidraw*` e só usa a interface que corresponda a essa identidade completa; VID/PID sozinhos não distinguem o console das outras interfaces.
- **Mensagem:** ASCII `DLM1:<id>\n`, com exatamente um dígito de `0` a `6` seguido de LF. Não há bitmask, conjunto de camadas, evento de tecla nem metadados adicionais. IDs: `0` Base, `1` Funções, `2` Navegação, `3` Mídia, `4` Ponteiro, `5` Números, `6` Símbolos.
- **TUI:** aplicação de terminal em Go com Bubble Tea e Lip Gloss; nomes, estados e legenda em português do Brasil. A área útil do terminal é medida em colunas × linhas, não em pixels.
- **Instalação Linux:** `./monitor-camadas/run.sh` prepara as dependências de software, compila e inicia a TUI. As instruções identificam os pré-requisitos e orientam o acesso ao `hidraw` do console por usuário comum. A preparação pode precisar de rede e autorização administrativa; o monitor em si não roda como root. O launcher não grava firmware nem concede acesso geral a dispositivos HID.

## Errors and limits

- **Terminal:** `640 × 360 px` (um quarto de uma tela 720p) corresponde ao alvo compacto de **64 colunas × 18 linhas**; `960 × 540 px` (um quarto de uma tela 1080p), ao alvo de **96 colunas × 27 linhas**. São alvos de área útil: com células de `10 × 20 px` a correspondência é exata; fonte, escala, margens e decoração da janela podem alterá-la. Abaixo de `64 × 18` células, a interface informa o espaço necessário em vez de cortar silenciosamente o mapa.
- **Quadros incompletos ou inválidos:** o host remonta linhas divididas em relatórios HID e tolera padding de bytes zero (`NUL`). Linhas de outros formatos ou IDs fora do intervalo não alteram a camada exibida; uma linha válida posterior deve ser aceita. Fragmentos pendentes da conexão anterior são descartados ao reconectar.
- **Sem dispositivo ou sem permissão:** a interface distingue espera pelo dispositivo de falta de acesso e, no segundo caso, orienta o ajuste da regra `udev`. Dados ausentes ou inválidos não são tratados como Base.

## Conformance criteria

- [ ] Um terminal de `64 × 18` e outro de `96 × 27` células mostram a camada atual, as 36 posições e a legenda completa na mesma tela; redimensionar entre eles não perde informação. Abaixo do mínimo aparece o aviso de tamanho.
- [ ] Uma leitura do console confirma apenas `DLM1:<id>\n` para o monitor, com envio na mudança efetiva e repetição a cada 500 ms; relatórios com padding e mensagens fragmentadas são aceitos, mensagens inválidas são ignoradas e fragmentos não sobrevivem à reconexão.
- [ ] A instalação com `./monitor-camadas/run.sh` permite operar no Linux como usuário comum sem liberar todos os dispositivos HID. Sem a permissão, aparece orientação específica; sem o teclado, aparece o estado de espera.
- [ ] **Confirmação física:** em um Dilemma V3 com o firmware correspondente gravado por ação explícita, abrir o monitor com uma camada já ativa, conferir as sete camadas por acionamento físico, inclusive Ponteiro tanto por `Z` quanto por `/`, e verificar retorno à Base ou à camada ainda efetiva ao soltar. Conferir o mapa de Símbolos e as posições transparentes contra o keymap vigente.
- [ ] **Confirmação física:** com outra aplicação em foco, confirmar que a digitação e o trackball funcionam normalmente e que o VIA continua utilizável enquanto a TUI atualiza a camada sem exibir o que foi digitado. Desconectar e reconectar o USB e confirmar que a TUI se recupera sem reinício nem acionamento de tecla.

## Out of scope

- Captura ou histórico de digitação, visualização de pressionamentos, leitura de estados além da camada efetiva, comandos para o firmware e edição de keymap ou configuração VIA pelo monitor.
- Exibição da pilha de camadas ativas, resolução de transparência para uma combinação concreta de camadas, mapas de giro dos encoders, outros sistemas operacionais e outros modelos ou keymaps.
- Gravação automática do firmware, distribuição de pacotes e mudanças nas ações de tecla para acomodar a TUI.

## Related references

- [[pop/MODIFICATIONS|M-1]] — origem do contrato; cartão [[M-1.1-monitor-camadas|M-1.1]].
- [[visual-layouts/layers_dillema_3x5_3_procyon|Mapa visual do Dilemma 3x5_3 Procyon]] — referência de apresentação; fonte normativa: `keyboards/bastardkb/dilemma/3x5_3_procyon/keymaps/vendor/keymap.c`.
- [[pop/specs/build-flash-script|Build e flash do firmware]] — gravação explícita, separada do launcher do monitor.
