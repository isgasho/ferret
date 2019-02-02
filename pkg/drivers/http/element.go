package http

import (
	"context"
	"encoding/json"
	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"hash/fnv"

	"github.com/MontFerret/ferret/pkg/drivers/common"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/PuerkitoBio/goquery"
)

type HTMLElement struct {
	selection *goquery.Selection
	attrs     *values.Object
	children  *values.Array
}

func NewHTMLElement(node *goquery.Selection) (*HTMLElement, error) {
	if node == nil {
		return nil, core.Error(core.ErrMissedArgument, "element selection")
	}

	return &HTMLElement{node, nil, nil}, nil
}

func (el *HTMLElement) MarshalJSON() ([]byte, error) {
	return json.Marshal(el.InnerText().String())
}

func (el *HTMLElement) Type() core.Type {
	return drivers.HTMLElementType
}

func (el *HTMLElement) String() string {
	return el.InnerHTML().String()
}

func (el *HTMLElement) Compare(other core.Value) int64 {
	switch other.Type() {
	case drivers.HTMLElementType:
		other := other.(drivers.HTMLNode)

		return el.InnerHTML().Compare(other.InnerHTML())
	default:
		return drivers.Compare(el.Type(), other.Type())
	}
}

func (el *HTMLElement) Unwrap() interface{} {
	return el.selection
}

func (el *HTMLElement) Hash() uint64 {
	str, err := el.selection.Html()

	if err != nil {
		return 0
	}

	h := fnv.New64a()

	h.Write([]byte(el.Type().Name()))
	h.Write([]byte(":"))
	h.Write([]byte(str))

	return h.Sum64()
}

func (el *HTMLElement) Copy() core.Value {
	c, _ := NewHTMLElement(el.selection.Clone())

	return c
}

func (el *HTMLElement) NodeType() values.Int {
	nodes := el.selection.Nodes

	if len(nodes) == 0 {
		return 0
	}

	return values.NewInt(common.ToHTMLType(nodes[0].Type))
}

func (el *HTMLElement) NodeName() values.String {
	return values.NewString(goquery.NodeName(el.selection))
}

func (el *HTMLElement) Length() values.Int {
	if el.children == nil {
		el.children = el.parseChildren()
	}

	return el.children.Length()
}

func (el *HTMLElement) Value() core.Value {
	val, ok := el.selection.Attr("value")

	if ok {
		return values.NewString(val)
	}

	return values.EmptyString
}

func (el *HTMLElement) InnerText() values.String {
	return values.NewString(el.selection.Text())
}

func (el *HTMLElement) InnerHTML() values.String {
	h, err := el.selection.Html()

	if err != nil {
		return values.EmptyString
	}

	return values.NewString(h)
}

func (el *HTMLElement) GetAttributes() core.Value {
	if el.attrs == nil {
		el.attrs = el.parseAttrs()
	}

	return el.attrs
}

func (el *HTMLElement) GetAttribute(name values.String) core.Value {
	v, ok := el.selection.Attr(name.String())

	if ok {
		return values.NewString(v)
	}

	return values.None
}

func (el *HTMLElement) GetChildNodes() core.Value {
	if el.children == nil {
		el.children = el.parseChildren()
	}

	return el.children
}

func (el *HTMLElement) GetChildNode(idx values.Int) core.Value {
	if el.children == nil {
		el.children = el.parseChildren()
	}

	return el.children.Get(idx)
}

func (el *HTMLElement) QuerySelector(selector values.String) core.Value {
	selection := el.selection.Find(selector.String())

	if selection == nil {
		return values.None
	}

	res, err := NewHTMLElement(selection)

	if err != nil {
		return values.None
	}

	return res
}

func (el *HTMLElement) QuerySelectorAll(selector values.String) core.Value {
	selection := el.selection.Find(selector.String())

	if selection == nil {
		return values.None
	}

	arr := values.NewArray(selection.Length())

	selection.Each(func(i int, selection *goquery.Selection) {
		el, err := NewHTMLElement(selection)

		if err == nil {
			arr.Push(el)
		}
	})

	return arr
}

func (el *HTMLElement) InnerHTMLBySelector(selector values.String) values.String {
	selection := el.selection.Find(selector.String())

	str, err := selection.Html()

	// TODO: log error
	if err != nil {
		return values.EmptyString
	}

	return values.NewString(str)
}

func (el *HTMLElement) InnerHTMLBySelectorAll(selector values.String) *values.Array {
	selection := el.selection.Find(selector.String())
	arr := values.NewArray(selection.Length())

	selection.Each(func(_ int, selection *goquery.Selection) {
		str, err := selection.Html()

		// TODO: log error
		if err == nil {
			arr.Push(values.NewString(str))
		}
	})

	return arr
}

func (el *HTMLElement) InnerTextBySelector(selector values.String) values.String {
	selection := el.selection.Find(selector.String())

	return values.NewString(selection.Text())
}

func (el *HTMLElement) InnerTextBySelectorAll(selector values.String) *values.Array {
	selection := el.selection.Find(selector.String())
	arr := values.NewArray(selection.Length())

	selection.Each(func(_ int, selection *goquery.Selection) {
		arr.Push(values.NewString(selection.Text()))
	})

	return arr
}

func (el *HTMLElement) CountBySelector(selector values.String) values.Int {
	selection := el.selection.Find(selector.String())

	if selection == nil {
		return values.ZeroInt
	}

	return values.NewInt(selection.Size())
}

func (el *HTMLElement) ExistsBySelector(selector values.String) values.Boolean {
	selection := el.selection.Closest(selector.String())

	if selection == nil {
		return values.False
	}

	return values.True
}

func (el *HTMLElement) GetIn(ctx context.Context, path []core.Value) (core.Value, error) {
	return common.GetIn(ctx, el, path)
}

func (el *HTMLElement) SetIn(ctx context.Context, path []core.Value, value core.Value) error {
	return common.SetIn(ctx, el, path, value)
}

func (el *HTMLElement) Iterate(ctx context.Context) (collections.CollectionIterator, error) {
	return common.NewIterator(el)
}

func (el *HTMLElement) parseAttrs() *values.Object {
	obj := values.NewObject()

	for _, name := range common.Attributes {
		val, ok := el.selection.Attr(name)

		if ok {
			obj.Set(values.NewString(name), values.NewString(val))
		}
	}

	return obj
}

func (el *HTMLElement) parseChildren() *values.Array {
	children := el.selection.Children()

	arr := values.NewArray(10)

	children.Each(func(i int, selection *goquery.Selection) {
		child, err := NewHTMLElement(selection)

		if err == nil {
			arr.Push(child)
		}
	})

	return arr
}