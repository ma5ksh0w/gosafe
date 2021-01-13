package gosafe_test

import (
	"fmt"
	"testing"

	"github.com/ma5ksh0w/gosafe"
)

func TestFailed(t *testing.T) {
	if gosafe.Failed(func() {
		fmt.Println("no panic, must be false")
	}) {
		t.Fatal("Failed func with no panic")
	}

	if !gosafe.Failed(func() {
		panic("this is a panic")
	}) {
		t.Fatal("Failed func returns true")
	}
}

func TestAsError(t *testing.T) {
	if err := gosafe.AsError(func() {
		fmt.Println("no panic, no error")
	}); err != nil {
		t.Fatal("No panic func returns error:", err)
	}

	err := gosafe.AsError(func() {
		panic("this is a panic")
	})

	if err == nil {
		t.Fatal("Panic cannot returns error")
	} else {
		fmt.Println("Panic returns error:", err)
	}
}

func TestCatch(t *testing.T) {
	gosafe.Catch(func() {
		fmt.Println("no panic, no error")
	}, func(e interface{}) {
		t.Fatal("Catch called callback with no panic, e:", e)
	})

	called := false
	gosafe.Catch(func() {
		panic("must be catched")
	}, func(e interface{}) {
		called = true
		fmt.Println("Catch:", e)
	})

	if !called {
		t.Fatal("Callback is not called on panic")
	}
}

func TestCatchCh(t *testing.T) {
	pCh := make(chan interface{})
	done := make(chan chan int)

	// example panic handler
	go func() {
		ctr := 0
		for {
			select {
			case e := <-pCh:
				fmt.Println("Panic handler: catched panic:", e)
				ctr++

			case ch := <-done:
				select {
				case <-ch:
				case ch <- ctr:
				}

				return
			}
		}
	}()

	gosafe.CatchCh(func() {
		fmt.Println("no panic, no error")
	}, pCh)

	gosafe.CatchCh(func() {
		panic("panic #1")
	}, pCh)

	gosafe.CatchCh(func() {
		panic("panic #2")
	}, pCh)

	ret := make(chan int)
	select {
	case <-done:
	case done <- ret:
	}

	cnt := <-ret
	if cnt != 2 {
		t.Fatal("invalid panics count, want 2, got", cnt)
	}
}

func TestIgnore(t *testing.T) {
	gosafe.Ignore(func() {
		panic("must be ignored")
	})
}
