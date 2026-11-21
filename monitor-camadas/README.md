# Monitor de camadas do Dilemma V3

O monitor mostra as 36 posições do keymap `vendor` do Dilemma V3 Procyon e a legenda de toque, segurar e transparência. Os rótulos das teclas e das camadas ficam no [catálogo do aplicativo](layout.go), não são lidos do firmware nem do VIA. O firmware envia **somente o número da camada efetiva**: `DLM1:<id>` seguido de LF, na mudança e a cada 500 ms. A TUI não lê teclas, texto ou movimentos, não envia comandos ao teclado e continua recebendo a camada com outra janela em foco.

Se o teclado ainda usa firmware antigo, **1. atualize o firmware** e **2. inicie o monitor**. Se ele já envia `DLM1`, pule a primeira etapa. É possível abrir a TUI sem atualizar, mas ela só mostra a camada depois de receber um anúncio válido.

## Atualizar o firmware (opcional e separado do monitor)

Com o [ambiente QMK da BastardKB já configurado](https://docs.bastardkb.com/fw/compile-firmware.html#compiling-your-firmware-using-the-local-command-line), rode na raiz deste checkout, inclusive em um clone novo:

```sh
git submodule update --init -- modules/bastardkb
QMK_USERSPACE="$PWD" qmk compile -kb bastardkb/dilemma/3x5_3_procyon -km vendor
```

Isso só prepara o submódulo e compila localmente, sem gravar no teclado. Após uma compilação bem-sucedida, procure `bastardkb_dilemma_3x5_3_procyon_vendor.uf2` no diretório do firmware QMK (`user.qmk_home`).

Para gravar, conecte por USB a metade direita. Se o firmware já gravado tiver o atalho `BOOT`, segure a tecla inferior externa (`/` na direita ou `Z` na esquerda) para ativar a camada Ponteiro (`PTR`) e, sem soltá-la, toque a tecla superior externa da mesma metade (`BOOT`). Se o atalho não estiver disponível, consulte o [guia oficial da BastardKB](https://docs.bastardkb.com/fw/flashing.html) para escolher um método compatível que preserve suas configurações. No bootloader, copie o `.uf2` para o volume `RPI-RP2` e aguarde o volume desaparecer e o teclado reiniciar antes de retirar o USB. Passe o USB para a metade esquerda e repita com o mesmo arquivo. Deixe o cabo TRRS entre as metades no lugar e volte a ligar o USB na direita para usar o teclado.

Se o seu checkout já tiver o script local `./flash.sh` (não incluído no repositório), você pode usá-lo como alternativa para compilar e gravar as duas metades; `./flash.sh --compile-only` só compila. O launcher do monitor não compila nem grava firmware.

## Iniciar no Linux

Na raiz do projeto, como **usuário comum**, em um terminal:

```sh
./monitor-camadas/run.sh
```

De dentro desta pasta use `./run.sh`; de outro diretório use o caminho para o mesmo script. Para sair, dê foco ao terminal do monitor e pressione **q** ou **Ctrl+C**.

O launcher requer Bash, utilitários usuais do Linux e **Go 1.25 ou superior**. No Arch, se o Go faltar ou for antigo, ele oferece instalar o pacote `go` pelo `pacman`, com sua autorização. Nas demais distribuições, prepare o Go antes. O script baixa as dependências declaradas pelo módulo, compila e inicia a TUI sem root; o executável fica em `${XDG_CACHE_HOME:-$HOME/.cache}/monitor-camadas/`. A primeira preparação pode precisar de rede. Com a toolchain e as dependências em cache, é possível iniciar offline; o monitor em execução não usa a rede. Não é necessário conectar o teclado para abrir a TUI.

## Sem camada na tela?

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
2. Digite em outra janela. Sobreponha camadas e solte uma delas: a tela deve voltar à camada ainda efetiva (ou à Base), sem exibir a digitação. Na Base, teste segurar `Z` e `/` separadamente para Ponteiro; o toque deve continuar normal. Use o trackpad (não uma bola) do Procyon e o VIA com o monitor aberto.
3. Remova o cabo e confira que a camada antiga não fica indicada como atual. Reconecte: a camada deve reaparecer sem reiniciar ou acionar tecla.
4. Redimensione a **área útil** do terminal para **64 × 18** e **96 × 27** colunas × linhas (`stty size` informa linhas e depois colunas). São aproximadamente quartos de telas 720p e 1080p com células de 10 × 20 px. Confira as 36 posições, o nome da camada e a legenda juntos, sem rolagem ou corte. Abaixo de 64 × 18 deve aparecer o aviso de espaço insuficiente.
