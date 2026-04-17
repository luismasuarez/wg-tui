package wg

import "golang.zx2c4.com/wireguard/wgctrl/wgtypes"

// GenerateKeypair returns a base64-encoded private and public key pair.
func GenerateKeypair() (priv, pub string, err error) {
	key, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		return "", "", err
	}
	return key.String(), key.PublicKey().String(), nil
}
