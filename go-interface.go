package elisp

import (
	"github.com/glycerine/zygomys/zygo"
	"reflect"
)

func bindGoInterface(intVar reflect.Value) (goStruct []zygo.Sexp) {
	count := intVar.NumMethod()
	goStruct = make([]zygo.Sexp, count*2)
	t := intVar.Type()
	i := 0
	for j := 0; j < count; j++ {
		name := lowerFirst(t.Method(j).Name)
		goStruct[i] = toValue(name, nil); i+=1
		_, f, _ := bindGoFunc(name, intVar.Method(j))
		goStruct[i] = zygo.MakeUserFunction(name, f); i += 1
	}
	return
}

