package gbidtf

import (
	"github.com/gofrs/uuid"
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
