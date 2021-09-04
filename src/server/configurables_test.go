package server

import "testing"

var (
	ExampleJson = []byte(`{
		"First Name": "Gaukas",
		"Last Name": "Wang",
		"GitHub": "@Gaukas",
		"Hobby": "Cooking"
	}`)
	ExampleJsonArr = []byte(`[
		{
			"FirstName": "Gaukas",
			"LastName": "Wang",
			"GitHub": "@Gaukas",
			"Hobby": "Cooking"
		},
		{
			"FirstName": "Aaron",
			"LastName": "Li",
			"GitHub": "@yl4579",
			"Hobby": "NotCooking"
		},
		{
			"FirstName": "Milkey",
			"LastName": "Tan",
			"GitHub": "@mili-tan",
			"Hobby": "Building Wheels"
		}
	]`)
)

func TestJsonToConfigurables(t *testing.T) {
	Configurables, err := JsonToConfigurables(ExampleJson)
	if err != nil {
		t.Errorf("JsonToConfigurables() returns error:%s\n", err)
	}
	ServerConfigurables := ServerConfigurables(Configurables)
	AccountConfigurables := AccountConfigurables(Configurables)

	for key, value := range ServerConfigurables {
		if value != AccountConfigurables[key] {
			t.Errorf("Inconsistent results after casting type.\n")
		} else {
			t.Logf(`ServerConf['%s'] == AccountConf['%s'] = '%s'`, key, key, value)
		}
	}
}

func TestJsonArrToConfigurablesSlice(t *testing.T) {
	ConfigurablesSlice, err := JsonArrToConfigurablesSlice(ExampleJsonArr)
	if err != nil {
		t.Errorf("JsonArrToConfigurablesSlice() returns error:%s\n", err)
	}

	for idx, conf := range ConfigurablesSlice {
		for key, value := range conf {
			t.Logf(`ConfigurablesSlice[%d]['%s'] == '%s'`, idx, key, value)
		}
	}
}
