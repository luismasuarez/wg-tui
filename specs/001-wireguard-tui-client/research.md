# Research: WireGuard TUI Client

## 1. nmcli — Parsing de conexiones WireGuard

**Decision**: Usar `nmcli -t -f NAME,TYPE,STATE connection show` para listar y `nmcli -t -f` para detalles.

**Formato terse (`-t`)**: campos separados por `:`, filas por `\n`. Predecible y fácil de parsear sin regex complejos.

**Comandos clave**:
```bash
# Listar conexiones WireGuard
nmcli -t -f NAME,TYPE,STATE connection show | grep wireguard

# Detalles de una conexión (incluye IP, endpoint)
nmcli -t connection show "<nombre>"

# Estado en tiempo real (handshake, tráfico) — requiere conexión activa
nmcli -t -f GENERAL,IP4,IP6 device show "<iface>"

# Conectar / desconectar
nmcli connection up "<nombre>"
nmcli connection down "<nombre>"

# Agregar conexión WireGuard
nmcli connection add type wireguard \
  con-name "<nombre>" \
  wireguard.private-key "<privkey>" \
  wireguard.peers "[{public-key=<pubkey>,endpoint=<host:port>,allowed-ips=<cidr>}]" \
  ipv4.method manual \
  ipv4.addresses "<ip/prefix>"

# Eliminar
nmcli connection delete "<nombre>"
```

**Rationale**: nmcli es el único backend requerido por la constitución. No hay alternativa aceptable (wg-quick rompería el principio V).

**Alternativas descartadas**: wg-tools (requiere root siempre), netlink directo (complejo, rompe principio V).

---

## 2. Generación de keypair WireGuard en Go

**Decision**: `golang.zx2c4.com/wireguard/wgctrl/wgtypes` para generar keypairs.

```go
import "golang.zx2c4.com/wireguard/wgctrl/wgtypes"

key, _ := wgtypes.GeneratePrivateKey()
privKey := key.String()   // base64
pubKey  := key.PublicKey().String()
```

**Rationale**: Librería oficial del proyecto WireGuard para Go. Pura Go, sin CGO. Genera claves Curve25519 correctas.

**Alternativas descartadas**: `os/exec wg genkey` (requiere wg-tools instalado, rompe principio V).

---

## 3. QR Code en terminal

**Decision**: `github.com/skip2/go-qrcode` — renderiza QR como bloques Unicode directamente en stdout.

```go
import qrcode "github.com/skip2/go-qrcode"

qr, _ := qrcode.New(configStr, qrcode.Medium)
fmt.Println(qr.ToSmallString(false))
```

**Rationale**: Puro Go, sin CGO, sin binarios externos. `ToSmallString` usa caracteres Unicode `▀▄█` que funcionan en cualquier terminal moderno.

**Alternativas descartadas**: shell a `qrencode` (rompe principio II — dependencia externa); `github.com/mdp/qrterminal` (similar pero menos mantenido).

---

## 4. Bubble Tea — Arquitectura de vistas

**Decision**: Modelo único con campo `view` (enum de vista activa). Cada vista es un componente que implementa su propio `View() string`.

```go
type viewState int
const (
    viewList viewState = iota
    viewDetail
    viewQR
    viewAddForm
    viewConfirmDelete
)
```

**Rationale**: La app tiene flujo lineal (lista → detalle/acción). Un modelo único evita complejidad de sub-modelos y mensajes entre componentes. Suficiente para la escala del proyecto.

**Alternativas descartadas**: múltiples programas Bubble Tea anidados (overhead innecesario para 5 vistas).

---

## 5. Refresh periódico del estado

**Decision**: `tea.Tick` con intervalo de 3 segundos para refrescar el estado de conexiones.

**Rationale**: Balance entre responsividad y no saturar nmcli. 3s es imperceptible para el usuario pero mantiene el estado fresco.

---

## 6. Distribución — Binario estático + install.sh

**Decision**: `CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w"` → binario ~5-8MB.

`install.sh`: descarga el binario del release de GitHub, lo coloca en `~/.local/bin/wg-tui`.

**Rationale**: Sin CGO el binario no depende de libc del sistema. `-s -w` elimina símbolos de debug para reducir tamaño. Compatible con cualquier Linux x86_64.
