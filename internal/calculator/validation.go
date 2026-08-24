package calculator

import "net/http"

func (input calculationRequest) operands() (float64, float64, *requestError) {
	if input.Left == nil {
		return 0, 0, requestFailure(http.StatusBadRequest, "missing_operand", "The left operand is required.")
	}

	if input.Operation == squareRoot {
		if input.Right != nil {
			return 0, 0, requestFailure(
				http.StatusBadRequest,
				"unexpected_operand",
				"Square root accepts only the left operand.",
			)
		}
		return *input.Left, 0, nil
	}

	if !isBinary(input.Operation) {
		return 0, 0, requestFailure(
			http.StatusBadRequest,
			"unsupported_operation",
			"Choose a supported operation.",
		)
	}
	if input.Right == nil {
		return 0, 0, requestFailure(http.StatusBadRequest, "missing_operand", "The right operand is required.")
	}

	return *input.Left, *input.Right, nil
}

func isBinary(op operation) bool {
	switch op {
	case add, subtract, multiply, divide, power, percentage:
		return true
	default:
		return false
	}
}
