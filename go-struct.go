package elisp

import (
	"github.com/glycerine/zygomys/zygo"
	"reflect"
	"strings"
	// "fmt"
)

func bindGoStruct(name string, structVar reflect.Value, env *zygo.Zlisp) (goStruct []zygo.Sexp) {
	var structE reflect.Value
	if structVar.Kind() == reflect.Ptr {
		structE = structVar.Elem()
	} else {
		structE = structVar
	}
	structT := structE.Type()

	if structE == structVar {
		// struct is unaddressable, so make a copy of struct to an Elem of struct-pointer.
		// NOTE: changes of the copied struct cannot effect the original one. it is recommended to use the pointer of struct.
		structVar = reflect.New(structT) // make a struct pointer
		structVar.Elem().Set(structE)    // copy the old struct
		structE = structVar.Elem()       // structE is the copied struct
	}

	if len(name) == 0 {
		n := structT.Name()
		if len(n) > 0 {
			if pos := strings.LastIndex(n, "."); pos >= 0 {
				name = n[pos+1:]
			} else {
				name = n
			}
		}
		if len(name) == 0 {
			name = "noname"
		}
	}

	goStruct = getAttrs(structVar, structE, structT, env)
	return
}

func lowerFirst(name string) string {
	return strings.ToLower(name[:1]) + name[1:]
}
func upperFirst(name string) string {
	return strings.ToUpper(name[:1]) + name[1:]
}

func getAttrs(structVar, structE reflect.Value, structT reflect.Type, env *zygo.Zlisp) []zygo.Sexp {
	count := structT.NumField() + structVar.NumMethod() + structE.NumMethod()
	sexps := make([]zygo.Sexp, count*2)
	i := 0
	// fmt.Printf("structT.NumField(): %d\n", structT.NumField())
	for j:=0; j<structT.NumField(); j++ {
		name := lowerFirst(structT.Field(j).Name)
		v := structE.Field(j)
		if !v.CanInterface() {
			continue
		}
		sexps[i] = toValue(name, env); i += 1
		// fmt.Printf("j: %d, v: %v\n", j, v)
		sexps[i] = toValue(v.Interface(), env); i += 1
	}
	for j:=0; j<structE.NumMethod(); j++ {
		name := lowerFirst(structT.Method(j).Name)
		sexps[i] = toValue(name, env); i += 1
		_, f, _ := bindGoFunc(name, structE.Method(j))
		sexps[i] = zygo.MakeUserFunction(name, f); i += 1
	}
	t := structVar.Type()
	for j:=0; j<structVar.NumMethod(); j++ {
		name := lowerFirst(t.Method(j).Name)
		sexps[i] = toValue(name, env); i += 1
		_, f, _ :=  bindGoFunc(name, structVar.Method(j))
		sexps[i] = zygo.MakeUserFunction(name, f); i += 1
	}
	return sexps[:i]
}

