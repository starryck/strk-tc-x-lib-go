package gbjson

import (
	"encoding/json"

	"github.com/bytedance/sonic"
)

var mJSON JSON

type JSON = sonic.API

type RawMessage = json.RawMessage

func Marshal(v any) ([]byte, error) {
	bytes, err := getJSON().Marshal(v)
	return bytes, err
}

func Unmarshal(data []byte, v any) error {
	err := getJSON().Unmarshal(data, v)
	return err
}

func getJSON() JSON {
	if mJSON == nil {
		mJSON = newJSON()
	}
	return mJSON
}

func newJSON() JSON {
	json := sonic.Config{
		CompactMarshaler: true,
	}.Froze()
	return json
}
