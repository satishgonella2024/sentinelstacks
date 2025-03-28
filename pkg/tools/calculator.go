package tools

import (
	"fmt"
	"math"
)

// CalculatorTool provides basic arithmetic operations
type CalculatorTool struct{}

// ID returns the unique identifier for the calculator tool
func (c *CalculatorTool) ID() string {
	return "calculator"
}

// Name returns a user-friendly name
func (c *CalculatorTool) Name() string {
	return "Calculator"
}

// Description returns a detailed description
func (c *CalculatorTool) Description() string {
	return "Performs basic arithmetic operations: addition, subtraction, multiplication, division, power, square root"
}

// Version returns the semantic version
func (c *CalculatorTool) Version() string {
	return "0.1.0"
}

// ParameterSchema returns the JSON schema for parameters
func (c *CalculatorTool) ParameterSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"operation": map[string]interface{}{
				"type": "string",
				"enum": []string{"add", "subtract", "multiply", "divide", "power", "sqrt"},
				"description": "The arithmetic operation to perform",
			},
			"a": map[string]interface{}{
				"type": "number",
				"description": "First operand (not used for sqrt)",
			},
			"b": map[string]interface{}{
				"type": "number",
				"description": "Second operand (not used for sqrt)",
			},
		},
		"required": []string{"operation"},
	}
}

// Execute runs the calculator with the provided parameters
func (c *CalculatorTool) Execute(params map[string]interface{}) (interface{}, error) {
	operation, ok := params["operation"].(string)
	if !ok {
		return nil, fmt.Errorf("operation parameter is required")
	}

	switch operation {
	case "sqrt":
		a, ok := params["a"].(float64)
		if !ok {
			return nil, fmt.Errorf("parameter 'a' must be a number")
		}
		if a < 0 {
			return nil, fmt.Errorf("cannot take square root of negative number")
		}
		return math.Sqrt(a), nil

	case "add", "subtract", "multiply", "divide", "power":
		a, okA := params["a"].(float64)
		b, okB := params["b"].(float64)
		
		if !okA || !okB {
			return nil, fmt.Errorf("parameters 'a' and 'b' must be numbers")
		}

		switch operation {
		case "add":
			return a + b, nil
		case "subtract":
			return a - b, nil
		case "multiply":
			return a * b, nil
		case "divide":
			if b == 0 {
				return nil, fmt.Errorf("division by zero is not allowed")
			}
			return a / b, nil
		case "power":
			return math.Pow(a, b), nil
		}
	}

	return nil, fmt.Errorf("unsupported operation: %s", operation)
}
