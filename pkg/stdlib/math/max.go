package math

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// Max returns the greatest (arithmetic mean) of the values in array.
// @param array (Array) - Array of numbers.
// @returns (Float) - The greatest of the values in array.
func Max(_ context.Context, args ...core.Value) (core.Value, error) {
	var err error
	err = core.ValidateArgs(args, 1, 1)

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[0], types.Array)

	if err != nil {
		return values.None, err
	}

	arr := args[0].(*values.Array)

	if arr.Length() == 0 {
		return values.None, nil
	}

	var max float64

	arr.ForEach(func(value core.Value, idx int) bool {
		err = core.ValidateType(value, types.Float, types.Int)

		if err != nil {
			return false
		}

		fv := toFloat(value)

		if fv > max {
			max = fv
		}

		return true
	})

	if err != nil {
		return values.None, nil
	}

	return values.NewFloat(max), nil
}
