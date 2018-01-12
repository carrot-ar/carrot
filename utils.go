package carrot

import (
	"crypto/rand"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"strings"
)


func validSession(serverToken SessionToken, clientToken SessionToken) error {

	if serverToken != clientToken {
		return fmt.Errorf("client-server token mismatch")
	}

	return nil
}

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

func getE_P(sender *Session, recipient *Session, event *offset) (*offset, error) {
	var err error
	var T_F *offset //transform foreign
	var T_L *offset //transform local
	if !sender.isPrimaryDevice() {
		if sender.T_L == nil || sender.T_P == nil {
			err = errors.New("The session did not complete the Picnic Protocol handshake")
			return nil, err
		}
		T_F = sender.T_P
		T_L = sender.T_L
	} else { //the sender is a primary device
		T_F = recipient.T_P
		T_L = recipient.T_L
	}

	log.Infof("t_p: x: %v y: %v z: %v", T_F.X, T_F.Y, T_F.Z)
	log.Infof("t_l: x: %v y: %v z: %v", T_L.X, T_L.Y, T_L.Z)

	// event is the e_l, event local
	// o_p = t_l - t_p, foreign recipient's origin placement in this coordinate system
	// e_p = e_l - o_p, event placement
	O_P := offsetSub(T_L, T_F)
	log.Infof("o_p: x: %v y: %v z: %v", O_P.X, O_P.Y, O_P.Z)

	E_P := offsetSub(event, O_P)
	log.Infof("e_l: x: %v y: %v z: %v", event.X, event.Y, event.Z)
	log.Infof("e_p: x: %v y: %v z: %v", E_P.X, E_P.Y, E_P.Z)

	return E_P, err
}
