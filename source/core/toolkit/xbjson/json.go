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

func Marshal(v any) ([]byte, error) {
	bytes, err := getJSON().Marshal(v)
	return bytes, err
}

func Unmarshal(data []byte, v any) error {
	err := getJSON().Unmarshal(data, v)
	return err
}
