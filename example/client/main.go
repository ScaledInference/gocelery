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
// celery -A mjm_worker worker --loglevel=debug --without-heartbeat --without-mingle
// or for a worker here,
//  python3 -m celery -A example.worker:app -b amqp://localhost:5672 worker --loglevel=debug

func main() {
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
	asyncResult, err := celeryClient.Delay("example.worker.add", arg1, arg2)
	if err != nil {
		panic(err)
	}

	log.Println("=== main.go:34 task submitted")
	res, err := asyncResult.Get(1 * time.Second)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Result: %v of type: %v\n", res, reflect.TypeOf(res))
	}

}
