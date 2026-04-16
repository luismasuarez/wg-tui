# Data Model: WireGuard TUI Client

> Estado siempre leído desde NetworkManager. No hay persistencia propia.

## Connection

Representa un perfil WireGuard en NetworkManager.

| Campo | Tipo | Fuente nmcli | Notas |
|-------|------|-------------|-------|
| Name | string | `NAME` | Identificador único en NM |
| Active | bool | `STATE` == "activated" | |
| Interface | string | `GENERAL.IP-IFACE` | Nombre de iface cuando activa |
| AssignedIP | string | `IP4.ADDRESS` | CIDR, ej. `10.0.0.2/24` |
| Endpoint | string | `wireguard.peers[0].endpoint` | `host:port` del servidor |
| PublicKey | string | `wireguard.public-key` | Clave pública propia |
| PeerPubKey | string | `wireguard.peers[0].public-key` | Clave pública del servidor |
| LastHandshake | time.Time | via `wg show` o NM si disponible | Solo cuando activa |
| RxBytes | int64 | `wireguard.peers[0].rx-bytes` | Solo cuando activa |
| TxBytes | int64 | `wireguard.peers[0].tx-bytes` | Solo cuando activa |

## NewPeerForm

Datos para crear una nueva conexión WireGuard.

| Campo | Tipo | Requerido | Notas |
|-------|------|-----------|-------|
| Name | string | ✅ | Nombre del perfil en NM |
| Endpoint | string | ✅ | `host:port` |
| ServerPublicKey | string | ✅ | Clave pública del servidor |
| AssignedIP | string | ✅ | CIDR a asignar al cliente |
| PresharedKey | string | ❌ | Opcional |
| DNS | string | ❌ | Opcional |

El keypair del cliente se genera automáticamente al crear.

## State Transitions

```
inactiva ──[Enter]──→ activando ──[ok]──→ activa
activa   ──[d]────→ desconectando ──[ok]──→ inactiva
* ────────[error nmcli]──────────────────→ estado anterior + mensaje error
```
