# Monitor de camadas do Dilemma V3

O monitor mostra as 36 posições do keymap `vendor` do Dilemma V3 Procyon e a legenda de toque, segurar e transparência. O firmware envia **somente o número da camada efetiva**: `DLM1:<id>` seguido de LF, na mudança e a cada 500 ms. O mapa é o do projeto, não uma leitura das configurações do VIA. A TUI não lê teclas, texto ou movimentos, não envia comandos ao teclado e continua recebendo a camada com outra janela em foco.

## Iniciar no Linux

Na raiz do projeto, como **usuário comum**, em um terminal:

```sh
./monitor-camadas/run.sh
```

De dentro desta pasta use `./run.sh`; de outro diretório use o caminho para o mesmo script. Para sair, dê foco ao terminal do monitor e pressione **q** ou **Ctrl+C**.

O launcher requer Bash, utilitários usuais do Linux e **Go 1.25 ou superior**. No Arch, se o Go faltar ou for antigo, ele oferece instalar o pacote `go` pelo `pacman`, com sua autorização. Nas demais distribuições, prepare o Go antes. O script baixa as dependências declaradas pelo módulo, compila e inicia a TUI sem root; o executável fica em `${XDG_CACHE_HOME:-$HOME/.cache}/monitor-camadas/`. A primeira preparação pode precisar de rede. Com a toolchain e as dependências em cache, é possível iniciar offline; o monitor em execução não usa a rede. Não é necessário conectar o teclado para abrir a TUI.

### Firmware: gravação separada e opcional

Se o teclado ainda não tem o firmware deste projeto com o anúncio de camadas, **só se quiser gravá-lo**, execute na raiz do checkout que contém o keymap atualizado:

```sh
./flash.sh
```

Esse comando prepara, compila e orienta a gravação das metades direita e esquerda. Siga as instruções de bootloader no terminal. `./flash.sh --compile-only` apenas compila, sem gravar. O launcher do monitor nunca executa o [script de flash](../flash.sh).

### Sem camada na tela?

| Situação | O que fazer |
| --- | --- |
| Nenhum console compatível | Conecte o Dilemma via USB. Sem o teclado, a TUI permanece em espera; firmware sem interface de console também não é detectado. |
| Console conectado, sem mensagem válida | Com o firmware correspondente, o heartbeat informa a camada já ativa sem precisar pressionar tecla. Se a espera persistir, confira o firmware: uma versão antiga pode ter console, mas não enviar `DLM1`. |
| Permissão negada | Com o teclado conectado, feche a TUI e rode o launcher novamente na sua sessão local ativa. Ele oferece a regra específica descrita abaixo. |
| USB desconectado durante o uso | Não há camada confirmada. Reconecte: a leitura recomeça, sem reiniciar a TUI, assim que chegar outra mensagem válida. |

**Permissão do console.** Se o usuário já puder ler o `hidraw` do console, o launcher não instala regra. Caso contrário, ele identifica o USB `a8f8:1836` e valida que o descritor tem exclusivamente a coleção do console QMK (usage page `ff31`, usage `74`). Só então oferece instalar, com `sudo`, uma regra em `/etc/udev/rules.d/70-dilemma-console-<hash>.rules`. A regra exige o VID/PID **e o SHA-256 desse descritor** a cada conexão: o número de `/dev/hidrawN` pode mudar, e teclado, ponteiro e VIA não recebem acesso por essa regra.

No Arch com `systemd-logind`, a tag `uaccess` concede acesso à sessão local ativa. Ao aceitar, o script recarrega as regras com `udevadm` e as reaplica ao nó identificado. Se a leitura ainda for negada, desconecte e reconecte o cabo USB. Sem `sudo`/`udevadm`, ou em uma sessão sem gerenciamento de acesso local (por exemplo, SSH), peça ao administrador acesso de leitura apenas à mesma coleção. Recusar a instalação não impede abrir a TUI. Se uma atualização mudar o descritor, rode novamente o launcher com o teclado conectado.

## Conferência no teclado físico

1. Com uma camada já ativa, abra o monitor: ela deve aparecer sem nova tecla. Confira Base, Funções, Navegação, Mídia, Ponteiro, Números e Símbolos; compare as posições com o keymap `vendor`, inclusive os pares de símbolos e as posições transparentes.
2. Digite em outra janela. Sobreponha camadas e solte uma delas: a tela deve voltar à camada ainda efetiva (ou à Base), sem exibir a digitação. Na Base, teste segurar `Z` e `/` separadamente para Ponteiro; o toque deve continuar normal. Use o trackball e o VIA com o monitor aberto.
3. Remova o cabo e confira que a camada antiga não fica indicada como atual. Reconecte: a camada deve reaparecer sem reiniciar ou acionar tecla.
4. Redimensione a **área útil** do terminal para **64 × 18** e **96 × 27** colunas × linhas (`stty size` informa linhas e depois colunas). São aproximadamente quartos de telas 720p e 1080p com células de 10 × 20 px. Confira as 36 posições, o nome da camada e a legenda juntos, sem rolagem ou corte. Abaixo de 64 × 18 deve aparecer o aviso de espaço insuficiente.
