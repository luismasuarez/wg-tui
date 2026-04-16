# CLI Contract: wg-tui

## Invocación

```
wg-tui [--version] [--help]
```

Sin argumentos adicionales — toda la interacción es dentro de la TUI.

## Keybindings

| Tecla | Contexto | Acción |
|-------|----------|--------|
| `↑` / `k` | Lista | Mover selección arriba |
| `↓` / `j` | Lista | Mover selección abajo |
| `Enter` | Lista, conexión inactiva | Conectar |
| `d` | Lista, conexión activa | Desconectar |
| `i` | Lista | Ver detalles de conexión seleccionada |
| `q` | Lista, Detalle | Generar QR de la conexión seleccionada |
| `a` | Lista | Abrir formulario agregar peer |
| `x` | Lista | Eliminar conexión seleccionada (pide confirmación) |
| `Esc` | Detalle, QR, Formulario | Volver a lista |
| `Tab` | Formulario | Siguiente campo |
| `Enter` | Formulario (último campo) | Confirmar creación |
| `y` / `n` | Confirmación eliminar | Confirmar / cancelar |
| `?` | Cualquiera | Mostrar/ocultar ayuda de teclas |
| `Ctrl+C` / `q` | Lista | Salir |

## Exit Codes

| Código | Condición |
|--------|-----------|
| 0 | Salida normal |
| 1 | Error fatal (nmcli no disponible, NetworkManager no responde) |

## Stderr

Errores fatales al arranque se escriben en stderr antes de salir.
Errores de operaciones nmcli se muestran dentro de la TUI como mensajes de estado.

## Requisitos del entorno

- `nmcli` en `$PATH`
- NetworkManager corriendo (`systemctl is-active NetworkManager`)
- Terminal con soporte UTF-8 y mínimo 80×24
