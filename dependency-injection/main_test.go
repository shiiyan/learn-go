package main

import (
	"testing"
	"time"
)

func TestNewGreeter_Grumpy(t *testing.T) {
	now = func() time.Time {
		return time.Unix(2, 0)
	}

	defer func() { now = time.Now }()

	g := NewGreeter("hello")

	if !g.Grumpy {
		t.Errorf("expected Grumpy=true, got Grumpy=%v", g.Grumpy)
	}
}
