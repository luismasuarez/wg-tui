package qr

import (
	qrcode "github.com/skip2/go-qrcode"
)

// Render returns a QR code as a Unicode block string suitable for terminal output.
func Render(text string) (string, error) {
	q, err := qrcode.New(text, qrcode.Medium)
	if err != nil {
		return "", err
	}
	return q.ToSmallString(false), nil
}
