package udwConsole

import (
	"reflect"
	"strings"
)

type CreateCompletionReq struct {
	AfterArgListMap  map[string][]string
	ReflectArgOfFunc interface{}
	SelfDefineFnList []CompletionFn
}

func CreateCompletion(req CreateCompletionReq) CompletionFn {
	var fnList []CompletionFn
	for pre, after := range req.AfterArgListMap {
		fnList = append(fnList, afterArgList(pre, after))
	}
	fn := reflectArg(req.ReflectArgOfFunc)
	if fn != nil {
		fnList = append(fnList, fn)
	}
	fnList = append(fnList, req.SelfDefineFnList...)
	return func(args []string) (accept bool, waitSelect []string) {
		for _, fn := range fnList {
			accept, waitSelect = fn(args)
			if accept {
				break
			}
		}
		return true, waitSelect
	}
}

type CompletionFn func(args []string) (accept bool, waitSelect []string)

func reflectArg(fn interface{}) CompletionFn {
	fv := reflect.ValueOf(fn)
	if fv.Kind() != reflect.Func || fv.Type().NumIn() != 1 {
		return nil
	}
	av := fv.Type().In(0)
	if av.Kind() == reflect.Ptr {
		av = av.Elem()
	}
	if av.Kind() != reflect.Struct {
		return nil
	}
	return func(args []string) (accept bool, waitSelect []string) {
		if len(args) < 1 {
			return false, nil
		}
		last := args[len(args)-1]
		if !strings.HasPrefix(last, `-`) {
			return false, nil
		}
		last = strings.TrimPrefix(strings.ToLower(last), `-`)
		for idx := 0; idx < av.NumField(); idx++ {
			name := av.Field(idx).Name
			tmp := strings.Split(av.Field(idx).Tag.Get(`CmdFlag`), `,`)
			if len(tmp) > 0 && tmp[0] != `` {
				name = tmp[0]
			}
			if strings.HasPrefix(strings.ToLower(name), last) {
				waitSelect = append(waitSelect, `-`+name)
			}
		}
		return true, waitSelect
	}
}

func afterArgList(name string, list []string) CompletionFn {
	return func(args []string) (accept bool, waitSelect []string) {
		if len(args) < 2 {
			return false, nil
		}
		last := args[len(args)-2]
		if strings.ToLower(last) != strings.ToLower(name) {
			return false, nil
		}
		waitSelect = list
		return true, waitSelect
	}
}
