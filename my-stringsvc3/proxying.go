package main

import (
	"context"
	"errors"
	"net/url"
	"strings"
	"time"

	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/lb"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
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

func proxyingMiddleware(ctx context.Context, proxyUrls string, logger log.Logger) ServiceMiddleware {
	if proxyUrls == "" {
		logger.Log("proxy_to", "none")
		return func(next StringService) StringService { return next }
	}

	var (
		proxyUrlList = split(proxyUrls)
		endPointer   sd.FixedEndpointer
	)

	var (
		qps         = 10
		maxAttempts = 3
		maxTime     = 250 * time.Millisecond
	)

	for _, proxyUrl := range proxyUrlList {
		var e endpoint.Endpoint
		e = makeUppercaseProxy(proxyUrl)
		e = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(e)
		e = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), qps))(e)
		endPointer = append(endPointer, e)
	}

	balancer := lb.NewRoundRobin(endPointer)
	retry := lb.Retry(maxAttempts, maxTime, balancer)

	return func(next StringService) StringService {
		return proxymw{ctx, next, retry}
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

func split(s string) []string {
	l := strings.Split(s, ",")
	for i := range l {
		l[i] = strings.TrimSpace(l[i])
	}

	return l
}
