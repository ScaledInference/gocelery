package main

import (
	"fmt"
	"log"
	"math/rand"
	"reflect"
	"time"

	"github.com/ScaledInference/gocelery"
)

// Run Celery Worker First!
// celery -A worker worker --loglevel=debug --without-heartbeat --without-mingle

func main() {

	// create broker and backend
	//celeryBroker := gocelery.NewRedisCeleryBroker("redis://localhost:6379")
	//celeryBackend := gocelery.NewRedisCeleryBackend("redis://localhost:6379")

	// AMQP example
	celeryBroker := gocelery.NewAMQPCeleryBroker("amqp://")
	celeryBackend := gocelery.NewAMQPCeleryBackend("amqp://")

	// create client
	celeryClient, err := gocelery.NewCeleryClient(celeryBroker, celeryBackend, 0)
	if err != nil {
		panic(err)
	}

	r := rand.New(rand.NewSource(int64(time.Now().Unix())))
	arg1 := r.Intn(100)
	arg2 := r.Intn(100)

	log.Printf("=== main.go:34 Submitting task %v ,%v", arg1, arg2)
	asyncResult, err := celeryClient.Delay("server.tasks.longtime_add", arg1, arg2)
	if err != nil {
		panic(err)
	}

	log.Println("=== main.go:34 task submitted")
	res, err := asyncResult.Get(10 * time.Second)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Result: %v of type: %v\n", res, reflect.TypeOf(res))
	}

	// send task
	arg3 := r.Intn(100)
	arg4 := r.Intn(100)


	asyncResult, err = celeryClient.DelayKwargs("test_celery.tasks.add_reflect", map[string]interface{}{
		"x": arg3,
		"y": arg4,
	})
	if err != nil {
		panic(err)
	}

	// check if result is ready
	isReady, err := asyncResult.IsComplete()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Ready status: %v\n", isReady)

	// get result with 1s timeout
	res2, err := asyncResult.Get(10 * time.Second)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Result: %v of type: %v\n", res2, reflect.TypeOf(res2))
	}

}
