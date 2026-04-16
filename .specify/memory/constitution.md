<!-- Sync Impact Report
Version change: 0.0.0 → 1.0.0
Added sections: Core Principles, Constraints, Governance
Modified principles: (initial — no prior version)
Templates requiring updates:
  ✅ constitution.md populated
  ✅ No placeholder tokens remaining
Follow-up TODOs: none
-->

# wg-tui Constitution

## Core Principles

### I. Go Idiomático
El código DEBE seguir las convenciones estándar de Go: interfaces pequeñas, composición
sobre herencia, manejo explícito de errores, sin magia ni reflexión innecesaria.
Dependencias externas al sistema operativo DEBEN ser mínimas y justificadas.

### II. Binario Único
El artefacto final DEBE ser un binario estático compilado sin dependencias de runtime
externas al sistema. Toda lógica se compila dentro del binario. La distribución se
realiza vía un script `curl | sh` que descarga y coloca el binario en `$PATH`.

### III. UX Minimalista
La interfaz DEBE ser directa y operable con el teclado únicamente. Cero configuración
requerida para el primer uso. El tiempo desde invocación hasta primera acción DEBE
ser inferior a 1 segundo en hardware normal.

### IV. Privilegios Mínimos
El binario NO DEBE requerir `sudo` para operaciones de lectura (listar, ver estado).
`sudo` o permisos elevados SOLO se permiten cuando NetworkManager lo exija
explícitamente para conectar/desconectar. Se documenta cada caso que lo requiera.

### V. Compatibilidad Linux + NetworkManager
El backend DEBE usar exclusivamente `nmcli` como interfaz con NetworkManager.
Compatible con cualquier distribución Linux que tenga NetworkManager ≥ 1.2.
Sin dependencias de APIs específicas de Ubuntu o GNOME.

## Constraints

- Lenguaje: Go 1.21+
- TUI framework: Bubble Tea (charmbracelet/bubbletea)
- Backend: `nmcli` (no wg-tools directamente, no netlink)
- QR codes: librería Go pura (no shell a `qrencode`)
- Tests: unitarios para lógica de parsing de `nmcli`; integración marcada con build tag `integration`
- Sin archivos de configuración propios; estado leído siempre desde NetworkManager

## Governance

Esta constitución rige todas las decisiones de diseño e implementación del proyecto.
Cualquier desviación DEBE documentarse con justificación explícita en el PR correspondiente.
Enmiendas requieren actualizar este archivo con incremento de versión semántica:
- MAJOR: remoción o redefinición incompatible de un principio
- MINOR: adición de principio o sección
- PATCH: clarificaciones y correcciones de redacción

**Version**: 1.0.0 | **Ratified**: 2026-04-16 | **Last Amended**: 2026-04-16
