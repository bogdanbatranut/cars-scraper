package main

import (
	"fmt"
	"time"
)

func main() {

	pub1 := make(chan string)
	pub2 := make(chan string)

	sub1 := make(chan string)
	sub2 := make(chan string)
	sub3 := make(chan string)

	go listen("Subscriber 1", sub1)
	go listen("Subscriber 2", sub2)
	go listen("Subscriber 3", sub3)

	subscriptions := map[chan string][]chan string{
		pub1: []chan string{sub1, sub2},
		pub2: []chan string{sub2, sub3},
	}

	go publish(subscriptions)

	pub1 <- "Hello, World!"
	pub2 <- "Hi, Universe!"
	pub1 <- "Goodbye, Cruel World!"

	subscriptions[pub1] = append(subscriptions[pub1], sub3)

	pub1 <- "Just kidding!"

	time.Sleep(2 * time.Second)
}

func listen(name string, subscriber chan string) {
	for {
		select {
		case message := <-subscriber:
			fmt.Printf("%q: %q\n", name, message)
		case <-time.After(time.Millisecond):
		}
	}
}

func publish(subscriptions map[chan string][]chan string) {
	for {
		for publisher, subscribers := range subscriptions {
			select {
			case message := <-publisher:
				for _, subscriber := range subscribers {
					subscriber <- message
				}
			case <-time.After(time.Millisecond):
			}
		}
	}
}
