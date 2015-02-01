package token

import (
	"github.com/zimmski/container/list/linkedlist"

	"github.com/zimmski/tavor/log"
)

// LoopExists determines if a cycle exists in the internal token graph
func LoopExists(root Token) bool {
	lookup := make(map[Token]struct{})
	queue := linkedlist.New()

	queue.Unshift(root)

	for !queue.Empty() {
		v, _ := queue.Shift()
		t, _ := v.(Token)

		lookup[t] = struct{}{}

		switch tok := t.(type) {
		case PointerToken:
			if v := tok.InternalGet(); v != nil {
				if _, ok := lookup[v]; ok {
					log.Debugf("found a loop through (%p)%#v", t, t)

					return true
				}

				queue.Unshift(v)
			}
		case ForwardToken:
			if v := tok.InternalGet(); v != nil {
				queue.Unshift(v)
			}
		case ListToken:
			for i := tok.InternalLen() - 1; i >= 0; i-- {
				c, _ := tok.InternalGet(i)

				queue.Unshift(c)
			}
		}
	}

	return false
}
