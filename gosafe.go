package gosafe

import "fmt"

type ErrPanic struct {
	e interface{}
}

func (err *ErrPanic) Error() string {
	return fmt.Sprint("Panic:", err.e)
}

// Failed returns true, if given func be crashed with panic
func Failed(fn func()) (ok bool) {
	defer func() {
		if e := recover(); e != nil {
			ok = true
		}
	}()

	fn()
	return
}

// AsError returns panic as error
func AsError(fn func()) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = &ErrPanic{e}
		}
	}()

	fn()
	return
}

// Catch execute given func as callback
func Catch(fn func(), callback func(interface{})) {
	defer func() {
		if e := recover(); e != nil {
			callback(e)
		}
	}()

	fn()
}

// CatchCh sends panic to channel
func CatchCh(fn func(), ch chan interface{}) {
	defer func() {
		if e := recover(); e != nil {
			select {
			case <-ch:
			case ch <- e:
			}
		}
	}()

	fn()
}

// Ignore do nothing on panic
func Ignore(fn func()) {
	defer func() {
		_ = recover()
	}()

	fn()
}
