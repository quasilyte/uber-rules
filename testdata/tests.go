package target

import (
	"fmt"
	"sync"
)

func sink(args ...interface{}) {}

func ifacePtr() {
	type structType struct {
		_ *fmt.Stringer // want `\Qdon't use pointers to an interface`
	}

	type ifacePtrAlias = *fmt.Stringer // want `\Qdon't use pointers to an interface`

	{
		var x *interface{} // want `\Qdon't use pointers to an interface`
		_ = x
		_ = *x
	}

	{
		var x **interface{} // want `\Qdon't use pointers to an interface`
		_ = x
		_ = *x
	}
}

func newMutex() {
	mu := new(sync.Mutex)
	_ = mu

	mu2 := new(sync.Mutex) // want `\Quse zero mutex value instead, 'var mu2 sync.Mutex'`
	mu2.Lock()
}

func channelSize() {
	_ = make(chan int, 1)    // OK: size of 1
	_ = make(chan string, 0) // OK: explicit size of 0
	_ = make(chan float32)   // OK: unbuffered, implicit size of 0

	size := 1
	_ = make(chan int, size) // OK: can't analyze

	_ = make(chan int, 2)     // want `\Qchannels should have a size of one or be unbuffered`
	_ = make(chan []int, 128) // want `\Qchannels should have a size of one or be unbuffered`
}

func uncheckedTypeAssert() {
	var v interface{}

	_ = v.(int) // want `\Qavoid unchecked type assertions as they can panic`
	{
		x := v.(int) // want `\Qavoid unchecked type assertions as they can panic`
		_ = x
	}

	sink(v.(int))          // want `\Qavoid unchecked type assertions as they can panic`
	sink(0, v.(int))       // want `\Qavoid unchecked type assertions as they can panic`
	sink(v.(int), 0)       // want `\Qavoid unchecked type assertions as they can panic`
	sink(1, 2, v.(int), 3) // want `\Qavoid unchecked type assertions as they can panic`

	{
		type structSink struct {
			f0 interface{}
			f1 interface{}
			f2 interface{}
		}
		_ = structSink{v.(int), 0, 0}    // want `\Qavoid unchecked type assertions as they can panic`
		_ = structSink{0, v.(string), 0} // want `\Qavoid unchecked type assertions as they can panic`
		_ = structSink{0, 0, v.([]int)}  // want `\Qavoid unchecked type assertions as they can panic`

		_ = structSink{f0: v.(int)}                  // want `\Qavoid unchecked type assertions as they can panic`
		_ = structSink{f0: 0, f1: v.(int)}           // want `\Qavoid unchecked type assertions as they can panic`
		_ = structSink{f0: 0, f1: 0, f2: v.(int)}    // want `\Qavoid unchecked type assertions as they can panic`
		_ = structSink{f0: v.(string), f1: 0, f2: 0} // want `\Qavoid unchecked type assertions as they can panic`
	}

	{
		_ = []interface{}{v.(int)}       // want `\Qavoid unchecked type assertions as they can panic`
		_ = []interface{}{0, v.(int)}    // want `\Qavoid unchecked type assertions as they can panic`
		_ = []interface{}{v.(int), 0}    // want `\Qavoid unchecked type assertions as they can panic`
		_ = []interface{}{0, v.(int), 0} // want `\Qavoid unchecked type assertions as they can panic`

		_ = [...]interface{}{10: v.(int)}               // want `\Qavoid unchecked type assertions as they can panic`
		_ = [...]interface{}{10: 0, 20: v.(int)}        // want `\Qavoid unchecked type assertions as they can panic`
		_ = [...]interface{}{10: v.(int), 20: 0}        // want `\Qavoid unchecked type assertions as they can panic`
		_ = [...]interface{}{10: 0, 20: v.(int), 30: 0} // want `\Qavoid unchecked type assertions as they can panic`
	}
}

func unnecessaryElse() {
	var cond bool

	{
		var x int
		if cond {
			x = 10
		} else {
			x = 5
		}
		_ = x
	}
}
