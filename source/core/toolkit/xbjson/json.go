package xbjson

import (
	"encoding/json"

	"github.com/bytedance/sonic"
)

var mJSON JSON

type (
	JSON       = sonic.API
	RawMessage = json.RawMessage
)

func getJSON() JSON {
	if mJSON == nil {
		mJSON = newJSON()
	}
	return mJSON
}

func newJSON() JSON {
	json := sonic.Config{
		SortMapKeys:      true,
		CompactMarshaler: true,
	}.Froze()
	return json
}

func Marshal(value any) ([]byte, error) {
	data, err := getJSON().Marshal(value)
	return data, err
}

func Unmarshal(data []byte, value any) error {
	err := getJSON().Unmarshal(data, value)
	return err
}
