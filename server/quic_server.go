package main

import (
	"context"

	"io"
	"log"
	"time"

	"distributed-code-executor/shared"
	"github.com/quic-go/quic-go"
)

func HandleConnection(conn *quic.Conn) {
	//close the connection and send the string msg to remote peer
	defer conn.CloseWithError(0, "Peer Candidate has closed the connection")
	log.Printf("New connection from %s", conn.RemoteAddr())
	for {
		stream, err := conn.AcceptStream(context.Background())

		if err != nil {
			log.Printf("Failed to accept stream: %v", err)
			return
		}

		go handleStream(stream)
	}
}

func handleStream(stream *quic.Stream) {
	defer stream.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := shared.ReadRequest(stream)

	if err != nil {
		log.Printf("Failed to read request: %v", err)
		sendErrorResponse(stream, "", "Failed to read request")
		return
	}

	log.Printf("Received request %s for %s code execution", req.ID, req.Language)

	resp := executeCode(ctx, &req)

	// Send response
	if err := shared.SendResponse(stream, *resp); err != nil {
		log.Printf("Failed to send response: %v", err)
	}
}

func sendErrorResponse(w io.Writer, id, errorMsg string) {
	resp := shared.CodeResponse{
		ID:        id,
		Error:     errorMsg,
		ExitCode:  -1,
		Timestamp: time.Now(),
	}
	shared.SendResponse(w, resp)
}
