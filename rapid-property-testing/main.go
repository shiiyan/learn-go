package main

import (
	"fmt"
	"strconv"
)


func ParseDate(s string) (int, int, int, error) {
	y , err := strconv.Atoi(s[0:4])
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