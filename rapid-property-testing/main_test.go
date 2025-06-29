package main

import (
	"fmt"
	"sort"
	"testing"

	"pgregory.net/rapid"
)

func TestSortStrings(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		s := rapid.SliceOf(rapid.String()).Draw(t, "s")
		sort.Strings(s)
		if !sort.StringsAreSorted(s) {
			t.Fatalf("unsorted after sort: %v", s)
		}
	})
}

func testParseDate(t *rapid.T) {
	y := rapid.IntRange(0, 9999).Draw(t, "y")
	m := rapid.IntRange(1, 12).Draw(t, "m")
	d := rapid.IntRange(1, 31).Draw(t, "d")

	s := fmt.Sprintf("%04d-%02d-%02d", y, m, d)

	y_, m_, d_, err := ParseDate(s)
	if err != nil {
		t.Fatalf("failed to parse date %q: %v", s, err)
	}

	if y != y_ || m != m_ || d != d_ {
		t.Fatalf("expected %04d-%02d-%02d, got %04d-%02d-%02d", y, m, d, y_, m_, d_)
	}
}

func TestParseDate(t *testing.T) {
	rapid.Check(t, testParseDate)
}

func testQueue(t *rapid.T) {
	n := rapid.IntRange(1, 1000).Draw(t, "n")
	q := NewQueue(n)
	var state []int

	t.Repeat(map[string]func(*rapid.T){
		"get": func(t *rapid.T) {
			if q.Size() == 0 {
				t.Skip("queue is empty, cannot get")
			}

			i, _ := q.Get()
			if i != state[0] {
				t.Fatalf("expected Get() to return %d, got %d", state[0], i)
			}

			state = state[1:]
		},
		"put": func(t *rapid.T) {
			i := rapid.Int().Draw(t, "i")
			ok := q.Put(i)
			if ok {
				state = append(state, i)
			}
		},
		"": func(t *rapid.T) {
			if q.Size() != len(state) {
				t.Fatalf("expected queue size %d, got %d", len(state), q.Size())
			}
		},
	})
}

func TestQueue(t *testing.T) {
	rapid.Check(t, testQueue)
}
