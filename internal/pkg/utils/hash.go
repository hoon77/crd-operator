package utils

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
)

func HashMapString(data map[string]string) string {
	// data marshal -> []byte
	raw, _ := json.Marshal(data)
	// SHA-1 hash
	sum := sha1.Sum(raw)
	return hex.EncodeToString(sum[:])
}
