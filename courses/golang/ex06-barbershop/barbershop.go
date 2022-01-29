package main

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	CUTTING_TIME  = 20
	BARBERS_NUM   = 1
	HALL_SITS_NUM = 3
)

type Barber struct {
	val int
}

type Client struct {
	val int
}

func main() {
	clients := make(chan *Client)
	go clientMaker(clients)
	go BarberShop(clients)
	time.Sleep(2 * time.Second)
}

func clientMaker(clients chan *Client) {
	for {
		time.Sleep(time.Duration(rand.Intn(28)+7) * time.Millisecond)
		clients <- &Client{}
	}
}

func cutHair(barber *Barber, client *Client, finished chan *Barber) {
	// Cutting hair
	time.Sleep(CUTTING_TIME * time.Millisecond)
	finished <- barber
}

func BarberShop(clients <-chan *Client) {
	freeBarbers := []*Barber{}
	waitingClient := []*Client{}
	syncBarberChan := make(chan *Barber)

	//creating barbers
	for i := 0; i < BARBERS_NUM; i++ {
		freeBarbers = append(freeBarbers, &Barber{})
	}

	for {
		select {
		case client := <-clients:
			if len(freeBarbers) == 0 {
				if len(waitingClient) < HALL_SITS_NUM {
					// client is waiting in the hall
					waitingClient = append(waitingClient, client)
					fmt.Printf("Client is waiting in hall (%v)\n", len(waitingClient))
				} else {
					// No free space in the hall - bye client
					fmt.Println("No free space for client")
				}
			} else {
				barber := freeBarbers[0]
				freeBarbers = freeBarbers[1:]
				fmt.Println("Client goes to barber")
				go cutHair(barber, client, syncBarberChan)
			}
		// barber finishes work
		case barber := <-syncBarberChan:
			if len(waitingClient) > 0 {
				// get client from hall
				client := waitingClient[0]
				waitingClient = waitingClient[1:]
				fmt.Printf("Take client from room (%v)\n", len(waitingClient))
				go cutHair(barber, client, syncBarberChan)
			} else {
				// barber => go sleep
				fmt.Println("Barber idle")
				freeBarbers = append(freeBarbers, barber)
			}
		}
	}

}
