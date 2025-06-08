package main

import (
	"sync"
	"testing"
)

func TestAddStrategy(t *testing.T) {
	testTable := []struct {
		name        string
		left, right int
		expected    int
	}{
		{"positive", 7, 3, 10},
		{"negative", -1, -2, -3},
		{"zero", 0, 0, 0},
	}

	strategy := &AddStrategy{}
	for _, testCase := range testTable {
		result := strategy.Calculate(testCase.left, testCase.right)
		t.Logf("Calling AddStrategy(%v, %v), result %d\n", testCase.left, testCase.right, result)
		if result != testCase.expected {
			t.Errorf("Incorrect result. Expected %d, got %d", testCase.expected, result)
		}
	}
}

func TestSubtractStrategy(t *testing.T) {
	testTable := []struct {
		name        string
		left, right int
		expected    int
	}{
		{"positive", 7, 3, 4},
		{"negative", -1, -2, 1},
		{"zero", 0, 0, 0},
	}

	strategy := &SubtractStrategy{}
	for _, testCase := range testTable {
		result := strategy.Calculate(testCase.left, testCase.right)
		t.Logf("Calling SubtractStrategy(%v, %v), result %d\n", testCase.left, testCase.right, result)
		if result != testCase.expected {
			t.Errorf("Incorrect result. Expected %d, got %d", testCase.expected, result)
		}
	}
}

func TestMultiplyStrategy(t *testing.T) {
	testTable := []struct {
		name        string
		left, right int
		expected    int
	}{
		{"positive", 1, 2, 2},
		{"negative", -1, -2, 2},
		{"zero", 0, 1, 0},
	}

	strategy := &MultiplyStrategy{}
	for _, testCase := range testTable {
		result := strategy.Calculate(testCase.left, testCase.right)
		t.Logf("Calling MultiplyStrategy(%v, %v), result %d\n", testCase.left, testCase.right, result)
		if result != testCase.expected {
			t.Errorf("Incorrect result. Expected %d, got %d", testCase.expected, result)
		}
	}
}

func TestGetValue(t *testing.T) {
	calc := &Calculator{
		strategies:   make(map[string]OperationStrategy),
		commandByVar: make(map[string]*Command),
		variables:    sync.Map{},
	}

	calc.variables.Store("valid_var", 42)

	testTable := []struct {
		name        string
		input       interface{}
		expectedVal int
		expectedOk  bool
	}{
		{"number", float64(5), 5, true},
		{"string valid", "valid_var", 42, true},
		{"string invalid", "invalid_var", 0, false},
		{"bool", true, 0, false},
	}

	for _, testCase := range testTable {
		result, ok := calc.getValue(testCase.input)
		t.Logf("Calling getValue(%v), result (%d, %v)", testCase.input, result, ok)
		if result != testCase.expectedVal || ok != testCase.expectedOk {
			t.Errorf("Incorrect result. Expected (%d, %v), got (%d, %v)", testCase.expectedVal, testCase.expectedOk, result, ok)
		}
	}
}

func TestEvaluate(t *testing.T) {
	calc := &Calculator{
		strategies:   map[string]OperationStrategy{"+": &AddStrategy{}},
		commandByVar: make(map[string]*Command),
		variables:    sync.Map{},
	}

	calc.variables.Store("varA", 5)
	calc.variables.Store("varB", 3)
	calc.commandByVar["varC"] = &Command{Var: "varC", Op: "+", Left: "varA", Right: "varB"}

	testTable := []struct {
		varName       string
		expected      int
		expectSuccess bool
	}{
		{"varA", 5, true},
		{"varC", 8, true},
		{"varD", 0, false},
	}

	for _, testCase := range testTable {
		result, ok := calc.evaluate(testCase.varName)
		t.Logf("Calling evaluate(%v), result (%d, %v)", testCase.varName, result, ok)
		if ok != testCase.expectSuccess || (ok && result != testCase.expected) {
			t.Errorf("Incorrect result. Expected (%d, %v), got (%d, %v)", testCase.expected, testCase.expectSuccess, result, ok)
		}
	}
}
