package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Semaphore struct {
    count int
}

func (sem *Semaphore) wait(){
	for ;sem.count <= 0; {}
	sem.count -= 1
}

func (sem *Semaphore) signal(){
	sem.count += 1;
}

type Philosopher struct {
    id int  `json:"id"`
    state string    `json:"state"`
    leftFork, rightFork *Semaphore
}

var forks []*Semaphore
var philosophers []*Philosopher
var HungryChan chan int
var (msgChanMutex *Semaphore = &Semaphore{count: 1})

func makeHungry(phil *Philosopher) {
    fmt.Printf("Philosopher %d is hungry\n", phil.id)
    phil.state = "HUNGRY"
    messageChan <- (fmt.Sprintf("{\"id\":%d,\"state\":\"%s\"}", phil.id, phil.state))
    phil.leftFork.wait()
    phil.rightFork.wait()

    // time.Sleep(time.Millisecond * 500)
    phil.state = "EATING"
    fmt.Printf("Philosopher %d is eating\n", phil.id)
    messageChan <- (fmt.Sprintf("{\"id\":%d,\"state\":\"%s\"}", phil.id, phil.state))
    
	time.Sleep(time.Second * 5)

	phil.state = "THINKING"

    fmt.Printf("Philosopher %d is thinking\n", phil.id)
    messageChan <- (fmt.Sprintf("{\"id\":%d,\"state\":\"%s\"}", phil.id, phil.state))
    phil.leftFork.signal()
    phil.rightFork.signal()

    
}

func init() {
    forks = make([]*Semaphore, 5)
    philosophers = make([]*Philosopher, 5)

	for i := 0; i < 5; i++ {
        forks[i] = new(Semaphore)
		forks[i].count = 1
    }

	
    // id: 0-4
    for i := 0; i < 5; i++ {
        philosophers[i] = &Philosopher{
            id: i,
            state: "THINKING",
            leftFork: forks[i],
            rightFork: forks[(i+1)%5],
        }
    }
    fmt.Println(philosophers)
}


var messageChan chan string;
type AddToQueueRequest struct {
	PhilosopherID int `json:"philosopherId"`
}
func corsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Set CORS headers
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")

        // Handle preflight requests
        if r.Method == "OPTIONS" {
            return
        }

        // Call the next handler
        next.ServeHTTP(w, r)
    })
}

func StartListening() {
	messageChan = make(chan string)
	HungryChan = make(chan int)
    mux := http.NewServeMux()
	// Serve the SSE endpoint
	mux.HandleFunc("/sse", func(w http.ResponseWriter, r *http.Request) {
		// Set the response headers
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		// Create a channel to detect when the connection has been closed
		closeNotifyChan := w.(http.CloseNotifier).CloseNotify()

		// Start a loop to send SSE messages
		for {
			// Wait for a message to be received on the message channel or for the connection to be closed
			select {
			case message := <-messageChan:
				// Format the SSE message with the received message
                fmt.Println(message, " HERE")
				msg := fmt.Sprintf("data: %s\n\n", message)
            
				// Write the message to the response writer
				_, err := w.Write([]byte(msg))
				if err != nil {
					// If there's an error, the client has probably closed the connection
					return
				}

				// Flush the response writer to ensure that the message is sent immediately
				w.(http.Flusher).Flush()

			case <-closeNotifyChan:
				// If the connection has been closed, stop sending SSE messages
				return
			}
		}
	})

	mux.HandleFunc("/add-to-queue", func(w http.ResponseWriter, r *http.Request) {
        fmt.Println("Request to add to queue")
		// Parse JSON request body
        var req AddToQueueRequest
        fmt.Println(r.Body);
        err := json.NewDecoder(r.Body).Decode(&req)
        if err != nil {
            http.Error(w, "Invalid request body", http.StatusBadRequest)
            return
        }

        // Add philosopherId to queue
        fmt.Println(req.PhilosopherID, " added to queue")
        HungryChan <- req.PhilosopherID
       
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(fmt.Sprint("Added ", req.PhilosopherID)))
    })

    go func(){
        // for {
        //     select{
        //     case data :=<- HungryChan:
        //         fmt.Println(data," is here")
        //     }
        // }
        for{
            select {
                case philosopherId := <- HungryChan:
                    fmt.Println(philosopherId, " is hungry")
                    // wg.Add(1)
                    go func(phil *Philosopher) {
                        // defer wg.Done()
                        makeHungry(phil)
                    }(philosophers[philosopherId])
                }
        }
    }()
    

    // go func(){
    //     for{
    //         select{
    //         case msg := <- messageChan:
    //             fmt.Println("MSG: ", msg)
    //         }
    //     }
    // }()

    // Wrap the mux with the cors middleware
    corsHandler := corsMiddleware(mux)
	// Start the server
    fmt.Println("Server is listening...")
	http.ListenAndServe(":8080", corsHandler)
}

// func BroadcastMessage(message string) {
// 	messageChan <- message
// }

func main() {
    go StartListening()
    fmt.Println("program is running")
    // var wg sync.WaitGroup
    
        for{
            select {
                case philosopherId := <- HungryChan:
                    fmt.Println(philosopherId, " is hungry")
                    // wg.Add(1)
                    go func(phil *Philosopher) {
                        // defer wg.Done()
                        makeHungry(phil)
                    }(philosophers[philosopherId])
                }
        }
        // wg.Wait()
    }
