package main

import "fmt"

var (
	thirdPartyInitList map[string]*func()
)

func thirdPartyInit() {
	for fname, f := range thirdPartyInitList {
		Info("thirdPartyInit(): ", fname)
		(*f)()
	}
}

func RegisterThirdPartyInitFunc(uniqueName string, initFunc *func()) bool {
	if _, exist := thirdPartyInitList[uniqueName]; exist {
		return false
	}
	Debug("RegisterThirdPartyInitFunc(): ", uniqueName, "@", fmt.Sprintf("%x", initFunc))
	thirdPartyInitList[uniqueName] = initFunc
	return true
}
