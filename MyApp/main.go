package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	// Welcome! This file contains several examples of concurrency in Go.
	// Uncomment the function you want to run to see it in action.

	// demoBasicWaitGroup()
	// demoMultipleGoroutines()
	// demoChannels()
	demoContextTimeout()
}

// demoBasicWaitGroup demonstrates the basic usage of sync.WaitGroup
// to wait for a single goroutine to finish.
func demoBasicWaitGroup() {
	fmt.Println("--- Running Basic WaitGroup Demo ---")
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		process()
	}()

	wg.Wait()
	fmt.Println("Basic WaitGroup demo finished.")
}

// demoMultipleGoroutines demonstrates using a sync.WaitGroup
// to wait for multiple goroutines processing different data simultaneously.
func demoMultipleGoroutines() {
	fmt.Println("--- Running Multiple Goroutines Demo ---")
	var wg sync.WaitGroup
	clients := []string{"Client1", "Client2", "Client3"}

	for _, client := range clients {
		wg.Add(1)
		go processClient(client, &wg)
	}

	wg.Wait()
	fmt.Println("All goroutines finished executing.")
}

// demoChannels demonstrates passing data between goroutines using channels.
func demoChannels() {
	fmt.Println("--- Running Channels Demo ---")
	ch := make(chan int)

	// Start two goroutines that send their results to the channel
	go func(c int) {
		ch <- getResult(c)
	}(10)

	go func(c int) {
		ch <- getResult(c)
	}(20)

	// Receive the results from the channel
	fmt.Println("Result from a goroutine:", <-ch)
	fmt.Println("Result from another goroutine:", <-ch)
}

// demoContextTimeout demonstrates how to use context to cancel a running goroutine after a timeout.
func demoContextTimeout() {
	fmt.Println("--- Running Context Timeout Demo ---")
	// Create a context that will automatically cancel after 1 second
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		worker(ctx)
	}()

	wg.Wait()
	fmt.Println("Context timeout demo finished.")
}

// --- Helper Functions ---

func getResult(data int) int {
	return data
}

func process() {
	time.Sleep(2 * time.Second)
	fmt.Println("Processing...")
}

func processClient(client string, wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(1 * time.Second)
	fmt.Println("Processing client:", client)
}

func worker(ctx context.Context) {
	// Note on Context Cancellation:
	// If we used a simple time.Sleep(5 * time.Second) here, the goroutine
	// would not be aware of the context cancellation and would keep working.
	//
	// By using a select statement with time.After, we allow the goroutine
	// to listen for the ctx.Done() signal while waiting.

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Worker: context cancelled, stopping work.")
			return
		case <-time.After(5 * time.Second):
			fmt.Println("Worker: Working...")
		}
	}

	/* Timeline of execution:
		t=0
		│
		├── ctx.Done() → not ready
		│
		└── timer (time.After) → not ready
				│
				│ waiting in select block...
				│
		t=1s
		│
		└── ctx timeout fires
			↓
			ctx.Done() channel receives signal
			↓
			select wakes up on ctx.Done() case
			↓
			"Worker: context cancelled, stopping work."
	*/
}
