package gbrand

import (
	"crypto/rand"

	"github.com/gofrs/uuid/v5"
	"github.com/rs/xid"
	"github.com/segmentio/ksuid"
)

func MakeXID() string {
	id := xid.New().String()
	return id
}

func MakeKSUID() string {
	id := ksuid.New().String()
	return id
}

func MakeUUID4() string {
	id := uuid.Must(uuid.NewV4()).String()
	return id
}

var alphanumBytes = []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")

func MakeBytes(length int) []byte {
	buf := make([]byte, length)
	if _, err := rand.Read(buf); err != nil {
		panic(err)
	} else {
		return buf
	}
}

func MakeBaseBytes(charset []byte, base byte, length int) []byte {
	src := MakeBytes(length)
	dst := make([]byte, length)
	for i := 0; i < length; i++ {
		dst[i] = charset[src[i]%base]
	}
	return dst
}

func MakeBaseString(charset []byte, base byte, length int) string {
	return string(MakeBaseBytes(charset, base, length))
}

func MakeBase10Bytes(length int) []byte {
	return MakeBaseBytes(alphanumBytes, 10, length)
}

func MakeBase10String(length int) string {
	return string(MakeBase10Bytes(length))
}

func MakeBase16Bytes(length int) []byte {
	return MakeBaseBytes(alphanumBytes, 16, length)
}

func MakeBase16String(length int) string {
	return string(MakeBase16Bytes(length))
}

func MakeBase36Bytes(length int) []byte {
	return MakeBaseBytes(alphanumBytes, 36, length)
}

func MakeBase36String(length int) string {
	return string(MakeBase36Bytes(length))
}

func MakeBase62Bytes(length int) []byte {
	return MakeBaseBytes(alphanumBytes, 62, length)
}

func MakeBase62String(length int) string {
	return string(MakeBase62Bytes(length))
}
