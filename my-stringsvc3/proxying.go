package main

import (
	"context"
	"errors"
	"net/url"
	"strings"

	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/sony/gobreaker"
)

type proxymw struct {
	ctx       context.Context
	next      StringService
	uppercase endpoint.Endpoint
}

func (mw proxymw) Uppercase(s string) (string, error) {
	response, err := mw.uppercase(mw.ctx, uppercaseRequest{S: s})
	if err != nil {
		return "", err
	}

	resp := response.(uppercaseResponse)
	if resp.Err != "" {
		return resp.V, errors.New(resp.Err)
	}

	return resp.V, nil
}

func (mw proxymw) Count(s string) int {
	return mw.next.Count(s)
}

func proxyingMiddleware(ctx context.Context, proxyUrl string) ServiceMiddleware {
	if proxyUrl == "" {
		return func(next StringService) StringService {
			return next
		}
	}

	e := makeUppercaseProxy(proxyUrl)
	e = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(e)


	return func(next StringService) StringService {
		return proxymw{ctx, next, e}
	}
}

func makeUppercaseProxy(proxyURL string) endpoint.Endpoint {
	if !strings.HasPrefix(proxyURL, "http") {
		proxyURL = "http://" + proxyURL
	}

	u, err := url.Parse(proxyURL)
	if err != nil {
		panic(err)
	}

	if u.Path == "" {
		u.Path = "/uppercase"
	}

	return httptransport.NewClient(
		"GET",
		u,
		encodeRequest,
		decodeUppercaseResponse,
	).Endpoint()
}
