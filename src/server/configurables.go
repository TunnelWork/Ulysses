package server

import (
	"encoding/json"
	"errors"
)

// JSON utils

var (
	ErrNotJsonObject = errors.New("ERR_NOT_JSON_OBJ")
	ErrNotJsonArray  = errors.New("ERR_NOT_JSON_ARRAY")
)

// func isJSONObj(data []byte) bool {
// 	var js map[string]interface{}
// 	return json.Unmarshal(data, &js) == nil
// }

// func isJSONArray(data []byte) bool {
// 	var js []map[string]interface{}
// 	return json.Unmarshal(data, &js) == nil
// }

type ServerConfigurables map[string]string
type AccountConfigurables map[string]string

// JsonToConfigurables returns map[string]string. It is caller's job to cast into either Server or Account configurables.
func JsonToConfigurables(data []byte) (map[string]string, error) {
	var configurables map[string]string
	if json.Unmarshal(data, &configurables) != nil {
		return nil, ErrNotJsonObject
	}
	return configurables, nil
}

func JsonArrToConfigurablesSlice(data []byte) ([]map[string]string, error) {
	var configurablesSlice []map[string]string
	if json.Unmarshal(data, &configurablesSlice) != nil {
		return nil, ErrNotJsonArray
	}
	return configurablesSlice, nil
}
