package main

import "errors"

var ErrDivideByZero = errors.New("cannot divide by zero")

func divide(divident, divisor int) (int, error) {

	if divisor == 0 {
		return 0, ErrDivideByZero
	}

	return divident * divisor, nil
}
