package calculator

import (
	"errors"
	"math"
	"testing"
)

func TestCalculate(t *testing.T) {
	tests := []struct {
		name      string
		operation operation
		left      float64
		right     float64
		want      float64
		wantErr   error
	}{
		{name: "add", operation: add, left: 7, right: 5, want: 12},
		{name: "subtract", operation: subtract, left: 7, right: 5, want: 2},
		{name: "multiply", operation: multiply, left: 7, right: 5, want: 35},
		{name: "divide", operation: divide, left: 7, right: 2, want: 3.5},
		{name: "power", operation: power, left: 2, right: 8, want: 256},
		{name: "zero to zero", operation: power, left: 0, right: 0, want: 1},
		{name: "square root", operation: squareRoot, left: 81, want: 9},
		{name: "percentage", operation: percentage, left: 15, right: 200, want: 30},
		{name: "division by zero", operation: divide, left: 7, wantErr: errDivisionByZero},
		{name: "negative square root", operation: squareRoot, left: -1, wantErr: errNegativeSquareRoot},
		{name: "overflow", operation: multiply, left: math.MaxFloat64, right: 2, wantErr: errNonFiniteResult},
		{name: "invalid power", operation: power, left: -1, right: 0.5, wantErr: errNonFiniteResult},
		{name: "unsupported", operation: "modulo", wantErr: errUnsupportedOperation},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := calculate(test.operation, test.left, test.right)
			if !errors.Is(err, test.wantErr) {
				t.Fatalf("calculate() error = %v, want %v", err, test.wantErr)
			}
			if err == nil && got != test.want {
				t.Fatalf("calculate() = %v, want %v", got, test.want)
			}
		})
	}
}

func BenchmarkCalculate(b *testing.B) {
	for b.Loop() {
		_, _ = calculate(power, 12.5, 3.25)
	}
}
