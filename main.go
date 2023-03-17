package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

type Schedule struct {
	ID        string    `json:"id"`
	StartTime time.Time `json:"start_time"`
}

type Fooer interface {
	Foo(schedule Schedule)
	FooJoke(filename string)
	FooStderr()
}

type fooerImpl struct{}

func (f fooerImpl) Foo(schedule Schedule) {
	echoThis := fmt.Sprintf(`{"id: %s, scheduled_time: %s"}`, schedule.ID, schedule.StartTime.String())
	cmd := exec.Command("echo", "-n", echoThis)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("Error getting stdout pipe: %v\n", err)
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Printf("Error getting stderr pipe: %v\n", err)
		return
	}

	if err := cmd.Start(); err != nil {
		log.Printf("Error starting command: %v\n", err)
		return
	}

	se, err := io.ReadAll(stderr)
	if err != nil {
		log.Printf("Error reading stderr: %v\n", err)
	}
	fmt.Printf("stderr from echo => \n%s\n", se)

	so, err := io.ReadAll(stdout)
	if err != nil {
		log.Printf("Error reading stdout: %v\n", err)
	}
	fmt.Printf("stdout from echo => \n%s\n", so)

	if err := cmd.Wait(); err != nil {
		log.Printf("Error waiting for command: %v\n", err)
		return
	}

}

func (f fooerImpl) FooJoke(filename string) {

	cmd := exec.Command("python3", filename)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("Error getting stdout pipe: %v\n", err)
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Printf("Error getting stderr pipe: %v\n", err)
		return
	}

	if err := cmd.Start(); err != nil {
		log.Printf("Error starting command: %v\n", err)
		return
	}

	se, err := io.ReadAll(stderr)
	if err != nil {
		log.Printf("Error reading stderr: %v\n", err)
	}
	fmt.Printf("stderr => \n%s\n", se)

	so, err := io.ReadAll(stdout)
	if err != nil {
		log.Printf("Error reading stdout: %v\n", err)
	}
	fmt.Printf("stdout from python script => \n%s\n", so)

	if err := cmd.Wait(); err != nil {
		log.Printf("Error waiting for command: %v\n", err)
		return
	}

}

func (f fooerImpl) FooStderr() {

	cmd := exec.Command("python3", "./pyfiles/print_to_stdout_and_stderr.py")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("Error getting stdout pipe: %v\n", err)
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Printf("Error getting stderr pipe: %v\n", err)
		return
	}

	if err := cmd.Start(); err != nil {
		log.Printf("Error starting command: %v\n", err)
		return
	}

	se, err := io.ReadAll(stderr)
	if err != nil {
		log.Printf("Error reading stderr: %v\n", err)
	}
	fmt.Printf("stderr from python script => \n%s\n", se)

	so, err := io.ReadAll(stdout)
	if err != nil {
		log.Printf("Error reading stdout: %v\n", err)
	}
	fmt.Printf("stdout from python script => \n%s\n", so)

	if err := cmd.Wait(); err != nil {
		log.Printf("Error waiting for command: %v\n", err)
		return
	}

}

func main() {
	scheduleChan := make(chan Schedule)

	// Create a new cache
	c := NewCache[string, time.Time]()

	// Set a value in the cache
	c.Set("FOO", time.Now())

	// Get a value from the cache
	if value, ok := c.Get("FOO"); ok {
		fmt.Printf("cache has value: %s (willl be deleted next))\n", value)
	}

	// Delete a value from the cache
	c.Delete("FOO")

	http.HandleFunc("/api/v1/schedule/add", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		decoder := json.NewDecoder(r.Body)
		var schedule Schedule
		err := decoder.Decode(&schedule)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		scheduleChan <- schedule

		w.WriteHeader(http.StatusOK)
	})

	server := &http.Server{Addr: ":8080"}

	fooer := fooerImpl{}

	go func() {
		for schedule := range scheduleChan {
			c.Set(schedule.ID, schedule.StartTime)
			fooer.Foo(schedule)
			fooer.FooJoke("./pyfiles/joke_api.py")
			fooer.FooStderr()
		}
	}()

	go func() {
		log.Println("Starting server...")
		if err := server.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				log.Printf("Error starting server: %v\n", err)
			}
		}
	}()

	// Create a channel to receive OS signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Block until we receive a signal
	<-sigChan

	log.Println("Shutting down server...")

	// Create a context with a timeout of 5 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown the server gracefully
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Error shutting down server: %v\n", err)
	}

	log.Println("Server shut down successfully")

	log.Printf("CACHE HAS VALUES\n")
	for k, v := range c.data {
		fmt.Printf("\t%s, %s\n", k, v)
	}

}
