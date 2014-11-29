package token

import (
	"github.com/zimmski/container/list/linkedlist"

	"github.com/zimmski/tavor/log"
)

// ResetResetTokens resets all tokens in the token graph that fullfill the ResetToken interface
func ResetResetTokens(root Token) {
	var queue = linkedlist.New()

	queue.Push(root)

	for !queue.Empty() {
		v, _ := queue.Shift()

		switch tok := v.(type) {
		case ResetToken:
			log.Debugf("reset %#v(%p)", tok, tok)

			tok.Reset()
		}

		switch tok := v.(type) {
		case ForwardToken:
			if v := tok.Get(); v != nil {
				queue.Push(v)
			}
		case ListToken:
			for i := 0; i < tok.Len(); i++ {
				c, _ := tok.Get(i)
				queue.Push(c)
			}
		}
	}
}

// ResetScope resets all scopes of tokens in the token graph that fullfill the ScopeToken interface
func ResetScope(root Token) {
	SetScope(root, make(map[string]Token))
}

// SetScope sets all scopes of tokens in the token graph that fullfill the ScopeToken interface
func SetScope(root Token, scope map[string]Token) {
	queue := linkedlist.New()

	type set struct {
		token Token
		scope map[string]Token
	}

	queue.Push(set{
		token: root,
		scope: scope,
	})

	for !queue.Empty() {
		v, _ := queue.Shift()
		s := v.(set)

		if t, ok := s.token.(ScopeToken); ok {
			log.Debugf("setScope %#v(%p)", t, t)

			t.SetScope(s.scope)
		}

		nScope := make(map[string]Token, len(s.scope))
		for k, v := range s.scope {
			nScope[k] = v
		}

		switch t := s.token.(type) {
		case ForwardToken:
			if v := t.Get(); v != nil {
				queue.Push(set{
					token: v,
					scope: nScope,
				})
			}
		case ListToken:
			for i := 0; i < t.Len(); i++ {
				c, _ := t.Get(i)

				queue.Push(set{
					token: c,
					scope: nScope,
				})
			}
		}
	}
}

// SetInternalScope sets all scopes of internal tokens in the token graph that fullfill the ScopeToken interface
func SetInternalScope(root Token, scope map[string]Token) {
	queue := linkedlist.New()

	queue.Push(root)

	for !queue.Empty() {
		tok, _ := queue.Shift()

		if t, ok := tok.(ScopeToken); ok {
			log.Debugf("setScope %#v(%p)", t, t)

			t.SetScope(scope)
		}

		switch t := tok.(type) {
		case ForwardToken:
			if v := t.InternalGet(); v != nil {
				queue.Push(v)
			}
		case ListToken:
			for i := 0; i < t.InternalLen(); i++ {
				c, _ := t.InternalGet(i)

				queue.Push(c)
			}
		}
	}
}
