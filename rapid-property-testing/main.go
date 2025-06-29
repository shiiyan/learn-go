package main

import (
	"fmt"
	"strconv"
)

func ParseDate(s string) (int, int, int, error) {
	y, err := strconv.Atoi(s[0:4])
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to parse year from %q: %w", s, err)
	}

	m, err := strconv.Atoi(s[5:7])
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to parse month from %q: %w", s, err)
	}

	d, err := strconv.Atoi(s[8:10])
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to parse day from %q: %w", s, err)
	}

	return y, m, d, nil
}

type Queue struct {
	buf []int
	in  int
	out int
}

func NewQueue(n int) *Queue {
	return &Queue{
		buf: make([]int, n+1),
	}
}

func (q *Queue) Get() (int, bool) {
	if q.Size() == 0 {
		return 0, false
	}

	i := q.buf[q.out]
	q.out = (q.out + 1) % len(q.buf)
	return i, true
}

func (q *Queue) Put(i int) bool {
	if q.Size() == len(q.buf)-1 {
		return false
	}

	q.buf[q.in] = i
	q.in = (q.in + 1) % len(q.buf)
	return true
}

func (q *Queue) Size() int {
	return (q.in - q.out + len(q.buf)) % len(q.buf)
}
