package udwConsole

import (
	"reflect"
	"strings"
)

type CompletionBuilder struct {
	fnList []CompletionFn
}

type CompletionFn func(args []string) (accept bool, waitSelect []string)

func NewCompletionBuilder() *CompletionBuilder {
	return &CompletionBuilder{}
}

func (this *CompletionBuilder) Common(fn CompletionFn) *CompletionBuilder {
	this.fnList = append(this.fnList, fn)
	return this
}

func (this *CompletionBuilder) ReflectArg(fn interface{}) *CompletionBuilder {
	fv := reflect.ValueOf(fn)
	if fv.Kind() != reflect.Func || fv.Type().NumIn() != 1 {
		return this
	}
	av := fv.Type().In(0)
	if av.Kind() == reflect.Ptr {
		av = av.Elem()
	}
	if av.Kind() != reflect.Struct {
		return this
	}
	for idx := 0; idx < av.NumField(); idx++ {
		f := av.Field(idx)
		if f.Type.Kind() == reflect.Bool {
			this.AfterArgList(`-`+f.Name, []string{`True`, `False`})
		}
	}
	this.Common(func(args []string) (accept bool, waitSelect []string) {
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
	})
	return this
}

func (this *CompletionBuilder) AfterArgList(name string, list []string) *CompletionBuilder {
	return this.Common(func(args []string) (accept bool, waitSelect []string) {
		if len(args) < 2 {
			return false, nil
		}
		last := args[len(args)-2]
		if strings.ToLower(last) != strings.ToLower(name) {
			return false, nil
		}
		waitSelect = list
		return true, waitSelect
	})
}

func (this *CompletionBuilder) Finish() CompletionFn {
	return func(args []string) (accept bool, waitSelect []string) {
		for _, fn := range this.fnList {
			accept, waitSelect = fn(args)
			if accept {
				break
			}
		}
		return true, waitSelect
	}
}
