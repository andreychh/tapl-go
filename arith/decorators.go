package arith

import "errors"

// MultiStepTerm wraps a term to perform full evaluation to normal form
// through multiple reduction steps, unlike single-step Evaluate() method.
type MultiStepTerm struct {
	// Origin is the initial term to evaluate
	Origin Term
}

// Evaluate performs full evaluation of the origin term to normal form.
// Applies reduction rules sequentially until the term becomes a value
// or gets stuck.
//
// Returns the final term in normal form and nil error on success.
// Any error except ErrCannotEvaluate interrupts evaluation.
// ErrCannotEvaluate indicates normal form is reached.
func (m MultiStepTerm) Evaluate() (Term, error) {
	current, err := m.Origin.Evaluate()
	if err != nil {
		return nil, err
	}

	for {
		current, err = current.Evaluate()
		if err != nil {
			if errors.Is(err, ErrCannotEvaluate) {
				break
			}
			return nil, err
		}
	}

	return current, nil
}

// Format returns string representation of the origin term.
func (m MultiStepTerm) Format() string {
	return m.Origin.Format()
}

// isValue checks if the origin term is a value.
// Note: this checks the origin term, not the evaluation result.
func (m MultiStepTerm) isValue() bool {
	return m.Origin.isValue()
}

// isNumeric checks if the origin term is a numeric value.
// Note: this checks the origin term, not the evaluation result.
func (m MultiStepTerm) isNumeric() bool {
	return m.Origin.isNumeric()
}
