package utils

import (
	"github.com/tinwoan-go/basic-api/logger"
	"reflect"
)

func MapDestructor(target interface{}) {
	if m, ok := target.(map[string]interface{}); ok {
		singleMap(m)
		return
	}
	if lm, ok := target.([]map[string]interface{}); ok {
		listMap(lm)
		return
	}
	logger.Warn("%v is not a map or a list of map", reflect.TypeOf(target))
}

func singleMap(m map[string]interface{}) {
	for key := range m {
		delete(m, key)
	}
	m = nil
}

func listMap(lm []map[string]interface{}) {
	for _, m := range lm {
		singleMap(m)
	}
}