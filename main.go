package main

import (
	"fmt"
	"log"
	"main/http"
	"main/internal/queue"
	"main/internal/router"
	"main/internal/service"
	"main/internal/utils"
	goHttp "net/http"
)

// program can process command line arguments with following flags:
// '-p' - sets the port of the service (default: 9818)
// '-t' - sets the default reading timeout value in seconds (default: 10)
// '-lt' - sets the topic(queue) limit (default: 10)
// '-lm' - sets the messages limit for each topic(queue) (default: 10)
func main() {

	// init options for the program with defaults
	port := utils.GetCLArgWithDefault("-p", "9818")
	defaultTimeout := utils.GetCLArgWithDefaultAsInt("-t", 10)
	messagesLimit := utils.GetCLArgWithDefaultAsInt("-lm", 10)
	topicsLimit := utils.GetCLArgWithDefaultAsInt("-lt", 10)

	// init queue
	q := queue.New(messagesLimit, topicsLimit)

	// init service based on the queue
	s, err := service.New(q, defaultTimeout)
	if err != nil {
		log.Panic(err)
	}

	// init handlers based on the service
	readFromQueue := http.ReadFromQueue(s)
	writeToQueue := http.WriteToQueue(s)

	r := router.New()

	r.Handle(goHttp.MethodGet, "/queue/{topic}", readFromQueue)
	r.Handle(goHttp.MethodPut, "/queue/{topic}", writeToQueue)

	addr := fmt.Sprintf("localhost:%s", port)

	if err := goHttp.ListenAndServe(addr, r); err != nil {
		log.Println(err)
	}
}
