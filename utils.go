package carrot

import (
	"crypto/rand"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
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

func getE_P(currentSession *Session, offset *offset) (*offset, error) {
	var err error
	if currentSession.T_L == nil || currentSession.T_P == nil {
		err = errors.New("The session did not complete the Picnic Protocol handshake")
		return nil, err
	}
	primaryT_P := currentSession.T_P
	log.Infof("t_p: x: %v y: %v z: %v", primaryT_P.X, primaryT_P.Y, primaryT_P.Z)

	currentT_L := currentSession.T_L

	log.Infof("t_l: x: %v y: %v z: %v", currentT_L.X, currentT_L.Y, currentT_L.Z)

	// offset is the e_l
	// o_p = t_l - t_p
	// e_p = e_l - o_p
	o_p := offsetSub(currentT_L, primaryT_P)
	log.Infof("o_p: x: %v y: %v z: %v", o_p.X, o_p.Y, o_p.Z)

	e_p := offsetSub(offset, o_p)

	log.Infof("e_p: x: %v y: %v z: %v", e_p.X, e_p.Y, e_p.Z)
	log.Infof("e_l: x: %v y: %v z: %v", offset.X, offset.Y, offset.Z)

	return e_p, err
}
