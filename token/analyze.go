package token

import (
	"github.com/zimmski/container/list/linkedlist"

	"github.com/zimmski/tavor/log"
	//"github.com/zimmski/tavor/token/primitives"
)

// LoopExists determines if a cycle exists in the internal token graph
func LoopExists(root Token) bool {
	lookup := make(map[Token]struct{})
	queue := linkedlist.New()

	queue.Push(root)

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

				queue.Push(v)
			}
		case ForwardToken:
			if v := tok.InternalGet(); v != nil {
				queue.Push(v)
			}
		case List:
			for i := 0; i < tok.InternalLen(); i++ {
				c, _ := tok.InternalGet(i)

				queue.Push(c)
			}
		}
	}

	return false
}
