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
