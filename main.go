package main

//Package context defines the Context type, which carries deadlines, cancellation signals, and other request-scoped values across API boundaries and between processes.

import (
	"context"
	"fmt"
	"time"
)

type User struct {
	name string
	age  uint8
}

type contextKey string

//channel of type send
func gen_numbers(ctx context.Context, out chan<- int) {
	n := 1
	for {
		select {
		//returns a channel that's closed when work done on behalf of this
		// context should be canceled

		// WithCancel arranges for Done to be closed when cancel is called;
		// WithDeadline arranges for Done to be closed when the deadline
		// expires; WithTimeout arranges for Done to be closed when the timeout
		// elapses.
		case <-ctx.Done():
			fmt.Println("gen_numbers closed ", ctx.Err())
			return
		case out <- n:
			n++
			time.Sleep(300 * time.Millisecond)
		}
	}
}

func gen_numbers_ret_channel(ctx context.Context) chan int {
	n := 1
	dst := make(chan int)
	go func() {
		for {
			select {
			case val := <-ctx.Done():
				fmt.Println("gen_numbers closed ", ctx.Err(), " ", val)
				close(dst)
				return
			case dst <- n:
				n++
				time.Sleep(200 * time.Millisecond)
			}
		}
	}()
	return dst
}

func gen_use_ctx_value(ctx context.Context) chan User {
	ch := make(chan User)
	u := User{
		age:  3,
		name: "Gigel",
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("The context is done ", ctx.Err())
				close(ch)
				return
			case ch <- u:
				time.Sleep(200 * time.Millisecond)
				fmt.Printf("Send user %v\n", u)
			}
		}
	}()
	return ch
}

func main_2_use_cancel() {
	//ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*1750)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for n := range gen_numbers_ret_channel(ctx) {
		if n == 5 {
			break
		}
		fmt.Println("Got ", n)
	}
}

func main_3_use_timeout() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*1750)
	defer cancel()
	for n := range gen_numbers_ret_channel(ctx) {
		fmt.Println("Got ", n)
	}
}

func main_4_use_values() {
	k := contextKey("Helo")
	u := &User{
		age:  30,
		name: "Mark",
	}

	ctx := context.WithValue(context.Background(), k, u)
	d := time.Now().Add(2800 * time.Millisecond)
	ctx2, cancel := context.WithDeadline(ctx, d)
	defer cancel()
	for u := range gen_use_ctx_value(ctx2) {
		fmt.Printf("User %v\n", u)
	}
}

func main_1() {
	//out := make(chan<- int)
	out := make(chan int)
	should_break := false
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() //should not be necessary
	go gen_numbers(ctx, out)
	for !should_break {
		select {
		case <-ctx.Done():
			should_break = true
			cancel()
		case n := <-out:
			fmt.Println("Got ", n)
			if n > 10 {
				fmt.Println("Cancelling")
				cancel()
				//should_break = true
				break
			}
		}
	}

}

func main() {
	main_4_use_values()
}
