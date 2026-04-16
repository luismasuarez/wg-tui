# Tasks: WireGuard TUI Client

**Input**: Design documents from `/specs/001-wireguard-tui-client/`
**Prerequisites**: plan.md ✅, spec.md ✅, research.md ✅, data-model.md ✅, contracts/cli.md ✅

## Format: `[ID] [P?] [Story?] Description`

- **[P]**: Paralelizable (archivos distintos, sin dependencias incompletas)
- **[Story]**: User story asociada (US1–US5)
- Rutas relativas a la raíz del repositorio

---

## Phase 1: Setup

**Purpose**: Inicializar proyecto Go con módulos y estructura de directorios.

- [ ] T001 Inicializar módulo Go: `go mod init github.com/user/wg-tui` y crear estructura de directorios (`internal/nmcli`, `internal/wg`, `internal/qr`, `ui/`) en `go.mod`, `main.go`
- [ ] T002 [P] Agregar dependencias: `go get github.com/charmbracelet/bubbletea github.com/charmbracelet/lipgloss github.com/skip2/go-qrcode golang.zx2c4.com/wireguard/wgctrl` en `go.mod`, `go.sum`
- [ ] T003 [P] Crear `Makefile` con targets: `build`, `install`, `test`, `lint` en `Makefile`

**Checkpoint**: `go build ./...` compila sin errores.

---

## Phase 2: Foundational

**Purpose**: Backend nmcli y tipos core — bloquean todas las user stories.

**⚠️ CRÍTICO**: Las fases siguientes dependen de esta.

- [ ] T004 Definir tipo `Connection` y parsear `nmcli -t -f NAME,TYPE,STATE connection show` en `internal/nmcli/client.go`
- [ ] T005 Implementar `ListConnections() ([]Connection, error)` que filtra por tipo `wireguard` en `internal/nmcli/client.go`
- [ ] T006 [P] Unit tests para parseo de stdout de nmcli (tabla con casos: activa, inactiva, sin resultados, error) en `internal/nmcli/client_test.go`
- [ ] T007 Implementar `Connect(name string) error` y `Disconnect(name string) error` via `nmcli connection up/down` en `internal/nmcli/client.go`
- [ ] T008 [P] Implementar chequeo de precondiciones al arranque: nmcli en PATH y NetworkManager activo, salir con mensaje claro a stderr si falla, en `internal/nmcli/client.go`
- [ ] T009 Implementar modelo Bubble Tea base: `Model` struct con `viewState`, lista de conexiones, mensaje de estado; `Init()`, `Update()`, `View()` esqueleto en `ui/model.go`
- [ ] T010 [P] Definir keybindings (`ui/keys.go`) y barra de ayuda con `lipgloss` en `ui/keys.go`

**Checkpoint**: `go test ./internal/nmcli/...` pasa; `go run main.go` arranca y muestra error claro si no hay nmcli.

---

## Phase 3: User Story 1 — Ver y cambiar conexión activa (P1) 🎯 MVP

**Goal**: Lista de conexiones con estado, conectar con Enter, desconectar con `d`.

**Independent Test**: Con ≥2 perfiles WireGuard en NM, el usuario puede cambiar entre ellos con el teclado.

- [ ] T011 [US1] Implementar vista lista: renderizar conexiones con estado activa/inactiva, indicador de selección, usando `lipgloss` en `ui/list.go`
- [ ] T012 [US1] Conectar teclas `↑↓`/`jk` a navegación en lista; `Enter` dispara `Connect`; `d` dispara `Disconnect` en `ui/model.go`
- [ ] T013 [US1] Mostrar mensaje de estado/error inline en la lista (ej. "Conectando…", "Error: …") en `ui/model.go`, `ui/list.go`
- [ ] T014 [US1] Implementar refresh periódico con `tea.Tick` cada 3s que relanza `ListConnections` en `ui/model.go`
- [ ] T015 [US1] Completar `main.go`: chequear precondiciones, cargar lista inicial, arrancar `tea.NewProgram` en `main.go`

**Checkpoint**: `wg-tui` muestra lista real desde NM; Enter/d conectan y desconectan; estado se actualiza en pantalla.

---

## Phase 4: User Story 2 — Ver detalles de conexión (P2)

**Goal**: Panel de detalles con IP asignada, endpoint, último handshake y tráfico ↑↓.

**Independent Test**: Con una conexión activa, el usuario ve métricas tras pulsar `i`.

- [ ] T016 [US2] Implementar `GetDetails(name string) (Connection, error)` que parsea campos extendidos de nmcli en `internal/nmcli/client.go`
- [ ] T017 [US2] Implementar vista detalle: renderizar todos los campos de `Connection` con formato legible en `ui/detail.go`
- [ ] T018 [US2] Conectar tecla `i` para cambiar a `viewDetail`; `Esc` para volver a `viewList` en `ui/model.go`

**Checkpoint**: Pulsar `i` sobre cualquier conexión muestra el panel de detalles.

---

## Phase 5: User Story 3 — Generar QR (P3)

**Goal**: QR code en terminal con la configuración del peer en formato WireGuard estándar.

**Independent Test**: El QR generado puede escanearse con la app oficial de WireGuard.

- [ ] T019 [US3] Implementar `BuildConfig(c Connection) string` que genera el texto `.conf` WireGuard estándar en `internal/qr/render.go`
- [ ] T020 [US3] Implementar `RenderQR(config string) string` usando `skip2/go-qrcode` con `ToSmallString` en `internal/qr/render.go`
- [ ] T021 [US3] Implementar vista QR: mostrar el QR centrado con instrucción "Esc para volver" en `ui/qr.go`
- [ ] T022 [US3] Conectar tecla `q` (en lista y detalle) para cambiar a `viewQR`; `Esc` para volver en `ui/model.go`

**Checkpoint**: Pulsar `q` muestra el QR; escaneable con WireGuard para iOS/Android.

---

## Phase 6: User Story 4 — Agregar peer (P4)

**Goal**: Formulario para crear nuevo perfil WireGuard con keypair generado automáticamente.

**Independent Test**: Se puede crear un perfil nuevo en NM desde la TUI.

- [ ] T023 [US4] Implementar `GenerateKeypair() (priv, pub string, err error)` usando `wgctrl/wgtypes` en `internal/wg/keygen.go`
- [ ] T024 [US4] Implementar `AddConnection(form NewPeerForm, privKey string) error` via `nmcli connection add` en `internal/nmcli/client.go`
- [ ] T025 [US4] Implementar vista formulario: campos Name, Endpoint, ServerPublicKey, AssignedIP, PresharedKey (opcional), validación inline en `ui/form.go`
- [ ] T026 [US4] Conectar tecla `a` para abrir formulario; `Tab` navega campos; `Enter` en último campo confirma; `Esc` cancela en `ui/model.go`
- [ ] T027 [US4] Al confirmar: generar keypair, llamar `AddConnection`, mostrar clave pública generada para compartir con el servidor en `ui/model.go`

**Checkpoint**: Desde la TUI se puede crear un perfil que aparece en `nmcli connection show`.

---

## Phase 7: User Story 5 — Eliminar peer (P5)

**Goal**: Eliminar un perfil de NM con confirmación previa.

**Independent Test**: Un perfil existente puede eliminarse y desaparece de la lista.

- [ ] T028 [US5] Implementar `DeleteConnection(name string) error` via `nmcli connection delete` en `internal/nmcli/client.go`
- [ ] T029 [US5] Implementar vista confirmación: mostrar "¿Eliminar <nombre>? [y/n]" en `ui/model.go`
- [ ] T030 [US5] Conectar tecla `x` para pedir confirmación; `y` confirma y ejecuta delete + refresca lista; `n`/`Esc` cancela en `ui/model.go`
- [ ] T031 [US5] Si la conexión está activa, desconectar antes de eliminar en `ui/model.go`

**Checkpoint**: Pulsar `x` + `y` elimina el perfil; ya no aparece en la lista.

---

## Phase 8: Polish

**Purpose**: Distribución, ayuda y refinamiento final.

- [ ] T032 [P] Crear `install.sh`: descarga binario del release de GitHub y lo coloca en `~/.local/bin/wg-tui` en `install.sh`
- [ ] T033 [P] Agregar `--version` flag con versión embebida via `-ldflags` en `main.go`
- [ ] T034 [P] Implementar vista de ayuda (`?`) que muestra tabla de keybindings en `ui/model.go`, `ui/keys.go`
- [ ] T035 Manejar redimensionado de terminal (`tea.WindowSizeMsg`) y mínimo 80×24 en `ui/model.go`
- [ ] T036 Compilar binario estático release: `CGO_ENABLED=0 go build -ldflags="-s -w"` y verificar tamaño <10MB en `Makefile`

---

## Dependencies & Execution Order

### Orden de fases

- **Phase 1 (Setup)**: Sin dependencias — empezar aquí.
- **Phase 2 (Foundational)**: Depende de Phase 1 — **bloquea todas las user stories**.
- **Phases 3–7 (User Stories)**: Todas dependen de Phase 2. Pueden avanzar en orden P1→P2→P3→P4→P5.
- **Phase 8 (Polish)**: Depende de que las user stories deseadas estén completas.

### Dependencias entre user stories

- **US1 (P1)**: Sin dependencias en otras stories.
- **US2 (P2)**: Sin dependencias — usa `GetDetails` que es extensión de nmcli/client.go.
- **US3 (P3)**: Sin dependencias en otras stories.
- **US4 (P4)**: Sin dependencias en otras stories.
- **US5 (P5)**: Sin dependencias en otras stories.

---

## Implementation Strategy

### MVP (solo US1)

1. Phase 1: Setup
2. Phase 2: Foundational
3. Phase 3: US1 — listar y cambiar conexión
4. **VALIDAR**: `wg-tui` funciona para cambiar entre conexiones WireGuard.

### Entrega incremental

Cada fase añade valor sin romper lo anterior:
- MVP: US1 — reemplaza el applet gráfico ✅
- +US2 — diagnóstico sin salir de la TUI
- +US3 — exportar config al móvil con QR
- +US4 — onboarding de nuevas conexiones
- +US5 — mantenimiento de perfiles

---

## Resumen

| Fase | Tasks | Paralelizables |
|------|-------|----------------|
| Setup | T001–T003 | T002, T003 |
| Foundational | T004–T010 | T006, T008, T010 |
| US1 (P1) | T011–T015 | — |
| US2 (P2) | T016–T018 | — |
| US3 (P3) | T019–T022 | — |
| US4 (P4) | T023–T027 | — |
| US5 (P5) | T028–T031 | — |
| Polish | T032–T036 | T032, T033, T034 |
| **Total** | **36 tasks** | **8 paralelizables** |
