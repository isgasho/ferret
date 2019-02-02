package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// InnerTextAll returns an array of inner text of matched elements.
// @param doc (HTMLDocument|HTMLNode) - Parent document or element.
// @param selector (String) - String of CSS selector.
// @returns (String) - An array of inner text if any element found, otherwise empty array.
func InnerTextAll(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 2)

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[0], drivers.HTMLDocumentType, drivers.HTMLNodeType)

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[1], types.String)

	if err != nil {
		return values.None, err
	}

	doc := args[0].(drivers.HTMLNode)
	selector := args[1].(values.String)

	return doc.InnerTextBySelectorAll(selector), nil
}
