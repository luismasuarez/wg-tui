# wg-tui

TUI para gestionar conexiones WireGuard desde la terminal, usando NetworkManager como backend.

![demo](https://github.com/luismasuarez/wg-tui/assets/demo.gif)

## Requisitos

- Linux con NetworkManager
- `nmcli` en PATH
- `wg` en PATH (para ver detalles de conexiones activas)

## Instalación

```bash
curl -fsSL https://raw.githubusercontent.com/luismasuarez/wg-tui/main/install.sh | sh
sudo cp ~/.local/bin/wg-tui /usr/local/bin/wg-tui
```

## Uso

```bash
sudo wg-tui
```

## Teclas

| Tecla | Acción |
|-------|--------|
| `↑` / `k` | Arriba |
| `↓` / `j` | Abajo |
| `Enter` | Conectar |
| `d` | Desconectar |
| `i` | Ver detalles |
| `q` | Generar QR |
| `a` | Agregar peer |
| `x` | Eliminar peer |
| `r` | Refrescar |
| `?` | Ayuda |
| `Ctrl+C` | Salir |

## Compilar desde fuente

```bash
git clone git@github.com:luismasuarez/wg-tui.git
cd wg-tui
make build
sudo cp wg-tui /usr/local/bin/wg-tui
```
