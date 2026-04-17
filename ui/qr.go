package ui

import (
	"github.com/luismasuarez/wg-tui/internal/qr"
)

func init() {
	generateQR = func(text string) string {
		out, err := qr.Render(text)
		if err != nil {
			return "Error generando QR: " + err.Error()
		}
		return out
	}
}
