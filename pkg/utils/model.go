package utils

import (
	"reflect"
	"strings"
)

func GetPkgPath(model interface{}) string {
	sourceType := reflect.Indirect(reflect.ValueOf(model)).Type()
	pkg := sourceType.PkgPath()
	modelIndex := strings.Index(pkg, "model/")
	var name = SnackedName(sourceType.Name())
	if modelIndex > 0 {
		name = pkg[modelIndex+6:] + "/" + name
	}
	return name
}
