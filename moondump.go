package moondump

import (
	"fmt"
	"reflect"
)

type Some struct {
	value string
	id    int
}

func main() {
	instance := Some{
		value: "Hey there",
		id:    1,
	}

	typeObj := reflect.TypeOf(instance)
	name := typeObj.Name()
	pkg := typeObj.PkgPath()

	fmt.Println(name, pkg)

}
