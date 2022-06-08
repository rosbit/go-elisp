package elisp

import (
	"github.com/glycerine/zygomys/zygo"
	"reflect"
	"time"
)

func toValue(v interface{}, env *zygo.Zlisp) zygo.Sexp {
	if v == nil {
		return zygo.SexpNull
	}

	switch vv := v.(type) {
	/*
	case rune:
		return &zygo.SexpChar{Val:vv}
	*/
	case int,int8,int16,int32,int64:
		return &zygo.SexpInt{Val:reflect.ValueOf(v).Int()}
	case uint,uint8,uint16,uint32,uint64:
		return &zygo.SexpUint64{Val:reflect.ValueOf(v).Uint()}
	case float32,float64:
		return &zygo.SexpFloat{Val:reflect.ValueOf(v).Float()}
	case string:
		return &zygo.SexpStr{S:vv}
	case []byte:
		return &zygo.SexpRaw{Val:vv}
	case bool:
		return &zygo.SexpBool{Val:vv}
	case time.Time:
		return &zygo.SexpTime{Tm:vv}
	case zygo.Sexp:
		return vv
	default:
		v2 := reflect.ValueOf(v)
		switch v2.Kind() {
		case reflect.Slice, reflect.Array:
			l := v2.Len()
			sexps := make([]zygo.Sexp, l)
			for i:=0; i<l; i++ {
				sexps[i] = toValue(v2.Index(i), env)
			}
			return zygo.MakeList(sexps)
		case reflect.Map:
			// h := &zygo.SexpHash{}
			sexps := make([]zygo.Sexp, v2.Len()*2)
			i := 0
			it := v2.MapRange()
			for it.Next() {
				// h.HashSet(toValue(it.Key().Interface(), env), toValue(it.Value().Interface(), env))
				sexps[i] = toValue(it.Key().Interface(), env)
				i += 1
				sexps[i] = toValue(it.Value().Interface(), env)
				i += 1
			}
			if res, err := zygo.MakeHash(sexps, "hash", env); err == nil {
				return res
			}
			return zygo.SexpNull
			// return h
		case reflect.Struct:
			sexps := bindGoStruct("", v2, env)
			if res, err := zygo.MakeHash(sexps, "struct", env); err == nil {
				return res
			}
			return zygo.SexpNull
		case reflect.Ptr:
			e := v2.Elem()
			if e.Kind() == reflect.Struct {
				sexps := bindGoStruct("", v2, env)
				if res, err := zygo.MakeHash(sexps, "struct", env); err == nil {
					return res
				}
				return zygo.SexpNull
			}
			return toValue(e.Interface(), env)
		case reflect.Func:
			if name, f, err := bindGoFunc("", v); err == nil {
				return zygo.MakeUserFunction(name, f)
			}
			return zygo.SexpNull
		case reflect.Interface:
			sexps := bindGoInterface(v2)
			if res, err := zygo.MakeHash(sexps, "struct", env); err == nil {
				return res
			}
			return zygo.SexpNull
		default:
			return zygo.SexpNull
		}
	}
}

func fromValue(v zygo.Sexp) (interface{}) {
	if v == nil || v == zygo.SexpNull {
		return nil
	}
	switch vv := v.(type) {
	case *zygo.SexpReflect:
		return vv.Val.Interface()
	case *zygo.SexpError:
		return vv
	case *zygo.SexpBool:
		return vv.Val
	case *zygo.SexpChar:
		return vv.Val
	case *zygo.SexpRaw:
		return vv.Val
	case *zygo.SexpInt:
		return vv.Val
	case *zygo.SexpUint64:
		return vv.Val
	case *zygo.SexpFloat:
		return vv.Val
	case *zygo.SexpStr:
		return vv.S
	case *zygo.SexpSymbol:
		return vv.Name()
	case *zygo.SexpArray:
		l := len(vv.Val)
		// fmt.Printf("l: %d\n", l)
		if l == 0 {
			return []interface{}{}
		}
		res := make([]interface{}, l)
		for i:=0; i<l; i++ {
			res[i] = fromValue(vv.Val[i])
		}
		// fmt.Printf("res: %v\n", res)
		return res
	case *zygo.SexpPair:
		if zygo.IsList(v) {
			l, err := zygo.ListToArray(v)
			if err != nil {
				res := make([]interface{}, len(l))
				for i, ll := range l {
					res[i] = fromValue(ll)
				}
				return res
			}
		}
		res := make(map[interface{}]interface{})
		key := fromValue(vv.Head)
		val := fromValue(vv.Tail)
		res[key] = val
		return res
	case *zygo.SexpHash:
		res := make(map[interface{}]interface{})
		lenKeyOrder := len(vv.KeyOrder)
		for i := 0; i < vv.NumKeys; i++ {
			var key, val zygo.Sexp
			var err error
			found := false
			for k:=i; i<lenKeyOrder; k++ {
				key = vv.KeyOrder[k]
				if val, err = vv.HashGet(nil, key); err == nil {
					found = true
					break
				}
			}
			if found {
				res[fromValue(key)] = fromValue(val)
			}
		}
		return res
	case *zygo.SexpFunction:
		return vv
	case *zygo.SexpTime:
		return time.Time(vv.Tm)
	default:
		return nil
	}
}

