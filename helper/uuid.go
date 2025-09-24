package helper

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
)

func Uuid() string {
	uuidBytes := make([]byte, 16)

	uuidBytes[6] = (6 << 4) | (uuidBytes[6] & 0x0f)

	uuidBytes[8] = (uuidBytes[8] & 0xbf) | 0x80

	rand.Seed(time.Now().UnixNano())
	rand.Read(uuidBytes)

	return uuid.UUID(uuidBytes).String()
}
