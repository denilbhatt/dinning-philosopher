// package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// )


// var messageChan chan string;
// type AddToQueueRequest struct {
// 	PhilosopherID int `json:"philosopherId"`
// }

// func StartListening() {
// 	messageChan := make(chan string)

// 	go func() {
// 		for {
// 			// Wait for a message to be received on the channel
// 			message := <-messageChan

// 			// Broadcast the message to all connected clients
// 			BroadcastMessage(message)
// 		}
// 	}()

// 	// Serve the SSE endpoint
// 	http.HandleFunc("/sse", func(w http.ResponseWriter, r *http.Request) {
// 		// Set the response headers
// 		w.Header().Set("Content-Type", "text/event-stream")
// 		w.Header().Set("Cache-Control", "no-cache")
// 		w.Header().Set("Connection", "keep-alive")

// 		// Create a channel to detect when the connection has been closed
// 		closeNotifyChan := w.(http.CloseNotifier).CloseNotify()

// 		// Start a loop to send SSE messages
// 		for {
// 			// Wait for a message to be received on the message channel or for the connection to be closed
// 			select {
// 			case message := <-messageChan:
// 				// Format the SSE message with the received message
// 				msg := fmt.Sprintf("data: %s\n\n", message)

// 				// Write the message to the response writer
// 				_, err := w.Write([]byte(msg))
// 				if err != nil {
// 					// If there's an error, the client has probably closed the connection
// 					return
// 				}

// 				// Flush the response writer to ensure that the message is sent immediately
// 				w.(http.Flusher).Flush()

// 			case <-closeNotifyChan:
// 				// If the connection has been closed, stop sending SSE messages
// 				return
// 			}
// 		}
// 	})

// 	http.HandleFunc("/add-to-queue", func(w http.ResponseWriter, r *http.Request) {
// 		// Parse JSON request body
// 	var req AddToQueueRequest
// 	err := json.NewDecoder(r.Body).Decode(&req)
// 	if err != nil {
// 		http.Error(w, "Invalid request body", http.StatusBadRequest)
// 		return
// 	}

// 	// Add philosopherId to queue
// 	HungryChan <- req.PhilosopherID

// 	// Return success response
// 	w.WriteHeader(http.StatusOK)
// 	})

// 	// Start the server
// 	http.ListenAndServe(":8080", nil)
// }

// func BroadcastMessage(message string) {
// 	messageChan <- message
// }
