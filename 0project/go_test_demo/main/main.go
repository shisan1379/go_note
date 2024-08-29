package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {

	go func() {
		defer func() {
			fmt.Print("123")
		}()
		time.Sleep(1 * time.Second)
		fmt.Print("1")
		panic("err")
	}()

	runtime.Goexit()

}
