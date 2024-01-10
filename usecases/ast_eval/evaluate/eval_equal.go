package evaluate

import (
	"context"
	"math"

	"github.com/checkmarble/marble-backend/models/ast"
	"github.com/cockroachdb/errors"
)

type Equal struct{}

const floatEqualityThreshold = 1e-8

func (f Equal) Evaluate(ctx context.Context, arguments ast.Arguments) (any, []error) {

	leftAny, rightAny, err := leftAndRight(arguments.Args)
	if err != nil {
		return MakeEvaluateError(err)
	}

	if left, right, errs := adaptLeftAndRight(leftAny, rightAny, adaptArgumentToString); len(errs) == 0 {
		return MakeEvaluateResult(left == right)
	}

	if left, right, errs := adaptLeftAndRight(leftAny, rightAny, adaptArgumentToBool); len(errs) == 0 {
		return MakeEvaluateResult(left == right)
	}

	if left, right, errs := adaptLeftAndRight(leftAny, rightAny, promoteArgumentToInt64); len(errs) == 0 {
		return MakeEvaluateResult(left == right)
	}

	if left, right, errs := adaptLeftAndRight(leftAny, rightAny, promoteArgumentToFloat64); len(errs) == 0 {
		return MakeEvaluateResult(math.Abs(left-right) < floatEqualityThreshold)
	}

	if left, right, errs := adaptLeftAndRight(leftAny, rightAny, adaptArgumentToTime); len(errs) == 0 {
		return MakeEvaluateResult(left.Equal(right))
	}

	return MakeEvaluateError(errors.Wrap(ast.ErrArgumentInvalidType, "all arguments to Equal Evaluate must be string, boolean, time, int or float"))
}
