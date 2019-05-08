package tool

import (
	"crypto/rand"
	"encoding/base64"
	"io"

	"github.com/champly/lib4go/security/md5"
)

func GetGUID() string {
	b := make([]byte, 48)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""

	}
	return md5.Encrypt(base64.URLEncoding.EncodeToString(b))
}
