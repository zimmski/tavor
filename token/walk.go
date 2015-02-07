package token

import (
	"github.com/zimmski/container/list/linkedlist"
)

// Walk traverses a token graph beginning from the given token and calls for every newly visited token the given function.
// A depth-first algorithm is used to traverse the graph. If the given walk function returns an error, the whole walk process ends by returning the error back to the caller
func Walk(root Token, walkFunc func(tok Token) error) error {
	queue := linkedlist.New()

	queue.Unshift(root)

	walked := make(map[Token]struct{})

	for !queue.Empty() {
		v, _ := queue.Shift()
		tok := v.(Token)

		if err := walkFunc(tok); err != nil {
			return err
		}

		switch t := tok.(type) {
		case ForwardToken:
			if v := t.Get(); v != nil {
				if _, ok := walked[v]; !ok {
					queue.Unshift(v)
				}
			}
		case ListToken:
			for i := t.Len() - 1; i >= 0; i-- {
				c, _ := t.Get(i)

				if _, ok := walked[c]; !ok {
					queue.Unshift(c)
				}
			}
		}
	}

	return nil
}

// WalkInternal traverses a internal token graph beginning from the given token and calls for every newly visited token the given function.
// A depth-first algorithm is used to traverse the graph. If the given walk function returns an error, the whole walk process ends by returning the error back to the caller
func WalkInternal(root Token, walkFunc func(tok Token) error) error {
	queue := linkedlist.New()

	queue.Unshift(root)

	walked := make(map[Token]struct{})

	for !queue.Empty() {
		v, _ := queue.Shift()
		tok := v.(Token)

		if err := walkFunc(tok); err != nil {
			return err
		}

		if t, ok := v.(Follow); ok && !t.Follow() {
			continue
		}

		switch t := tok.(type) {
		case ForwardToken:
			if v := t.InternalGet(); v != nil {
				if _, ok := walked[v]; !ok {
					queue.Unshift(v)
				}
			}
		case ListToken:
			for i := t.InternalLen() - 1; i >= 0; i-- {
				c, _ := t.InternalGet(i)

				if _, ok := walked[c]; !ok {
					queue.Unshift(c)
				}
			}
		}
	}

	return nil
}

// WalkInternalTail traverses a internal token graph beginning from the given token and calls for every newly visited token the given function after it has traversed all children.
// A depth-first algorithm is used to traverse the graph. If the given walk function returns an error, the whole walk process ends by returning the error back to the caller
func WalkInternalTail(root Token, walkFunc func(tok Token) error) error {
	if t, ok := root.(Follow); !ok || t.Follow() {
		switch t := root.(type) {
		case ForwardToken:
			if v := t.InternalGet(); v != nil {
				if err := WalkInternalTail(v, walkFunc); err != nil {
					return err
				}
			}
		case ListToken:
			for i := 0; i < t.InternalLen(); i++ {
				c, _ := t.InternalGet(i)

				if err := WalkInternalTail(c, walkFunc); err != nil {
					return err
				}
			}
		}
	}

	if err := walkFunc(root); err != nil {
		return err
	}

	return nil
}

// ReleaseTokens traverses the token graph and calls Release for every release token
func ReleaseTokens(root Token) {
	_ = Walk(root, func(tok Token) error {
		if t, ok := tok.(ReleaseToken); ok {
			t.Release()
		}

		return nil
	})
	_ = WalkInternal(root, func(tok Token) error {
		if t, ok := tok.(ReleaseToken); ok {
			t.Release()
		}

		return nil
	})
}
