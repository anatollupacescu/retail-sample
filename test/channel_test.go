package test

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func _(_ *testing.T) {
	messages := make(chan int)
	quit := make(chan struct{})

	go func() {
		var i int
		for {
			select {
			case messages <- i:
				i++
				fmt.Println("sent message")
			case <-quit:
				fmt.Println("done")
				close(messages)
				return
			}
		}
	}()

	for i := 0; i < 5; i++ {
		msg := <-messages
		log.Println(msg)
	}

	quit <- struct{}{}
	close(quit)
	time.Sleep(time.Second)
}
