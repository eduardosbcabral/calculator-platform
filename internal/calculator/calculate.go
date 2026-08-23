package calculator

import (
	"errors"
	"math"
)

type operation string

const (
	add        operation = "add"
	subtract   operation = "subtract"
	multiply   operation = "multiply"
	divide     operation = "divide"
	power      operation = "power"
	squareRoot operation = "square_root"
	percentage operation = "percentage"
)

var (
	errUnsupportedOperation = errors.New("unsupported operation")
	errDivisionByZero       = errors.New("division by zero")
	errNegativeSquareRoot   = errors.New("square root of a negative number")
	errNonFiniteResult      = errors.New("result is not finite")
)

func calculate(op operation, left, right float64) (float64, error) {
	var result float64

	switch op {
	case add:
		result = left + right
	case subtract:
		result = left - right
	case multiply:
		result = left * right
	case divide:
		if right == 0 {
			return 0, errDivisionByZero
		}
		result = left / right
	case power:
		result = math.Pow(left, right)
	case squareRoot:
		if left < 0 {
			return 0, errNegativeSquareRoot
		}
		result = math.Sqrt(left)
	case percentage:
		result = left / 100 * right
	default:
		return 0, errUnsupportedOperation
	}

	if math.IsInf(result, 0) || math.IsNaN(result) {
		return 0, errNonFiniteResult
	}

	return result, nil
}
