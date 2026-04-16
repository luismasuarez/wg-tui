# Feature Specification: WireGuard TUI Client

**Feature Branch**: `001-wireguard-tui-client`
**Created**: 2026-04-16
**Status**: Draft
**Input**: TUI en Go para gestionar conexiones WireGuard en Ubuntu como cliente. Usa nmcli como backend. Lista las conexiones disponibles con su estado (activa/inactiva), permite conectar y desconectar con un keypress, muestra detalles de cada conexión, permite agregar/eliminar peers, y genera QR codes en terminal para exportar configuraciones a móvil.

---

## User Scenarios & Testing

### User Story 1 — Ver y cambiar conexión activa (Priority: P1)

El usuario abre `wg-tui` en la terminal y ve inmediatamente la lista de conexiones WireGuard con su estado (activa/inactiva). Selecciona una conexión inactiva y pulsa `Enter` para activarla. La conexión anterior se desactiva automáticamente si corresponde. Con `d` desconecta la activa.

**Why this priority**: Es el flujo central que reemplaza el applet gráfico. Sin esto no hay herramienta.

**Independent Test**: Con al menos dos perfiles WireGuard en NetworkManager, el usuario puede cambiar entre ellos usando solo el teclado.

**Acceptance Scenarios**:

1. **Given** la aplicación está abierta, **When** el usuario navega con ↑↓ y pulsa `Enter`, **Then** la conexión seleccionada se activa y el estado se actualiza en pantalla.
2. **Given** hay una conexión activa, **When** el usuario pulsa `d`, **Then** la conexión se desactiva y el estado cambia a inactiva.
3. **Given** NetworkManager no tiene perfiles WireGuard, **When** se abre la TUI, **Then** se muestra un mensaje claro indicando que no hay conexiones disponibles.

---

### User Story 2 — Ver detalles de una conexión (Priority: P2)

El usuario selecciona una conexión y pulsa `Enter` o una tecla dedicada para ver sus detalles: dirección IP asignada, endpoint del servidor, último handshake, tráfico enviado/recibido.

**Why this priority**: Necesario para diagnóstico básico sin salir de la TUI.

**Independent Test**: Con una conexión activa, el usuario puede ver métricas de tráfico y tiempo del último handshake.

**Acceptance Scenarios**:

1. **Given** una conexión activa seleccionada, **When** el usuario pulsa `i`, **Then** se muestra un panel con IP, endpoint, último handshake y tráfico ↑↓.
2. **Given** una conexión inactiva seleccionada, **When** el usuario pulsa `i`, **Then** se muestran los datos de configuración disponibles (endpoint, IP configurada) sin métricas de tráfico.

---

### User Story 3 — Generar QR code para exportar a móvil (Priority: P3)

El usuario selecciona una conexión y genera un QR code en la terminal que contiene la configuración del peer para importarla en la app WireGuard del móvil.

**Why this priority**: Funcionalidad de conveniencia de alto valor para el flujo de trabajo del usuario.

**Independent Test**: Dado un perfil WireGuard existente, se muestra un QR scannable en la terminal.

**Acceptance Scenarios**:

1. **Given** una conexión seleccionada, **When** el usuario pulsa `q`, **Then** se renderiza un QR code en la terminal con la configuración del perfil en formato WireGuard estándar.
2. **Given** el QR generado, **When** se escanea con la app WireGuard en iOS/Android, **Then** la configuración se importa correctamente.

---

### User Story 4 — Agregar peer (Priority: P4)

El usuario puede agregar un nuevo peer WireGuard. La herramienta genera automáticamente el keypair (privada/pública) y presenta los datos para configurar el servidor.

**Why this priority**: Cubre el flujo de onboarding de una nueva conexión.

**Independent Test**: Se puede crear un perfil WireGuard nuevo en NetworkManager desde la TUI con un keypair generado automáticamente.

**Acceptance Scenarios**:

1. **Given** el usuario pulsa `a`, **When** completa el formulario (nombre, endpoint, IP asignada, clave pública del servidor, preshared key opcional), **Then** se crea el perfil en NetworkManager con un keypair generado.
2. **Given** el formulario incompleto, **When** el usuario intenta confirmar, **Then** se indica el campo faltante sin abandonar el formulario.

---

### User Story 5 — Eliminar peer (Priority: P5)

El usuario selecciona una conexión y la elimina de NetworkManager con confirmación previa.

**Why this priority**: Mantenimiento básico de perfiles.

**Independent Test**: Un perfil existente puede eliminarse y ya no aparece en la lista.

**Acceptance Scenarios**:

1. **Given** una conexión seleccionada e inactiva, **When** el usuario pulsa `x` y confirma, **Then** el perfil se elimina de NetworkManager.
2. **Given** una conexión activa seleccionada, **When** el usuario pulsa `x`, **Then** se pide confirmación explícita antes de desconectar y eliminar.

---

### Edge Cases

- ¿Qué pasa si `nmcli` no está instalado o NetworkManager no está corriendo? → Mensaje de error claro y salida limpia.
- ¿Qué pasa si una operación de conexión falla (ej. servidor inalcanzable)? → Se muestra el error de nmcli al usuario.
- ¿Qué pasa si se pierde el acceso a NetworkManager mientras la TUI está abierta? → La TUI detecta el error en el siguiente refresh y lo notifica.
- ¿Qué pasa si el terminal es demasiado pequeño? → Se muestra un mensaje mínimo indicando el tamaño mínimo requerido.

---

## Requirements

### Functional Requirements

- **FR-001**: La herramienta DEBE listar todas las conexiones WireGuard gestionadas por NetworkManager con su estado (activa/inactiva).
- **FR-002**: La herramienta DEBE permitir activar una conexión seleccionada con un keypress.
- **FR-003**: La herramienta DEBE permitir desactivar la conexión activa con un keypress.
- **FR-004**: La herramienta DEBE mostrar detalles de una conexión seleccionada: IP asignada, endpoint, último handshake, tráfico ↑↓.
- **FR-005**: La herramienta DEBE generar un QR code en la terminal con la configuración de un peer en formato WireGuard estándar.
- **FR-006**: La herramienta DEBE permitir agregar un nuevo peer generando el keypair automáticamente.
- **FR-007**: La herramienta DEBE permitir eliminar un peer con confirmación previa.
- **FR-008**: La herramienta DEBE salir con un mensaje de error claro si `nmcli` no está disponible o NetworkManager no responde.
- **FR-009**: La herramienta DEBE refrescar el estado de las conexiones periódicamente sin requerir acción del usuario.
- **FR-010**: Todas las operaciones de escritura (conectar, desconectar, agregar, eliminar) DEBEN reportar éxito o error visible al usuario.

### Key Entities

- **Conexión WireGuard**: Perfil gestionado por NetworkManager. Atributos: nombre, estado (activa/inactiva), IP asignada, endpoint del servidor, último handshake, tráfico ↑↓, clave pública propia.
- **Peer**: Extremo remoto de una conexión WireGuard. Atributos: clave pública, clave preshared (opcional), endpoint, allowed IPs.

---

## Success Criteria

### Measurable Outcomes

- **SC-001**: El usuario puede cambiar de conexión WireGuard activa en menos de 3 pulsaciones de tecla.
- **SC-002**: La lista de conexiones con estado se muestra en menos de 1 segundo tras abrir la herramienta.
- **SC-003**: El usuario puede completar el flujo de agregar un nuevo peer sin consultar documentación externa.
- **SC-004**: Un QR generado puede ser escaneado e importado correctamente en la app WireGuard oficial de iOS y Android.
- **SC-005**: La herramienta no requiere configuración previa para funcionar en cualquier sistema Linux con NetworkManager y perfiles WireGuard existentes.

---

## Assumptions

- El sistema tiene NetworkManager instalado y corriendo con al menos un perfil WireGuard configurado.
- El usuario tiene permisos para ejecutar `nmcli connection up/down` (puede requerir que el usuario sea propietario de la conexión en NetworkManager o pertenecer al grupo `networkmanager`).
- El terminal soporta al menos 80×24 caracteres.
- La herramienta es para uso personal de un único usuario; no hay gestión de múltiples usuarios ni permisos entre usuarios.
- El formato de QR exportado es el archivo `.conf` estándar de WireGuard (compatible con las apps oficiales).
- Soporte únicamente para conexiones WireGuard; otras conexiones VPN de NetworkManager se ignoran.
