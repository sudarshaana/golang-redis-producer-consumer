package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func main() {
	redisAdd := "localhost:6379"
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAdd,
	})

	// handle redis graceful shutdown
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt)
	doneChan := make(chan bool, 1)

	go func() {
		for {

			select {
			case <-stopChan:
				fmt.Println("Received interrupt signal, shutting down...")
				doneChan <- true
				return

			default:
				// generate a random value
				value := rand.Intn(1000)

				// push the value to redis list
				err := rdb.LPush(ctx, "queue", value).Err()
				if err != nil {
					log.Printf("Error while pushing to redis: %v", err)
				} else {
					log.Printf("pushed Value %d\n", value)
				}

				// sleep for a random time
				sleepDuration := time.Duration(1+rand.Intn(5)) * time.Second
				time.Sleep(sleepDuration)

			}

		}
	}()

	<-doneChan
	fmt.Println("Producer has shut down.")

}
