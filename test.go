package main

import (
     "time"
     "fmt"
)

func main() {
    timer := time.NewTimer(1 * time.Second)
    timer2 := time.NewTimer(500*  time.Millisecond)
//    for i := 0; i < 2; i++ {
//            timer.Reset(1 * time.Second)
    L:
    for {
        select {
        case <- timer.C:
            fmt.Printf("Timer 1 done!\n")
            break L
        case <-timer2.C:
            fmt.Printf("Resetting timers!\n")
            timer.Reset(5* time.Second)
        }
    }
    
    fmt.Printf("Done!\n")
//    timer.Reset(3*time.Second)
//    <-timer.C
//    }
}
