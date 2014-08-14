package lists

import (
	"bytes"
	"math"
	"strconv"

	"github.com/zimmski/tavor"
	"github.com/zimmski/tavor/rand"
	"github.com/zimmski/tavor/token"
	"github.com/zimmski/tavor/token/primitives"
)

type Repeat struct {
	from  token.Token
	to    token.Token
	token token.Token
	value []token.Token

	reducing              bool
	reducingOriginalValue []token.Token
}

func NewRepeat(tok token.Token, from int64, to int64) *Repeat {
	return NewRepeatWithTokens(tok, primitives.NewConstantInt(int(from)), primitives.NewConstantInt(int(to)))
}

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

func (l *Repeat) From() int64 {
	iFrom, err := strconv.Atoi(l.from.String())
	if err != nil {
		panic(err) // TODO
	}

	return int64(iFrom)
}

func (l *Repeat) To() int64 {
	iTo, err := strconv.Atoi(l.to.String())
	if err != nil {
		panic(err) // TODO
	}

	return int64(iTo)
}

// Token interface methods

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

func (l *Repeat) Fuzz(r rand.Rand) {
	i := r.Intn(int(l.To() - l.From() + 1))

	l.permutation(i)
}

func (l *Repeat) FuzzAll(r rand.Rand) {
	l.Fuzz(r)

	for _, tok := range l.value {
		tok.FuzzAll(r)
	}
}

func (l *Repeat) Parse(pars *token.InternalParser, cur int) (int, []error) {
	var toks []token.Token

	for i := 1; i <= int(l.From()); i++ {
		tok := l.token.Clone()

		nex, errs := tok.Parse(pars, cur)

		if len(errs) != 0 {
			return cur, errs
		}

		cur = nex
		toks = append(toks, tok)
	}

	for i := l.From(); i < l.To(); i++ {
		tok := l.token.Clone()

		nex, errs := tok.Parse(pars, cur)

		if len(errs) != 0 {
			break
		}

		cur = nex
		toks = append(toks, tok)
	}

	l.value = toks

	return cur, nil
}

func (l *Repeat) permutation(i int) {
	toks := make([]token.Token, i+int(l.From()))

	for i := range toks {
		toks[i] = l.token.Clone()
	}

	l.value = toks
}

func (l *Repeat) Permutation(i int) error {
	permutations := l.Permutations()

	if i < 1 || i > permutations {
		return &token.PermutationError{
			Type: token.PermutationErrorIndexOutOfBound,
		}
	}

	l.permutation(i - 1)

	return nil
}

func (l *Repeat) Permutations() int {
	return int(l.To() - l.From() + 1)
}

func (l *Repeat) PermutationsAll() int {
	sum := 0
	from := l.From()

	if l.From() == 0 {
		sum++
		from++
	}

	tokenPermutations := l.token.PermutationsAll()

	for i := from; i <= l.To(); i++ {
		sum += int(math.Pow(float64(tokenPermutations), float64(i)))
	}

	return sum
}

func (l *Repeat) String() string {
	var buffer bytes.Buffer

	for _, tok := range l.value {
		buffer.WriteString(tok.String())
	}

	return buffer.String()
}

// List interface methods

func (l *Repeat) Get(i int) (token.Token, error) {
	if i < 0 || i >= len(l.value) {
		return nil, &ListError{ListErrorOutOfBound}
	}

	return l.value[i], nil
}

func (l *Repeat) Len() int {
	return len(l.value)
}

func (l *Repeat) InternalGet(i int) (token.Token, error) {
	if i != 0 {
		return nil, &ListError{ListErrorOutOfBound}
	}

	return l.token, nil
}

func (l *Repeat) InternalLen() int {
	return 1
}

func (l *Repeat) InternalLogicalRemove(tok token.Token) token.Token {
	if l.token == tok {
		return nil
	}

	return l
}

func (l *Repeat) InternalReplace(oldToken, newToken token.Token) {
	if l.token == oldToken {
		l.token = newToken

		for i := range l.value {
			l.value[i] = l.token.Clone()
		}
	}
}

// OptionalToken interface methods

func (l *Repeat) IsOptional() bool { return l.From() == 0 }
func (l *Repeat) Activate() {
	if l.From() != 0 {
		return
	}

	l.value = []token.Token{
		l.token.Clone(),
	}
}
func (l *Repeat) Deactivate() {
	if l.From() != 0 {
		return
	}

	l.value = []token.Token{}
}

// ReduceToken interface methods

func combinations(n int, k int) <-chan []int {
	ret := make(chan []int)

	go func() {
		is := make([]int, k)

		for i := 0; i < k; i++ {
			is[i] = i
		}

		for {
			// send the current progress
			cur := make([]int, k)
			copy(cur, is)
			ret <- cur

			// special case, of no elements to choose
			if k == 0 {
				// we reached the end

				break
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

						break
					} else {
						j--
						is[j]++
					}
				}

				// reset values
				for ; j < k-1; j++ {
					is[j+1] = is[j] + 1
				}

				// if after a reset the last value is still to high we are done
				if is[k-1] == n {
					// we reached the end

					break
				}
			}
		}

		close(ret)
	}()

	return ret
}

func (l *Repeat) Reduce(i int) error {
	count := 0
	reduces := l.reduces()
	for _, le := range reduces {
		count += le
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

	for c := range combinations(len(l.reducingOriginalValue), j+int(l.From())) {
		if i == 0 {
			sel = c

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

func factorial(n int) int {
	c := n
	n--

	for n != 0 {
		c *= n
		n--
	}

	return c
}

func (l *Repeat) reduces() []int {
	n := len(l.value)
	if l.reducing {
		n = len(l.reducingOriginalValue)
	}

	le := int(n - int(l.From()) + 1)
	reduces := make([]int, le)

	j := 0
	k := int(l.From())

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

func (l *Repeat) Reduces() int {
	if l.reducing || int(l.From()) < len(l.value) {
		count := 0
		r := l.reduces()
		for _, le := range r {
			count += le
		}

		return count
	}

	return 0
}

// ResetToken interface methods

func (l *Repeat) Reset() {
	// TODO reset the list if we depend on something else. this could and should be done in another way...
	_, okFrom := l.from.(*primitives.ConstantInt)
	_, okTo := l.to.(*primitives.ConstantInt)

	if !okFrom || !okTo {
		l.permutation(l.Permutations() - 1)
	}
}
