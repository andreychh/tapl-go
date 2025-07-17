package arith

import (
	"errors"
	"fmt"
)

var ErrCannotEvaluate = errors.New("cannot evaluate term")

type Term interface {
	Evaluate() (Term, error)
	Format() string
	isValue() bool
	isNumeric() bool
}

type True struct{}

func (t True) Evaluate() (Term, error) { return nil, ErrCannotEvaluate }
func (t True) Format() string          { return "true" }
func (t True) isValue() bool           { return true }
func (t True) isNumeric() bool         { return false }

type False struct{}

func (f False) Evaluate() (Term, error) { return nil, ErrCannotEvaluate }
func (f False) Format() string          { return "false" }
func (f False) isValue() bool           { return true }
func (f False) isNumeric() bool         { return false }

type Zero struct{}

func (z Zero) Evaluate() (Term, error) { return nil, ErrCannotEvaluate }
func (z Zero) Format() string          { return "0" }
func (z Zero) isValue() bool           { return true }
func (z Zero) isNumeric() bool         { return true }

type If struct {
	Cond Term
	Then Term
	Else Term
}

func (i If) Evaluate() (Term, error) {
	switch cond := i.Cond.(type) {
	case True:
		return i.Then, nil
	case False:
		return i.Else, nil
	default:
		evaluated, err := cond.Evaluate()
		if err != nil {
			return nil, err
		}

		return If{Cond: evaluated, Then: i.Then, Else: i.Else}, nil
	}
}

func (i If) Format() string {
	return fmt.Sprintf("if %s then %s else %s", i.Cond.Format(), i.Then.Format(), i.Else.Format())
}
func (i If) isValue() bool   { return false }
func (i If) isNumeric() bool { return false }

type Succ struct {
	Operand Term
}

func (s Succ) Evaluate() (Term, error) {
	if s.isValue() {
		return s, nil
	}

	evaluated, err := s.Operand.Evaluate()
	if err != nil {
		return nil, err
	}

	return Succ{Operand: evaluated}, nil
}

func (s Succ) Format() string {
	return fmt.Sprintf("succ %s", s.Operand.Format())
}
func (s Succ) isValue() bool   { return s.isNumeric() }
func (s Succ) isNumeric() bool { return s.Operand.isNumeric() }

type Pred struct {
	Operand Term
}

func (p Pred) Evaluate() (Term, error) {
	switch operand := p.Operand.(type) {
	case Zero:
		return Zero{}, nil
	case Succ:
		if operand.isNumeric() {
			return operand.Operand, nil
		}
	default:
		evaluated, err := operand.Evaluate()
		if err != nil {
			return nil, err
		}

		return Pred{Operand: evaluated}, nil
	}

	return nil, ErrCannotEvaluate
}

func (p Pred) Format() string {
	return fmt.Sprintf("pred %s", p.Operand.Format())
}
func (p Pred) isValue() bool   { return false }
func (p Pred) isNumeric() bool { return false }

type IsZero struct {
	Operand Term
}

func (i IsZero) Evaluate() (Term, error) {
	switch operand := i.Operand.(type) {
	case Zero:
		return True{}, nil
	case Succ:
		if operand.isNumeric() {
			return False{}, nil
		}
	default:
		evaluated, err := operand.Evaluate()
		if err != nil {
			return nil, err
		}

		return IsZero{Operand: evaluated}, nil
	}

	return nil, ErrCannotEvaluate
}

func (i IsZero) Format() string {
	return fmt.Sprintf("iszero %s", i.Operand.Format())
}
func (i IsZero) isValue() bool   { return false }
func (i IsZero) isNumeric() bool { return false }
