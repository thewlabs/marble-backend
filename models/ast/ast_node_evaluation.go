package ast

type NodeEvaluation struct {
	ReturnValue any
	Errors      []error

	Children      []NodeEvaluation
	NamedChildren map[string]NodeEvaluation
}

func (root NodeEvaluation) AllErrors() (errs []error) {

	var addEvaluationErrors func(NodeEvaluation)

	addEvaluationErrors = func(child NodeEvaluation) {
		if child.Errors != nil {
			errs = append(errs, child.Errors...)
		}

		for _, child := range child.Children {
			addEvaluationErrors(child)
		}

		for _, child := range child.NamedChildren {
			addEvaluationErrors(child)
		}
	}

	addEvaluationErrors(root)
	return errs
}
