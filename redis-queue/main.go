package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type Mydata struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	// check if successfully
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		fmt.Printf("Error connecting redis. %v", err)
		return
	}

	fmt.Println("Connected to redis: ", pong)

	// push data

	queueName := "myqueue"
	data := "mydata"

	err = rdb.LPush(ctx, queueName, data).Err()
	if err != nil {
		fmt.Printf("error pushing data to redis: %v", err)
		return
	}

	// complex data
	myData := Mydata{
		ID:    1,
		Name:  "nameSm",
		Email: "sudarshaana@gmail.com",
	}
	myDataJson, err := json.Marshal(myData)
	if err != nil {
		fmt.Println("Error serializing data:", err)
		return
	}

	err = rdb.LPush(ctx, queueName, myDataJson).Err()
	if err != nil {
		fmt.Printf("error pushing data to redis: %v", err)
		return
	}

	fmt.Println("Data pushed to queue: ", queueName)

	// read complex data
	// Read data from the Redis queue
	jsonData, err := rdb.RPop(ctx, queueName).Result()
	if err != nil {
		if err == redis.Nil {
			fmt.Println("Queue is empty")
		} else {
			fmt.Println("Error reading data from queue:", err)
		}
		return
	}

	var readDataJson Mydata
	err = json.Unmarshal([]byte(jsonData), &readDataJson)
	if err != nil {
		fmt.Println("Error deserializing data:", err)
		return
	}
	fmt.Println("Data read from queue:", readDataJson)

	// read data from redis queue
	readData, readErr := rdb.RPop(ctx, queueName).Result()
	if readErr != nil {
		if readErr == redis.Nil {
			fmt.Println("Queue is empty")
		} else {
			fmt.Println("Error reading data from queue:", err)
		}
	}

	fmt.Println()

	fmt.Println("Data read from queue:", readData)

}
