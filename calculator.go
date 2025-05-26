package main

import (
	"sync"
)

type OperationStrategy interface {
	Calculate(left, right int) int
}

type AddStrategy struct{}

func (a *AddStrategy) Calculate(left, right int) int { return left + right }

type SubtractStrategy struct{}

func (s *SubtractStrategy) Calculate(left, right int) int { return left - right }

type MultiplyStrategy struct{}

func (m *MultiplyStrategy) Calculate(left, right int) int { return left * right }

type Calculator struct {
	strategies   map[string]OperationStrategy
	commandByVar map[string]*Command
	variables    sync.Map
	printOrder   []string
}

func NewCalculator(commands []Command) *Calculator {
	c := &Calculator{
		strategies: map[string]OperationStrategy{
			"+": &AddStrategy{},
			"-": &SubtractStrategy{},
			"*": &MultiplyStrategy{},
		},
		commandByVar: make(map[string]*Command),
	}

	for i := range commands {
		if commands[i].Type == "calc" {
			c.commandByVar[commands[i].Var] = &commands[i]
		} else if commands[i].Type == "print" {
			c.printOrder = append(c.printOrder, commands[i].Var)
		}
	}

	return c
}

func (c *Calculator) evaluate(varName string) (int, bool) {
	if val, ok := c.variables.Load(varName); ok {
		return val.(int), true
	}

	cmd, exists := c.commandByVar[varName]
	if !exists {
		return 0, false
	}

	left, leftOk := c.getValue(cmd.Left)
	right, rightOk := c.getValue(cmd.Right)

	if !leftOk || !rightOk {
		return 0, false
	}

	strategy, ok := c.strategies[cmd.Op]
	if !ok {
		return 0, false
	}

	result := strategy.Calculate(left, right)
	c.variables.Store(varName, result)
	return result, true
}

func (c *Calculator) getValue(val interface{}) (int, bool) {
	switch v := val.(type) {
	case float64:
		return int(v), true
	case string:
		return c.evaluate(v)
	default:
		return 0, false
	}
}

func (c *Calculator) Process() Output {
	output := Output{}
	for _, varName := range c.printOrder {
		if val, ok := c.evaluate(varName); ok {
			output.Items = append(output.Items, Item{Var: varName, Value: val})
		}
	}
	return output
}
