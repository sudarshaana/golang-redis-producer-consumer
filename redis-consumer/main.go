package main

import (
	"context"
	"fmt"
	"log"
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

	// handle graceful shutdown
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt)
	doneChan := make(chan bool, 1)

	log.Println("Consumer started...")

	go func() {
		for {
			select {
			case <-stopChan:
				log.Println("Received interrupt signal, shutting down consumer...")
				doneChan <- true
				return

			default:
				value, err := rdb.LPop(ctx, "queue").Result()

				if err == redis.Nil {
					// wait for some moments and then retry
					// log.Println("wait for some moments and then retry. redis.Nil")
					time.Sleep(1 * time.Millisecond)
					continue

				} else if err != nil {
					log.Printf("Error in reading value from redis: %v", err)
					time.Sleep(1 * time.Second)
					log.Printf("Error in reading value from redis. %v", err)
					continue

				}

				fmt.Printf("Fetched value: %s\n", value)

			}
		}

	}()
	<-doneChan
	fmt.Println("Consumer has shut down")
}
