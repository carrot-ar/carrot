package carrot

import (
	"crypto/rand"
	"fmt"
	"io"
	"strings"
)

func InSlice(str string, items []string) bool {
	for _, item := range items {
		if item == str {
			return true
		}
	}

	return false
}

// generate UUID fulfilling RFC 4122
func generateUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}

	// variant bits
	uuid[8] = uuid[8]&^0xc0 | 0x80

	// version 4
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return strings.ToUpper(fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])), nil
}

// a - b
func offsetSub(a *offset, b *offset) *offset {
	return &offset{
		X: a.X - b.X,
		Y: a.Y - b.Y,
		Z: a.Z - b.Z,
	}
}
