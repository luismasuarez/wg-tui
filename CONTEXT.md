# wg-tui — WireGuard TUI Client

## Contexto del proyecto

TUI para gestionar conexiones WireGuard como cliente en Ubuntu, desde la terminal.

## Problema

En Ubuntu, cambiar entre conexiones WireGuard requiere usar el applet gráfico de la barra de tareas. El objetivo es reemplazarlo con una TUI rápida desde la terminal.

## Stack técnico decidido

- **Lenguaje:** Go
- **TUI framework:** Bubble Tea (Charm)
- **Backend:** `nmcli` (NetworkManager CLI)
- **Distribución:** binario único, instalable via `curl | sh`

## Conexiones existentes (en este equipo)

| Nombre       | Estado   |
|--------------|----------|
| RouterOS     | Activa   |
| Internacional| Inactiva |
| Nacional     | Inactiva |

Gestionadas por NetworkManager. Comandos clave:
```bash
nmcli connection show | grep wireguard
nmcli connection up <nombre>
nmcli connection down <nombre>
```

## Funcionalidades requeridas (v1)

1. Listar conexiones WireGuard con estado (activa/inactiva, último handshake, tráfico ↑↓)
2. Conectar / desconectar con un keypress
3. Ver configuración de un peer
4. Agregar peer (genera keypair automáticamente)
5. Eliminar peer
6. Generar QR code en terminal para exportar a móvil

## Flujo de desarrollo — Spec-Driven (spec-kit)

Seguir este orden estrictamente:

### 0. Setup
```bash
uv tool install specify-cli --from git+https://github.com/github/spec-kit.git
specify init . --ai kiro
```

### 1. Constitution
```
/speckit.constitution
Principios: Go idiomático, binario único sin dependencias externas al sistema,
UX minimalista, requiere sudo solo cuando es estrictamente necesario,
compatibilidad con cualquier distro Linux con NetworkManager.
```

### 2. Specify
```
/speckit.specify
TUI en Go para gestionar conexiones WireGuard en Ubuntu como cliente.
Usa nmcli como backend. Lista las conexiones disponibles con su estado
(activa/inactiva), permite conectar y desconectar con un keypress,
muestra detalles de cada conexión, permite agregar/eliminar peers,
y genera QR codes en terminal para exportar configuraciones a móvil.
```

### 3. Plan
```
/speckit.plan
Go + Bubble Tea para la TUI. nmcli para interactuar con NetworkManager.
qrencode o librería Go pura para QR codes. Compilar a binario único.
Script de instalación via curl para distribución.
```

### 4. Tasks
```
/speckit.tasks
```

### 5. Implement
```
/speckit.implement
```
