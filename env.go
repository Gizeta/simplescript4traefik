package simplescript4traefik

import (
	"fmt"
	"net/http"
	"strings"
)

type Env map[string]Val

var nullVal = Val{
	tag: tNull,
}

var trueVal = Val{
	tag: tBool,
	str: "true",
}

var falseVal = Val{
	tag: tBool,
	str: "false",
}

func Builtin_Get(env *Env, args []Val) Val {
	key := args[0].AsSymbol()
	if val, ok := (*env)[key]; ok {
		return val
	}
	return nullVal
}

func Builtin_Set(env *Env, args []Val) Val {
	key := args[0].AsSymbol()
	val := Eval(args[1], env)
	(*env)[key] = val
	return nullVal
}

func Builtin_Not(env *Env, args []Val) Val {
	cond := Eval(args[0], env).AsBool()
	if cond {
		return falseVal
	}
	return trueVal
}

func Builtin_Equal(env *Env, args []Val) Val {
	val1 := Eval(args[0], env)
	val2 := Eval(args[1], env)
	if val1.tag != val2.tag {
		return falseVal
	}
	if val1.tag == tBool {
		if val1.str == "true" {
			if val2.str == "true" {
				return trueVal
			}
			return falseVal
		} else {
			if val2.str == "true" {
				return falseVal
			}
			return trueVal
		}
	} else if val1.tag == tNumber {
		if val1.num == val2.num {
			return trueVal
		}
		return falseVal
	} else if val1.tag == tNull {
		return trueVal
	} else if val1.tag == tString {
		if val1.str == val2.str {
			return trueVal
		}
		return falseVal
	}
	return falseVal
}

func Builtin_If(env *Env, args []Val) Val {
	condition := Eval(args[0], env).AsBool()
	if condition {
		return Eval(args[1], env)
	} else {
		return Eval(args[2], env)
	}
}

func Builtin_StrContains(env *Env, args []Val) Val {
	source := Eval(args[0], env).AsString()
	substr := Eval(args[1], env).AsString()
	if strings.Contains(source, substr) {
		return trueVal
	}
	return falseVal
}

func MakeBuiltinFunc(fn func(*Env, []Val) Val) Val {
	return Val{
		tag: tFunc,
		fn: fn,
	}
}

func Builtin_GetReqPath(rw *http.ResponseWriter, req *http.Request, next *http.Handler, env *Env, args []Val) Val {
	return Val{
		tag: tString,
		str: req.URL.Path,
	}
}

func Builtin_GetReqHeader(rw *http.ResponseWriter, req *http.Request, next *http.Handler, env *Env, args []Val) Val {
	key := Eval(args[0], env).AsString()
	return Val{
		tag: tString,
		str: req.Header.Get(key),
	}
}

func Builtin_SetReqHeader(rw *http.ResponseWriter, req *http.Request, next *http.Handler, env *Env, args []Val) Val {
	key := Eval(args[0], env).AsString()
	val := Eval(args[1], env).AsString()
	req.Header.Set(key, val)
	return nullVal
}

func Builtin_Halt(rw *http.ResponseWriter, req *http.Request, next *http.Handler, env *Env, args []Val) Val {
	code := Eval(args[0], env).AsNumber()
	msg := Eval(args[1], env).AsString()
	http.Error(*rw, msg, int(code))
	return nullVal
}

func Builtin_RunNext(rw *http.ResponseWriter, req *http.Request, next *http.Handler, env *Env, args []Val) Val {
	(*next).ServeHTTP(*rw, req)
	return nullVal
}

func RegisterTraefikBuiltin(env *Env, rw *http.ResponseWriter, req *http.Request, next *http.Handler) {
	MakeTraefikBuiltFunc := func(fn func(*http.ResponseWriter, *http.Request, *http.Handler, *Env, []Val) Val) Val {
		return Val{
			tag: tFunc,
			fn: func(env *Env, args []Val) Val {
				return fn(rw, req, next, env, args)
			},
		}
	}

	(*env)["get_req_path"] = MakeTraefikBuiltFunc(Builtin_GetReqPath)
	(*env)["get_req_header"] = MakeTraefikBuiltFunc(Builtin_GetReqHeader)
	(*env)["set_req_header"] = MakeTraefikBuiltFunc(Builtin_SetReqHeader)
	(*env)["halt"] = MakeTraefikBuiltFunc(Builtin_Halt)
	(*env)["run_next"] = MakeTraefikBuiltFunc(Builtin_RunNext)
}

func CreateEnv() Env {
	env := make(Env)
	env["null"] = nullVal
	env["true"] = trueVal
	env["false"] = falseVal
	env["get"] = MakeBuiltinFunc(Builtin_Get)
	env["set"] = MakeBuiltinFunc(Builtin_Set)
	env["!"] = MakeBuiltinFunc(Builtin_Not)
	env["="] = MakeBuiltinFunc(Builtin_Equal)
	env["if"] = MakeBuiltinFunc(Builtin_If)
	env["str_contains"] = MakeBuiltinFunc(Builtin_StrContains)
	return env
}
