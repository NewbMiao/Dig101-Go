package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestModifyPointerArg(m *testing.T) {
	arg := ArgType{"arg0", 0}
	fmt.Println("init arg:", arg)

	modifyPointerArg1(&arg)
	fmt.Println("after call modifyPointerArg1", arg)
	modifyPointerArg2(&arg)
	fmt.Println("after call modifyPointerArg2", arg)
	modifyPointerArg3(&arg)
	fmt.Println("after call modifyPointerArg3", arg)
	modifyPointerArg4(&arg)
	fmt.Println("after call modifyPointerArg4", arg)

	arg10 := &ArgType{"arg10", 10}
	fmt.Println("\ninit arg10:", arg10)

	arg10.modifyPointerReceiver1()
	fmt.Println("after call modifyPointerReceiver1", arg10)
	arg10.modifyPointerReceiver2()
	fmt.Println("after call modifyPointerReceiver2", arg10)
	arg10.modifyPointerReceiver3()
	fmt.Println("after call modifyPointerReceiver3", arg10)
	arg10.modifyPointerReceiver4()
	fmt.Println("after call modifyPointerReceiver4", arg10)
}

type ArgType struct {
	A string
	b int
}

func modifyPointerArg1(arg *ArgType) {
	arg = &ArgType{"arg1", 1}
	fmt.Println("inside modifyPointerArg1:", arg)
}

func modifyPointerArg2(arg *ArgType) {
	*arg = ArgType{"arg2", 2}
	fmt.Println("inside modifyPointerArg2:", arg)
}

func modifyPointerArg3(arg *ArgType) {
	val := reflect.ValueOf(arg)
	val.Elem().FieldByName("A").SetString("arg3")
	fmt.Println("inside modifyPointerArg3:", arg)
	// val.Elem().FieldByName("b").SetInt(3)
	// panic: reflect: reflect.flag.mustBeAssignable using value obtained using unexported field
}

func modifyPointerArg4(arg *ArgType) {
	jsonStr := `{"A":"arg4","b":4}`
	json.Unmarshal([]byte(jsonStr), arg)
}

func (arg *ArgType) modifyPointerReceiver1() {
	arg = &ArgType{"arg1", 1}
	fmt.Println("inside modifyPointerArg1:", arg)
}

func (arg *ArgType) modifyPointerReceiver2() {
	*arg = ArgType{"arg2", 2}
	fmt.Println("inside modifyPointerArg2:", arg)
}

func (arg *ArgType) modifyPointerReceiver3() {
	val := reflect.ValueOf(arg)
	val.Elem().FieldByName("A").SetString("arg3")
	fmt.Println("inside modifyPointerArg3:", arg)
	// val.Elem().FieldByName("b").SetInt(3)
	// panic: reflect: reflect.flag.mustBeAssignable using value obtained using unexported field
}

func (arg *ArgType) modifyPointerReceiver4() {
	jsonStr := `{"A":"arg4","b":4}`
	json.Unmarshal([]byte(jsonStr), arg)
}
