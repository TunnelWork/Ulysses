package server

import (
	"encoding/json"
)

// JSON utils

type Configurables map[string]string

// JsonToConfigurables returns map[string]string representing a Json Object ({, , ,})
func JsonToConfigurables(data []byte) (map[string]string, error) {
	var configurables map[string]string
	if json.Unmarshal(data, &configurables) != nil {
		return nil, ErrBadJsonObject
	}
	return configurables, nil
}

// JsonArrToConfigurablesSlice returns a slice of map[string]string representing a Json Array ([{}, {}, {}])
func JsonArrToConfigurablesSlice(data []byte) ([]map[string]string, error) {
	var configurablesSlice []map[string]string
	if json.Unmarshal(data, &configurablesSlice) != nil {
		return nil, ErrBadJsonArray
	}
	return configurablesSlice, nil
}
