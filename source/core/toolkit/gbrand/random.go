package gbrand

import (
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
