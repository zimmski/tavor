package lists

import (
	"bytes"
	"math"
	"strconv"

	"github.com/zimmski/tavor/token"
	"github.com/zimmski/tavor/token/primitives"
)

// Repeat implements a list token which repeats a referenced token by a given range
type Repeat struct {
	from  token.Token
	to    token.Token
	token token.Token
	value []token.Token

	reducing              bool
	reducingOriginalValue []token.Token
}

// NewRepeat returns a new instance of a Repeat token referencing the given token and the given integer range
func NewRepeat(tok token.Token, from int64, to int64) *Repeat {
	return NewRepeatWithTokens(tok, primitives.NewConstantInt(int(from)), primitives.NewConstantInt(int(to)))
}

// NewRepeatWithTokens returns a new instance of a Repeat token referencing the given token and the given token range
// The tokens of the given range must return a valid integer values.
func NewRepeatWithTokens(tok token.Token, from token.Token, to token.Token) *Repeat {
	iFrom, err := strconv.Atoi(from.String())
	if err != nil {
		panic(err) // TODO
	}

	l := &Repeat{
		from:  from,
		to:    to,
		token: tok,
		value: make([]token.Token, iFrom),
	}

	for i := range l.value {
		l.value[i] = tok.Clone()
	}

	return l
}

// From returns the from value of the repeat range
func (l *Repeat) From() int64 {
	iFrom, err := strconv.Atoi(l.from.String())
	if err != nil {
		panic(err) // TODO
	}

	return int64(iFrom)
}

// To returns the to value of the repeat range
func (l *Repeat) To() int64 {
	iTo, err := strconv.Atoi(l.to.String())
	if err != nil {
		panic(err) // TODO
	}

	return int64(iTo)
}

// Token interface methods

// Clone returns a copy of the token and all its children
func (l *Repeat) Clone() token.Token {
	c := Repeat{
		from:  l.from,
		to:    l.to,
		token: l.token.Clone(),
		value: make([]token.Token, len(l.value)),
	}

	for i, tok := range l.value {
		c.value[i] = tok.Clone()
	}

	return &c
}

// Parse tries to parse the token beginning from the current position in the parser data.
// If the parsing is successful the error argument is nil and the next current position after the token is returned.
func (l *Repeat) Parse(pars *token.InternalParser, cur int) (int, []error) {
	var toks []token.Token

	i := 1

	for i <= int(l.From()) {
		tok := l.token.Clone()

		nex, errs := tok.Parse(pars, cur)

		if len(errs) > 0 {
			return cur, errs
		}

		cur = nex
		toks = append(toks, tok)

		i++
	}

	for i <= int(l.To()) {
		tok := l.token.Clone()

		nex, errs := tok.Parse(pars, cur)

		if len(errs) > 0 {
			break
		}

		cur = nex
		toks = append(toks, tok)

		i++
	}

	l.value = toks

	return cur, nil
}

func (l *Repeat) permutation(i uint) {
	toks := make([]token.Token, int(i)+int(l.From()))

	token.ReleaseTokens(l)

	for i := range toks {
		toks[i] = l.token.Clone()
	}

	l.value = toks
}

// Permutation sets a specific permutation for this token
func (l *Repeat) Permutation(i uint) error {
	permutations := l.Permutations()

	if i < 1 || i > permutations {
		return &token.PermutationError{
			Type: token.PermutationErrorIndexOutOfBound,
		}
	}

	l.permutation(i - 1)

	return nil
}

// Permutations returns the number of permutations for this token
func (l *Repeat) Permutations() uint {
	return uint(l.To() - l.From() + 1)
}

// PermutationsAll returns the number of all possible permutations for this token including its children
func (l *Repeat) PermutationsAll() uint {
	var sum uint
	from := l.From()

	if l.From() == 0 {
		sum++
		from++
	}

	tokenPermutations := l.token.PermutationsAll()

	for i := from; i <= l.To(); i++ {
		sum += uint(math.Pow(float64(tokenPermutations), float64(i)))
	}

	return sum
}

func (l *Repeat) String() string {
	var buffer bytes.Buffer

	for _, tok := range l.value {
		if _, err := buffer.WriteString(tok.String()); err != nil {
			panic(err)
		}
	}

	return buffer.String()
}

// List interface methods

// Get returns the current referenced token at the given index. The error return argument is not nil, if the index is out of bound.
func (l *Repeat) Get(i int) (token.Token, error) {
	if i < 0 || i >= len(l.value) {
		return nil, &ListError{ListErrorOutOfBound}
	}

	return l.value[i], nil
}

// Len returns the number of the current referenced tokens
func (l *Repeat) Len() int {
	return len(l.value)
}

// InternalGet returns the current referenced internal token at the given index. The error return argument is not nil, if the index is out of bound.
func (l *Repeat) InternalGet(i int) (token.Token, error) {
	if i != 0 {
		return nil, &ListError{ListErrorOutOfBound}
	}

	return l.token, nil
}

// InternalLen returns the number of referenced internal tokens
func (l *Repeat) InternalLen() int {
	return 1
}

// InternalLogicalRemove removes the referenced internal token and returns the replacement for the current token or nil if the current token should be removed.
func (l *Repeat) InternalLogicalRemove(tok token.Token) token.Token {
	if l.token == tok {
		return nil
	}

	return l
}

// InternalReplace replaces an old with a new internal token if it is referenced by this token. The error return argument is not nil, if the replacement is not suitable.
func (l *Repeat) InternalReplace(oldToken, newToken token.Token) error {
	if l.token == oldToken {
		l.token = newToken

		for i := range l.value {
			l.value[i] = l.token.Clone()
		}
	}

	return nil
}

// OptionalToken interface methods

// IsOptional checks dynamically if this token is in the current state optional
func (l *Repeat) IsOptional() bool { return l.From() == 0 }

// Activate activates this token
func (l *Repeat) Activate() {
	if l.From() != 0 {
		return
	}

	l.value = []token.Token{
		l.token.Clone(),
	}
}

// Deactivate deactivates this token
func (l *Repeat) Deactivate() {
	if l.From() != 0 {
		return
	}

	l.value = []token.Token{}
}

// ReduceToken interface methods

func combinations(n int, k int) (<-chan []int, chan<- struct{}) {
	ret := make(chan []int)
	cancel := make(chan struct{})

	go func() {
		is := make([]int, k)

		for i := 0; i < k; i++ {
			is[i] = i
		}

	CYCLE:
		for {
			// send the current progress
			cur := make([]int, k)
			copy(cur, is)
			select {
			case ret <- cur:
			case <-cancel:
				close(ret)

				return
			}

			// special case, of no elements to choose
			if k == 0 {
				// we reached the end

				break CYCLE
			}

			// increase the last element
			j := k - 1
			is[j]++

			if is[j] == n {
				// increase from the back to the front
				for j != 0 && is[j] == n {
					j--
					is[j]++
				}

				c := n
				found := false

				// do we have a increasing order up to the highest value at the end?
				for i := k - 1; i >= j; i-- {
					if is[i] != c {
						found = true

						break
					}

					c--
				}

				if !found {
					if j == 0 {
						// we reached the end

						break CYCLE
					} else {
						j--
						is[j]++
					}
				}

				jj := j

				for {
					j = jj

					// reset values
					for ; j < k-1; j++ {
						is[j+1] = is[j] + 1
					}

					// if after a reset the last value is still to high
					if is[k-1] == n {
						if jj > 0 {
							// start from an anterior index
							jj--

							is[jj]++

							continue
						} else {
							// we reached the end

							break CYCLE
						}
					}

					break
				}
			}
		}

		close(ret)
	}()

	return ret, cancel
}

// Reduce sets a specific reduction for this token
func (l *Repeat) Reduce(i uint) error {
	var count uint
	reduces := l.reduces()
	for _, le := range reduces {
		count += uint(le)
	}

	if count <= 1 || i < 1 || i > count {
		return &token.ReduceError{
			Type: token.ReduceErrorIndexOutOfBound,
		}
	}

	if !l.reducing {
		l.reducing = true
		l.reducingOriginalValue = l.value
	}

	j := 0

	if l.From() == 0 {
		if i == 1 {
			l.value = []token.Token{}

			return nil
		}

		i--
		j++
	}

	i--

	for i >= reduces[j] {
		i -= reduces[j]
		j++
	}

	var sel []int

	ch, cancel := combinations(len(l.reducingOriginalValue), j+int(l.From()))
	for c := range ch {
		if i == 0 {
			sel = c

			close(cancel)

			break
		}

		i--
	}

	tokens := make([]token.Token, len(sel))

	for i, c := range sel {
		tokens[i] = l.reducingOriginalValue[c]
	}

	l.value = tokens

	return nil
}

func factorial(n uint) uint {
	c := n
	n--

	for n > 0 {
		c *= n
		n--
	}

	return c
}

func (l *Repeat) reduces() []uint {
	n := uint(len(l.value))
	if l.reducing {
		n = uint(len(l.reducingOriginalValue))
	}

	le := n - uint(l.From()) + 1
	reduces := make([]uint, le)

	j := 0
	k := uint(l.From())

	if k == 0 {
		reduces[0] = 1

		j++
		k++
	}

	for ; k < n; k++ {
		reduces[j] = factorial(n) / (factorial(n-k) * factorial(k))
		j++
	}

	reduces[le-1] = 1

	return reduces
}

// Reduces returns the number of reductions for this token
func (l *Repeat) Reduces() uint {
	if l.reducing || int(l.From()) < len(l.value) {
		var count uint
		r := l.reduces()
		for _, le := range r {
			count += le
		}

		return count
	}

	return 0
}

// ResetToken interface methods

// Reset resets the (internal) state of this token and its dependences
func (l *Repeat) Reset() {
	// TODO reset the list if we depend on something else. this could and should be done in another way...
	_, okFrom := l.from.(*primitives.ConstantInt)
	_, okTo := l.to.(*primitives.ConstantInt)

	if !okFrom || !okTo {
		for _, tok := range l.value {
			token.ResetResetTokens(tok)
		}

		l.permutation(l.Permutations() - 1)
	}
}
