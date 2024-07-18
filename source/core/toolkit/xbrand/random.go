package xbrand

import (
	"crypto/rand"

	"github.com/gofrs/uuid/v5"
	"github.com/rs/xid"
	"github.com/segmentio/ksuid"

	"github.com/starryck/x-lib-go/source/core/toolkit/xbradix"
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

func MakeBytes(length int) []byte {
	buf := make([]byte, length)
	if _, err := rand.Read(buf); err != nil {
		panic(err)
	} else {
		return buf
	}
}

func MakeBaseNBytes(length int, radixes []byte) []byte {
	src := MakeBytes(length)
	dst := make([]byte, length)
	base := byte(len(radixes))
	for i := 0; i < length; i++ {
		dst[i] = radixes[src[i]%base]
	}
	return dst
}

func MakeBaseNString(length int, radixes []byte) string {
	return string(MakeBaseNBytes(length, radixes))
}

func MakeBase10Bytes(length int) []byte {
	return MakeBaseNBytes(length, xbradix.Base10Bytes)
}

func MakeBase10String(length int) string {
	return string(MakeBase10Bytes(length))
}

func MakeBase16Bytes(length int) []byte {
	return MakeBaseNBytes(length, xbradix.Base16Bytes)
}

func MakeBase16String(length int) string {
	return string(MakeBase16Bytes(length))
}

func MakeBase36Bytes(length int) []byte {
	return MakeBaseNBytes(length, xbradix.Base36Bytes)
}

func MakeBase36String(length int) string {
	return string(MakeBase36Bytes(length))
}

func MakeBase62Bytes(length int) []byte {
	return MakeBaseNBytes(length, xbradix.Base62Bytes)
}

func MakeBase62String(length int) string {
	return string(MakeBase62Bytes(length))
}
