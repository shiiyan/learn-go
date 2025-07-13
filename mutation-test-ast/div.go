package main

import "errors"

var ErrDivideByZero = errors.New("cannot divide by zero")

func divide(divident, divisor int) (int, error) {
	// mutant 1: reverse if condition
	if divisor == 0 {
		return 0, ErrDivideByZero
	}

	// mutant 2: change the operator
	return divident / divisor, nil
}
