package elisp

import (
	elutils "github.com/rosbit/go-embedding-utils"
	"github.com/glycerine/zygomys/zygo"
)

func bindGoFunc(name string, funcVar interface{}) (realName string, goFunc zygo.ZlispUserFunction, err error) {
	helper, e := elutils.NewGolangFuncHelper(name, funcVar)
	if e != nil {
		err = e
		return
	}

	realName = helper.GetRealName()
	goFunc = wrapGoFunc(helper)
	return
}

func wrapGoFunc(helper *elutils.GolangFuncHelper) zygo.ZlispUserFunction {
	return func(env *zygo.Zlisp, name string, args []zygo.Sexp) (val zygo.Sexp, err error) {
		getArgs := func(i int) interface{} {
			return fromValue(args[i])
		}

		v, e := helper.CallGolangFunc(len(args), name, getArgs)
		if e != nil {
			err = e
			return
		}
		if v == nil {
			val = zygo.SexpNull
			return
		}

		if vv, ok := v.([]interface{}); ok {
			retV := make([]zygo.Sexp, len(vv))
			for i, rv := range vv {
				retV[i] = toValue(rv, env)
			}
			val = zygo.MakeList(retV)
		} else {
			val = toValue(v, env)
		}
		return
	}
}
