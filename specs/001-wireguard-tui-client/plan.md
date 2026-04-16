# Implementation Plan: WireGuard TUI Client

**Branch**: `001-wireguard-tui-client` | **Date**: 2026-04-16 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/001-wireguard-tui-client/spec.md`

## Summary

TUI en Go para gestionar conexiones WireGuard a través de NetworkManager (`nmcli`).
Interfaz de teclado con Bubble Tea. Binario único estático sin dependencias de runtime.
Cubre: listar/conectar/desconectar conexiones, ver detalles, generar QR, agregar/eliminar peers.

## Technical Context

**Language/Version**: Go 1.21+
**Primary Dependencies**: charmbracelet/bubbletea, charmbracelet/lipgloss, skip2/go-qrcode
**Storage**: N/A — estado siempre leído desde NetworkManager
**Testing**: `go test` — unitarios para parsing de nmcli; integración con build tag `integration`
**Target Platform**: Linux con NetworkManager ≥ 1.2
**Project Type**: CLI / TUI
**Performance Goals**: Lista visible en <1s; sin polling agresivo (refresh cada 3s)
**Constraints**: Binario estático <10MB; sin archivos de config propios; sin root para lecturas
**Scale/Scope**: Uso personal, 1 usuario, ~10 conexiones máximo esperadas

## Constitution Check

| Principio | Estado | Notas |
|-----------|--------|-------|
| I. Go Idiomático | ✅ | Interfaces pequeñas, errores explícitos, sin reflexión |
| II. Binario Único | ✅ | `CGO_ENABLED=0 go build` → binario estático; QR puro en Go |
| III. UX Minimalista | ✅ | Solo teclado, cero config, arranque <1s |
| IV. Privilegios Mínimos | ✅ | Leer no requiere sudo; conectar/desconectar delega a nmcli (permisos NM) |
| V. Compatibilidad Linux+NM | ✅ | Solo nmcli, sin wg-tools, sin APIs Ubuntu/GNOME |

**GATE**: ✅ Sin violaciones. Procede a diseño.

## Project Structure

### Documentation (esta feature)

```text
specs/001-wireguard-tui-client/
├── plan.md          ← este archivo
├── research.md      ← Phase 0
├── data-model.md    ← Phase 1
├── contracts/       ← Phase 1 (CLI contract)
│   └── cli.md
└── tasks.md         ← generado por /speckit.tasks
```

### Source Code

```text
wg-tui/
├── main.go                  # entrypoint: inicializa y arranca bubbletea
├── cmd/
│   └── root.go              # cobra o flag parsing mínimo (flags: --version)
├── internal/
│   ├── nmcli/
│   │   ├── client.go        # exec nmcli, parsea stdout
│   │   └── client_test.go
│   ├── wg/
│   │   ├── keygen.go        # genera keypair WireGuard (golang.zx2c4.com/wireguard/wgctrl)
│   │   └── keygen_test.go
│   └── qr/
│       └── render.go        # genera QR en terminal (skip2/go-qrcode)
├── ui/
│   ├── model.go             # Bubble Tea model principal
│   ├── list.go              # vista: lista de conexiones
│   ├── detail.go            # vista: detalles de conexión
│   ├── form.go              # vista: formulario agregar peer
│   ├── qr.go                # vista: QR code
│   └── keys.go              # keybindings
├── go.mod
├── go.sum
├── Makefile
└── install.sh               # curl | sh installer
```

**Structure Decision**: Proyecto único. `internal/` aísla el backend (nmcli, keygen, qr) de la UI. La UI en `ui/` sigue el patrón Model-Update-View de Bubble Tea. Sin subdirectorios de features — la app es pequeña y cohesiva.

---

## Phase 0: Research

Ver [research.md](./research.md)

---

## Phase 1: Design

Ver [data-model.md](./data-model.md) y [contracts/cli.md](./contracts/cli.md)
