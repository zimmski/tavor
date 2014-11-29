package primitives_test

import (
	"testing"

	. "github.com/zimmski/tavor/test/assert"

	"github.com/zimmski/tavor/token"
	"github.com/zimmski/tavor/token/lists"
	"github.com/zimmski/tavor/token/primitives"
)

func TestPointerList(t *testing.T) {
	a := primitives.NewRangeInt(4, 10)

	// empty pointers should have a nil token
	{
		var tok *token.Token

		o := primitives.NewEmptyPointer(tok)
		Nil(t, o.Get())

		err := o.Set(a)
		Nil(t, err)
		Equal(t, a, o.Get())

		var list *token.ListToken
		o = primitives.NewEmptyPointer(list)

		err = o.Set(a)
		NotNil(t, err)
		Nil(t, o.Get())

		l := lists.NewAll(a)

		err = o.Set(l)
		Nil(t, err)
		Equal(t, l, o.Get())
	}
}
