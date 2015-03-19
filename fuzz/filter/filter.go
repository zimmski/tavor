package filter

import (
	"fmt"
	"sort"

	"github.com/zimmski/container/list/linkedlist"

	"github.com/zimmski/tavor/token"
)

// Filter defines a fuzzing filter
type Filter interface {
	// Apply applies the fuzzing filter onto the token and returns a replacement token, or nil if there is no replacement.
	// If a fatal error is encountered the error return argument is not nil.
	Apply(tok token.Token) (token.Token, error)
}

var filterLookup = make(map[string]func() Filter)

// New returns a new fuzzing filter instance given the registered name of the filter.
// The error return argument is not nil, if the name does not exist in the registered fuzzing filter list.
func New(name string) (Filter, error) {
	filt, ok := filterLookup[name]
	if !ok {
		return nil, fmt.Errorf("unknown fuzzing filter %q", name)
	}

	return filt(), nil
}

// List returns a list of all registered fuzzing filter names.
func List() []string {
	keyFilterLookup := make([]string, 0, len(filterLookup))

	for key := range filterLookup {
		keyFilterLookup = append(keyFilterLookup, key)
	}

	sort.Strings(keyFilterLookup)

	return keyFilterLookup
}

// Register registers a fuzzing filter instance function with the given name.
func Register(name string, filt func() Filter) {
	if filt == nil {
		panic("register fuzzing filter is nil")
	}

	if _, ok := filterLookup[name]; ok {
		panic("fuzzing filter " + name + " already registered")
	}

	filterLookup[name] = filt
}

// ApplyFilters applies a set of filters onto a token.
// Filters are applied in the order in which they are given. If multiple filters are replacing the same token, only the first replacement will be applied.
// Filters are not applied onto filter generated tokens.
func ApplyFilters(filters []Filter, root token.Token) (token.Token, error) {
	type Pair struct {
		token  token.Token
		parent token.Token
	}

	var known = make(map[token.Token]struct{})

	var queue = linkedlist.New()

	queue.Unshift(&Pair{
		token:  root,
		parent: nil,
	})

	for !queue.Empty() {
		v, _ := queue.Shift()
		pair := v.(*Pair)

		tok := pair.token

		// only apply filters if the token is not from one
		if _, ok := known[tok]; !ok {
			// apply filters
			for i := range filters {
				replacement, err := filters[i].Apply(tok)
				if err != nil {
					return nil, fmt.Errorf("error in fuzzing filter %v: %s", filters[i], err)
				}

				// replace if there is something to replace with
				if replacement != nil {
					tok = replacement
					known[tok] = struct{}{}

					if pair.parent == nil {
						root = tok
					} else {
						if pTok, ok := pair.parent.(token.InternalReplace); ok {
							err := pTok.InternalReplace(pair.token, tok)
							if err != nil {
								return nil, err
							}
						} else {
							panic(fmt.Sprintf("Token %#v does not implement InternalReplace interface", pair.parent))
						}
					}
					break // stop filtering this token and go to the next one
				}
			}
		}

		// go deeper into the graph
		switch t := tok.(type) {
		case token.ForwardToken:
			c := t.InternalGet()

			queue.Unshift(&Pair{
				token:  c,
				parent: tok,
			})
		case token.ListToken:
			for i := t.InternalLen() - 1; i >= 0; i-- {
				c, _ := t.InternalGet(i)

				queue.Unshift(&Pair{
					token:  c,
					parent: tok,
				})
			}
		}
	}

	return root, nil
}
