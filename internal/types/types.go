package types

import (
	"fmt"
	"reflect"
)

func GetType(typeName string, pkgName string) (reflect.Type, error) {
	typeObj, found := reflect.TypeOf((*interface{})(nil)).Elem().FieldByName(typeName)

	if !found {
		return nil, fmt.Errorf("type not found: %s", typeName)
	}

	return typeObj.Type, nil
}
