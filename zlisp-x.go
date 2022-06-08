package elisp

import (
	"github.com/glycerine/zygomys/zygo"
	"fmt"
	"os"
	"reflect"
)

func New() *XZlisp {
	return &XZlisp {
		lisp: zygo.NewZlisp(),
	}
}

func (lpw *XZlisp) LoadFile(path string, vars map[string]interface{}) (err error) {
	fp, e := os.Open(path)
	if e != nil {
		err = e
		return
	}
	defer fp.Close()

	if err = lpw.lisp.LoadFile(fp); err != nil {
		return
	}

	err = lpw.setEnv(vars)
	return
}

func (lpw *XZlisp) LoadScript(script string, vars map[string]interface{}) (err error) {
	if err = lpw.lisp.LoadString(script); err != nil {
		return
	}
	err = lpw.setEnv(vars)
	return
}

func (lpw *XZlisp) Dump() {
	lpw.lisp.DumpEnvironment()
}

func (lpw *XZlisp) GetGlobal(name string) (res interface{}, err error) {
	r, e := lpw.getVar(name)
	if e != nil {
		err = e
		return
	}
	res = fromValue(r)
	return
}

func (lpw *XZlisp) EvalFile(path string, env map[string]interface{}) (res interface{}, err error) {
	if err = lpw.LoadFile(path, env); err != nil {
		return
	}
	v, e := lpw.lisp.Run()
	if e != nil {
		err = e
		return
	}
	res = fromValue(v)
	return
}

func (lpw *XZlisp) Eval(script string, env map[string]interface{}) (res interface{}, err error) {
	if err = lpw.LoadScript(script, env); err != nil {
		return
	}
	v, e := lpw.lisp.Run()
	if e != nil {
		err = e
		return
	}
	res = fromValue(v)
	return
}

func (lpw *XZlisp) CallFunc(funcName string, args ...interface{}) (res interface{}, err error) {
	v, e := lpw.getVar(funcName)
	if e != nil {
		err = e
		return
	}
	fn, ok := v.(*zygo.SexpFunction)
	if !ok {
		err = fmt.Errorf("var %s is not with type function", funcName)
		return
	}

	r, e := lpw.callFunc(fn, args...)
	if e != nil {
		err = e
		return
	}
	res = fromValue(r)
	return
}

// bind a var of golang func with a ZLisp function name, so calling ZLisp function
// is just calling the related golang func.
// @param funcVarPtr  in format `var funcVar func(....) ...; funcVarPtr = &funcVar`
func (lpw *XZlisp) BindFunc(funcName string, funcVarPtr interface{}) (err error) {
	if funcVarPtr == nil {
		err = fmt.Errorf("funcVarPtr must be a non-nil poiter of func")
		return
	}
	t := reflect.TypeOf(funcVarPtr)
	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Func {
		err = fmt.Errorf("funcVarPtr expected to be a pointer of func")
		return
	}

	v, e := lpw.getVar(funcName)
	if e != nil {
		err = e
		return
	}
	fn, ok := v.(*zygo.SexpFunction)
	if !ok {
		err = fmt.Errorf("var %s is not with type function", funcName)
		return
	}
	lpw.bindFunc(fn, funcVarPtr)
	return
}

func (lpw *XZlisp) BindFuncs(funcName2FuncVarPtr map[string]interface{}) (err error) {
	for funcName, funcVarPtr := range funcName2FuncVarPtr {
		if err = lpw.BindFunc(funcName, funcVarPtr); err != nil {
			return
		}
	}
	return
}

// make a golang func as a built-in ZLisp function, so the function can be called in ZLisp script.
func (lpw *XZlisp) MakeUserFunc(funcName string, funcVar interface{}) (err error) {
	_, goFunc, e := bindGoFunc(funcName, funcVar)
	if e != nil {
		err = e
		return
	}
	lpw.lisp.AddFunction(funcName, goFunc)
	return
}

func (lpw *XZlisp) setEnv(vars map[string]interface{}) (err error) {
	if len(vars) == 0 {
		return nil
	}
	for k, v := range vars {
		lispV := toValue(v, lpw.lisp)
		lpw.lisp.AddGlobal(k, lispV)
	}
	return nil
}

func (lpw *XZlisp) getVar(name string) (v zygo.Sexp, err error) {
	var ok bool
	v, ok = lpw.lisp.FindObject(name)
	if !ok {
		err = fmt.Errorf("no var named %s found", name)
		return
	}
	return
}
