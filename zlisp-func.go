package elisp

import (
	elutils "github.com/rosbit/go-embedding-utils"
	"github.com/glycerine/zygomys/zygo"
	// "fmt"
	"reflect"
)

func (lpw *XZlisp) bindFunc(fn *zygo.SexpFunction, funcVarPtr interface{}) (err error) {
	helper, e := elutils.NewEmbeddingFuncHelper(funcVarPtr)
	if e != nil {
		err = e
		return
	}
	helper.BindEmbeddingFunc(lpw.wrapFunc(fn, helper))
	return
}

func (lpw *XZlisp) wrapFunc(fn *zygo.SexpFunction, helper *elutils.EmbeddingFuncHelper) elutils.FnGoFunc {
	return func(args []reflect.Value) (results []reflect.Value) {
		var lpArgs []zygo.Sexp

		// make lisp args
		itArgs := helper.MakeGoFuncArgs(args)
		for arg := range itArgs {
			lpArgs = append(lpArgs, toValue(arg, lpw.lisp))
		}

		// call lisp function
		res, err := lpw.lisp.Apply(fn, lpArgs)
		// fmt.Printf("res: %v, err: %v\n", res, err)

		// convert result to golang
		_, isResArray := res.(*zygo.SexpArray)
		results = helper.ToGolangResults(fromValue(res), isResArray, err)
		return
	}
}

func (lpw *XZlisp) callFunc(fn *zygo.SexpFunction, args ...interface{}) (res zygo.Sexp, err error) {
	lpArgs := make([]zygo.Sexp, len(args))
	for i, arg := range args {
		lpArgs[i] = toValue(arg, lpw.lisp)
	}

	return lpw.lisp.Apply(fn, lpArgs)
}
