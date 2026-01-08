package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Quak1/learn-http-go/internal/request"
	"github.com/Quak1/learn-http-go/internal/response"
	"github.com/Quak1/learn-http-go/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func handler(w *response.Writer, req *request.Request) {
	target := req.RequestLine.RequestTarget
	if target == "/yourproblem" {
		yourProblemHandler(w, req)
	} else if target == "/myproblem" {
		myProblemHandler(w, req)
	} else if strings.HasPrefix(target, "/httpbin") {
		proxyHandler(w, req)
	} else {
		successHandler(w, req)
	}
}

func myProblemHandler(w *response.Writer, req *request.Request) {
	body := `<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>
`

	w.WriteStatusLine(500)

	headers := response.GetDefaultHeaders(len(body))
	headers.Replace("Content-Type", "text/html")
	w.WriteHeaders(headers)

	w.WriteBody([]byte(body))
}

func yourProblemHandler(w *response.Writer, req *request.Request) {
	body := `<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>
`

	w.WriteStatusLine(400)

	headers := response.GetDefaultHeaders(len(body))
	headers.Replace("Content-Type", "text/html")
	w.WriteHeaders(headers)

	w.WriteBody([]byte(body))
}

func successHandler(w *response.Writer, req *request.Request) {
	body := `<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>
`

	w.WriteStatusLine(200)

	headers := response.GetDefaultHeaders(len(body))
	headers.Replace("Content-Type", "text/html")
	w.WriteHeaders(headers)

	w.WriteBody([]byte(body))
}

func proxyHandler(w *response.Writer, req *request.Request) {
	endpoint := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin/")
	url := "https://httpbin.org/" + endpoint
	resp, err := http.Get(url)
	if err != nil {
		yourProblemHandler(w, req)
		return
	}
	defer resp.Body.Close()

	w.WriteStatusLine(response.StatusOK)

	headers := response.GetDefaultHeaders(0)
	headers.Delete("Content-Length")
	headers.Set("Transfer-Econding", "chunked")
	w.WriteHeaders(headers)

	buff := make([]byte, 1024)
	n, err := resp.Body.Read(buff)
	for {
		n, err = resp.Body.Read(buff)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error: couldn't read response")
			break
		}
		w.WriteChunkedBody(buff[:n])
	}

	w.WriteChunkedBodyDone()
}
