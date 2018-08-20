package main

import (
	"fmt"
	"time"
)

// func synchronizer(l chan int, s chan bool){
// 	for <- request {

// 	}
// }

func fill(messages chan int, s chan int, val int) {
	fmt.Printf("Created wait routine: %v\n", val)
	p:= <- s
	messages <- val
	fmt.Printf("Put value:%v, taken from queue: %v\n", val, p)
}

func main() {
	messages := make(chan int, 2)
	s := make(chan int, 5)
	
	// go syncroniser(l, s)

	for i := 5;i <= 6;i++ {
		s <- i
	}

	for i:=0;i<=4;i++{
		go fill(messages, s, i)
	}

	

	

	for p:=0;p<=4;p++{
		time.Sleep(time.Second * 1)
		s <- p
		val:= <- messages
		fmt.Printf("Got from goroutine: %v\n", val)
		// switch {
		// case :
			
		// }
	}


	fmt.Printf("DONE! \n")
}
