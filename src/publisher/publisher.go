package publisher

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"fmt"

	"github.com/devopsgig/utilities/src/utilities"
)

// Meta hold the base fields with meta information for a publisher type
type Meta struct {
	UUID     string `json:"UUID"`
	TaskType string `json:"task"`
}

// Task interface set the contract for types
// which want to implement the publisher.Task interface
type Task interface {
	SetUUID(UUID string)
	SetTaskType()
	GetTaskType() string
}

const addr = "http://127.0.0.1"
const port = "8080"

// Send function receives tasks of type publisher.Task,
// send them to the respective endpoint in the REST API
// and transmit the server response through the response channel
func Send(task Task, respCh chan string) {
	// Get new UUID
	UUID, err := utilities.UUID()
	if err != nil {
		log.Fatal(err)
	}
	// Set UUID to task
	task.SetUUID(UUID)
	// Parse task to JSON format
	taskJSON, err := json.Marshal(task)
	if err != nil {
		log.Fatal(err)
	}
	// Get the url endpoint for the task type
	endpoint, err := endpoint(task.GetTaskType())
	if err != nil {
		log.Fatal(err)
	}
	// Form the URL
	URL := fmt.Sprintf("%s:%s/%s", addr, port, endpoint)
	// Send task in json format as POST request to REST API
	data := url.Values{"task": []string{string(taskJSON)}}
	resp, err := http.PostForm(URL, data)
	if err != nil {
		log.Fatal(err)
	}
	respCh <- fmt.Sprintf("Response: %d", resp.StatusCode)
}

// endpoint returns the endpoint for the taskType
func endpoint(taskType string) (string, error) {
	switch taskType {
	case "arithmetic":
		return "arith", nil
	}
	return "", fmt.Errorf("Task %s is not supported", taskType)
}
