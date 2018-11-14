package main

import (
	"fmt"
	"log"
	"time"

	"github.com/ScaledInference/gocelery"
)

// Celery Task using args
func add(a int, b int) int {
	return a + b
}

// AddTask is Celery Task using kwargs
type AddTask struct {
	a int // x
	b int // y
}

// ParseKwargs parses kwargs
func (a *AddTask) ParseKwargs(kwargs map[string]interface{}) error {
	kwargA, ok := kwargs["x"]
	if !ok {
		return fmt.Errorf("undefined kwarg x")
	}
	kwargAFloat, ok := kwargA.(float64)
	if !ok {
		return fmt.Errorf("malformed kwarg x")
	}
	a.a = int(kwargAFloat)
	kwargB, ok := kwargs["y"]
	if !ok {
		return fmt.Errorf("undefined kwarg y")
	}
	kwargBFloat, ok := kwargB.(float64)
	if !ok {
		return fmt.Errorf("malformed kwarg y")
	}
	a.b = int(kwargBFloat)
	return nil
}

// RunTask executes actual task
func (a *AddTask) RunTask() (interface{}, error) {
	result := a.a + a.b
	return result, nil
}

func main() {
	// create broker and backend
	//celeryBroker := gocelery.NewRedisCeleryBroker("redis://localhost:6379")
	//celeryBackend := gocelery.NewRedisCeleryBackend("redis://localhost:6379")

	// AMQP example
	celeryBroker := gocelery.NewAMQPCeleryBroker("amqp://")
	celeryBackend := gocelery.NewAMQPCeleryBackend("amqp://")

	log.Println("=== main.go:59 configuring")
	// Configure with 2 celery workers
	celeryClient, err := gocelery.NewCeleryClient(celeryBroker, celeryBackend, 2)
	if err != nil {
		panic(err)
	}

	log.Println("=== main.go:59 registering")
	// worker.add name reflects "add" task method found in "worker.py"
	// this worker uses args
	celeryClient.Register("test_celery.tasks.longtime_add", add)
	celeryClient.Register("test_celery.tasks.add_reflect", &AddTask{})

	log.Println("=== main.go:73 running")
	// Start Worker - blocking method
	go celeryClient.StartWorker()
	// Wait 30 seconds and stop all workers
	time.Sleep(30 * time.Second)
	celeryClient.StopWorker()
}
