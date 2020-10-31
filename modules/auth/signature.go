package auth

import (
	"crypto/hmac"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"time"

	"github.com/zeebo/blake3"
)

// error
var (
	ErrSasTokenExpired  = errors.New("shared access signature token expired")
	ErrSasTokenNotMatch = errors.New("shared access signature token not match")
)

// Validator interface
type Validator interface {
}

func blake3New() hash.Hash {
	return blake3.New()
}

// MakeSasToken create shared access signature token
func MakeSasToken(secureKey, host, path string, expired int64) (string, error) {
	sectext := fmt.Sprintf("%s://%s\n%d", host, path, expired)
	h := hmac.New(blake3New, []byte(secureKey))
	if _, err := h.Write([]byte(sectext)); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// VerifySasToken todo
func VerifySasToken(sastoken, secureKey, host, path string, expired int64) error {
	if time.Now().Unix() > expired {
		return ErrSasTokenExpired
	}
	st, err := MakeSasToken(secureKey, host, path, expired)
	if err != nil {
		return err
	}
	if st != sastoken {
		return ErrSasTokenNotMatch
	}
	return nil
}
