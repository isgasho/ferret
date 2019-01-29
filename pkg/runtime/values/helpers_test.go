package values_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type CustomValue struct {
	type_      core.Type
	properties map[core.Value]core.Value
}

func (t *CustomValue) MarshalJSON() ([]byte, error) {
	return nil, core.ErrNotImplemented
}

func (t *CustomValue) Type() core.Type {
	return t.type_
}

func (t *CustomValue) String() string {
	return ""
}

func (t *CustomValue) Compare(other core.Value) int64 {
	return other.Compare(t) * -1
}

func (t *CustomValue) Unwrap() interface{} {
	return t
}

func (t *CustomValue) Hash() uint64 {
	return 0
}

func (t *CustomValue) Copy() core.Value {
	return values.None
}

func (t *CustomValue) GetIn(ctx context.Context, path []core.Value) (core.Value, error) {
	if path == nil || len(path) == 0 {
		return values.None, nil
	}

	propKey := path[0]
	propValue, ok := t.properties[propKey]

	if !ok {
		return values.None, nil
	}

	if len(path) == 1 {
		return propValue, nil
	}

	return values.GetIn(context.Background(), propValue, path[1:])
}

func (t *CustomValue) SetIn(ctx context.Context, path []core.Value, value core.Value) error {
	if path == nil || len(path) == 0 {
		return nil
	}

	propKey := path[0]
	propValue, ok := t.properties[propKey]

	if !ok {
		return nil
	}

	if len(path) == 1 {
		t.properties[propKey] = value

		return nil
	}

	return values.SetIn(context.Background(), propValue, path[1:], value)
}

func TestHelpers(t *testing.T) {
	customType := core.NewType("custom")

	Convey("Helpers", t, func() {
		Convey("Getter", func() {
			Convey("It should get a value by a given path", func() {
				ct := &CustomValue{
					type_: customType,
					properties: map[core.Value]core.Value{
						values.NewString("foo"): values.NewInt(1),
						values.NewString("bar"): &CustomValue{
							type_: customType,
							properties: map[core.Value]core.Value{
								values.NewString("qaz"): values.NewInt(2),
							},
						},
					},
				}

				foo, err := values.GetIn(context.Background(), ct, []core.Value{
					values.NewString("foo"),
				})

				So(err, ShouldBeNil)
				So(foo, ShouldEqual, values.NewInt(1))

				qaz, err := values.GetIn(context.Background(), ct, []core.Value{
					values.NewString("bar"),
					values.NewString("qaz"),
				})

				So(err, ShouldBeNil)
				So(qaz, ShouldEqual, values.NewInt(2))
			})
		})

		Convey("Setter", func() {
			Convey("It should get a value by a given path", func() {
				ct := &CustomValue{
					type_: customType,
					properties: map[core.Value]core.Value{
						values.NewString("foo"): values.NewInt(1),
						values.NewString("bar"): &CustomValue{
							type_: customType,
							properties: map[core.Value]core.Value{
								values.NewString("qaz"): values.NewInt(2),
							},
						},
					},
				}

				err := values.SetIn(context.Background(), ct, []core.Value{
					values.NewString("foo"),
				}, values.NewInt(2))

				So(err, ShouldBeNil)
				So(ct.properties[values.NewString("foo")], ShouldEqual, values.NewInt(2))

				err = values.SetIn(context.Background(), ct, []core.Value{
					values.NewString("bar"),
					values.NewString("qaz"),
				}, values.NewString("foobar"))

				So(err, ShouldBeNil)

				qaz, err := values.GetIn(context.Background(), ct, []core.Value{
					values.NewString("bar"),
					values.NewString("qaz"),
				})

				So(err, ShouldBeNil)
				So(qaz, ShouldEqual, values.NewString("foobar"))
			})
		})
	})
}
