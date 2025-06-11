package main

import (
	"fmt"
	"net/http"
	"os"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
)

func main() {
	logger := log.NewLogfmtLogger(os.Stderr)

	var svc StringService
	svc = stringService{}
	svc = loggingMiddleware{logger, svc}

	uppercaseHandler := httptransport.NewServer(
		makeUppercaseEndpoint(svc),
		decodeUppercaseRequest,
		encodeResponse,
	)
	countHandler := httptransport.NewServer(
		makeCountEndpoint(svc),
		decodeCountRequest,
		encodeResponse,
	)

	http.Handle("/uppercase", uppercaseHandler)
	http.Handle("/count", countHandler)

	fmt.Println("Server starting on http://localhost:8080")
	logger.Log(http.ListenAndServe(":8080", nil))
}
