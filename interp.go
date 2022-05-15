package simplescript4traefik

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	tNull = iota
	tBool
	tNumber
	tString
	tSymbol
	tFunc
	tList
	tExpr
)

type Val struct {
	tag byte
	num float64
	str string
	fn func(*Env, []Val) Val
	args *[]Val
}

func (v Val) String() string {
	if v.tag == tNull {
		return "null"
	} else if v.tag == tNumber {
		return fmt.Sprintf("%f", v.num)
	} else if v.tag == tBool {
		if v.str == "true" {
			return "true"
		}
		return "false"
	} else if v.tag == tString {
		return "\"" + v.str + "\""
	} else if v.tag == tSymbol {
		return v.str
	} else if v.tag == tFunc {
		return "(func)"
	} else if v.tag == tList {
		var l []string
		for _, i := range *v.args {
			l = append(l, i.String())
		}
		return "(list:" + strings.Join(l, ",") + ")"
	} else if v.tag == tExpr {
		var l []string
		for _, i := range (*v.args)[1:] {
			l = append(l, i.String())
		}
		return "(expr " + (*v.args)[0].String() + ":" + strings.Join(l, ",") + ")"
	}
	return ""
}

func (v Val) AsString() string {
	if v.tag == tString {
		return v.str
	} else if v.tag == tNumber {
		return fmt.Sprintf("%f", v.num)
	} else if v.tag == tBool {
		if v.str == "true" {
			return "true"
		}
		return "false"
	} else if v.tag == tSymbol {
		return v.str
	}
	panic("string type mismatched")
}

func (v Val) AsBool() bool {
	if v.tag == tBool {
		return v.str == "true"
	} else if v.tag == tString {
		return v.str != ""
	} else if v.tag == tNumber {
		return v.num != 0
	} else if v.tag == tList {
		return len(*(v.args)) > 0
	} else if v.tag == tNull {
		return false
	}
	panic("bool type mismatched")
}

func (v Val) AsSymbol() string {
	if v.tag == tSymbol {
		return v.str
	} else if v.tag == tString {
		return v.str
	}
	panic("symbol type mismatched")
}

func (v Val) AsNumber() float64 {
	if v.tag == tNumber {
		return v.num
	} else if v.tag == tBool {
		if v.str == "true" {
			return 1
		}
		return 0
	}
	panic("number type mismatched")
}

func Tokenize(chars string) []string {
	chars = strings.ReplaceAll(chars, "(", " ( ")
	chars = strings.ReplaceAll(chars,  ")", " ) ")
	chars = strings.ReplaceAll(chars,  "\n", " ")
	var l []string
	for _, str := range strings.Split(chars, " ") {
		if str != "" {
			l = append(l, str)
		}
	}
	return l
}

func ReadFromTokens(tokens []string, args ...int) ([]Val, int) {
	var l []Val
	var pos int
	if len(args) > 0 {
		pos = args[0]
	} else {
		pos = 0
	}
	for pos < len(tokens) {
		token := tokens[pos]
		pos++
		if token == "(" {
			var a []Val
			a, pos = ReadFromTokens(tokens, pos)
			l = append(l, Val{
				tag: tExpr,
				args: &a,
			})
		} else if token == ")" {
			return l, pos
		} else {
			if strings.HasPrefix(token, "\"") {
				next := token
				for pos < len(tokens) && !strings.HasSuffix(next, "\"") {
					next = tokens[pos]
					pos++
					token += " " + next
				}
				token = strings.TrimPrefix(token, "\"")
				token = strings.TrimSuffix(token, "\"")
				l = append(l, Val{
					tag: tString,
					str: token,
				})
			} else {
				num, err := strconv.ParseFloat(token, 64)
				if err != nil {
					l = append(l, Val{
						tag: tSymbol,
						str: token,
					})
				} else {
					l = append(l, Val{
						tag: tNumber,
						num: num,
					})
				}
			}
		}
	}
	return l, pos
}

func Eval(v Val, env *Env) Val {
	if v.tag == tExpr {
		if len(*(v.args)) == 0 {
			return nullVal
		}
		first := Eval((*v.args)[0], env)
		if first.tag == tFunc {
			var args []Val
			for _, i := range (*v.args)[1:] {
				args = append(args, i)
			}
			return first.fn(env, args)
		} else {
			var l []Val
			l = append(l, first)
			for _, i := range (*v.args)[1:] {
				l = append(l, Eval(i, env))
			}
			return Val{
				tag: tList,
				args: &l,
			}
		}
	} else if v.tag == tSymbol {
		return (*env)[v.str]
	} else if v.tag == tNull {
		return nullVal
	}
	return v
}

func RunScript(code string, env *Env) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("panic occurred:", err)
		}
	}()

	tokens := Tokenize(code)
	expr, _ := ReadFromTokens(tokens)
	program := Val{
		tag: tExpr,
		args: &expr,
	}
	Eval(program, env)
}
