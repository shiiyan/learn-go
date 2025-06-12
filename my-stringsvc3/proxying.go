package main

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"
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
