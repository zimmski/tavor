package token

import (
	"github.com/zimmski/container/list/linkedlist"

	"github.com/zimmski/tavor/log"
)

// ResetResetTokens resets all tokens in the token graph that fullfill the ResetToken interface
func ResetResetTokens(root Token) {
	var queue = linkedlist.New()

	queue.Unshift(root)

	for !queue.Empty() {
		v, _ := queue.Shift()

		switch tok := v.(type) {
		case ResetToken:
			log.Debugf("reset %#v(%p)", tok, tok)

			tok.Reset()
		}

		if t, ok := v.(Follow); ok && !t.Follow() {
			continue
		}

		switch tok := v.(type) {
		case ForwardToken:
			if v := tok.Get(); v != nil {
				queue.Unshift(v)
			}
		case ListToken:
			for i := tok.Len() - 1; i >= 0; i-- {
				c, _ := tok.Get(i)

				queue.Unshift(c)
			}
		}
	}
}

// ResetScope resets all scopes of tokens in the token graph that fullfill the ScopeToken interface
func ResetScope(root Token) {
	log.Debug("start reseting scope")

	SetScope(root, make(map[string]Token))

	log.Debug("finished reseting scope")
}

// SetScope sets all scopes of tokens in the token graph that fullfill the ScopeToken interface
func SetScope(root Token, scope map[string]Token) {
	queue := linkedlist.New()

	type set struct {
		token Token
		scope map[string]Token
	}

	queue.Unshift(set{
		token: root,
		scope: scope,
	})

	for !queue.Empty() {
		v, _ := queue.Shift()
		s := v.(set)

		if t, ok := s.token.(ScopeToken); ok {
			log.Debugf("setScope %#v(%p) with %#v", t, t, s.scope)

			t.SetScope(s.scope)
		}

		if t, ok := s.token.(Follow); ok && !t.Follow() {
			continue
		}

		nScope := make(map[string]Token, len(s.scope))
		for k, v := range s.scope {
			nScope[k] = v
		}

		switch t := s.token.(type) {
		case ForwardToken:
			if v := t.Get(); v != nil {
				queue.Unshift(set{
					token: v,
					scope: nScope,
				})
			}
		case ListToken:
			for i := t.Len() - 1; i >= 0; i-- {
				c, _ := t.Get(i)

				queue.Unshift(set{
					token: c,
					scope: nScope,
				})
			}
		}
	}
}

// ResetInternalScope resets all scopes of interal tokens in the token graph that fullfill the ScopeToken interface
func ResetInternalScope(root Token) {
	log.Debug("start reseting scope")

	SetInternalScope(root, make(map[string]Token))

	log.Debug("finished reseting scope")
}

// SetInternalScope sets all scopes of internal tokens in the token graph that fullfill the ScopeToken interface
func SetInternalScope(root Token, scope map[string]Token) {
	queue := linkedlist.New()

	type set struct {
		token Token
		scope map[string]Token
	}

	queue.Unshift(set{
		token: root,
		scope: scope,
	})

	for !queue.Empty() {
		v, _ := queue.Shift()
		s := v.(set)

		if t, ok := s.token.(ScopeToken); ok {
			log.Debugf("setScope %#v(%p) with %#v", t, t, s.scope)

			t.SetScope(s.scope)
		}

		if t, ok := s.token.(Follow); ok && !t.Follow() {
			continue
		}

		nScope := make(map[string]Token, len(s.scope))
		for k, v := range s.scope {
			nScope[k] = v
		}

		switch t := s.token.(type) {
		case ForwardToken:
			if v := t.InternalGet(); v != nil {
				queue.Unshift(set{
					token: v,
					scope: nScope,
				})
			}
		case ListToken:
			for i := t.InternalLen() - 1; i >= 0; i-- {
				c, _ := t.InternalGet(i)

				queue.Unshift(set{
					token: c,
					scope: nScope,
				})
			}
		}
	}
}
