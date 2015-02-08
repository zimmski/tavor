package expressions

import (
	"strconv"

	"github.com/zimmski/tavor/token"
	"github.com/zimmski/tavor/token/lists"
)

// AddArithmetic implements an arithmetic token adding the values of two tokens
type AddArithmetic struct {
	a token.Token
	b token.Token
}

// NewAddArithmetic returns a new instance of a AddArithmetic token
func NewAddArithmetic(a, b token.Token) *AddArithmetic {
	return &AddArithmetic{
		a: a,
		b: b,
	}
}

// Clone returns a copy of the token and all its children
func (e *AddArithmetic) Clone() token.Token {
	return &AddArithmetic{
		a: e.a.Clone(),
		b: e.b.Clone(),
	}
}

// Parse tries to parse the token beginning from the current position in the parser data.
// If the parsing is successful the error argument is nil and the next current position after the token is returned.
func (e *AddArithmetic) Parse(pars *token.InternalParser, cur int) (int, []error) {
	panic("TODO implement")
}

// Permutation sets a specific permutation for this token
func (e *AddArithmetic) Permutation(i uint) error {
	permutations := e.Permutations()

	if i < 1 || i > permutations {
		return &token.PermutationError{
			Type: token.PermutationErrorIndexOutOfBound,
		}
	}
	// do nothing

	return nil
}

// Permutations returns the number of permutations for this token
func (e *AddArithmetic) Permutations() uint {
	return 1
}

// PermutationsAll returns the number of all possible permutations for this token including its children
func (e *AddArithmetic) PermutationsAll() uint {
	return e.a.PermutationsAll() * e.b.PermutationsAll()
}

func (e *AddArithmetic) String() string {
	as := e.a.String()
	bs := e.b.String()

	if as == "" || bs == "" || as == "TODO" || bs == "TODO" {
		return "TODO"
	}

	a, err := strconv.Atoi(as)
	if err != nil {
		panic(err)
	}
	b, err := strconv.Atoi(bs)
	if err != nil {
		panic(err)
	}

	return strconv.Itoa(a + b)
}

// List interface methods

// Get returns the current referenced token at the given index. The error return argument is not nil, if the index is out of bound.
func (e *AddArithmetic) Get(i int) (token.Token, error) {
	switch i {
	case 0:
		return e.a, nil
	case 1:
		return e.b, nil
	default:
		return nil, &lists.ListError{
			Type: lists.ListErrorOutOfBound,
		}
	}
}

// Len returns the number of the current referenced tokens
func (e *AddArithmetic) Len() int {
	return 2
}

// InternalGet returns the current referenced internal token at the given index. The error return argument is not nil, if the index is out of bound.
func (e *AddArithmetic) InternalGet(i int) (token.Token, error) {
	return e.Get(i)
}

// InternalLen returns the number of referenced internal tokens
func (e *AddArithmetic) InternalLen() int {
	return e.Len()
}

// InternalLogicalRemove removes the referenced internal token and returns the replacement for the current token or nil if the current token should be removed.
func (e *AddArithmetic) InternalLogicalRemove(tok token.Token) token.Token {
	if tok == e.a || tok == e.b {
		return nil
	}

	return e
}

// InternalReplace replaces an old with a new internal token if it is referenced by this token. The error return argument is not nil, if the replacement is not suitable.
func (e *AddArithmetic) InternalReplace(oldToken, newToken token.Token) error {
	if oldToken == e.a {
		e.a = newToken
	}
	if oldToken == e.b {
		e.b = newToken
	}

	return nil
}

// SubArithmetic implements an arithmetic token subtracting the values of two tokens
type SubArithmetic struct {
	a token.Token
	b token.Token
}

// NewSubArithmetic returns a new instance of a SubArithmetic token
func NewSubArithmetic(a, b token.Token) *SubArithmetic {
	return &SubArithmetic{
		a: a,
		b: b,
	}
}

// Clone returns a copy of the token and all its children
func (e *SubArithmetic) Clone() token.Token {
	return &SubArithmetic{
		a: e.a.Clone(),
		b: e.b.Clone(),
	}
}

// Parse tries to parse the token beginning from the current position in the parser data.
// If the parsing is successful the error argument is nil and the next current position after the token is returned.
func (e *SubArithmetic) Parse(pars *token.InternalParser, cur int) (int, []error) {
	panic("TODO implement")
}

// Permutation sets a specific permutation for this token
func (e *SubArithmetic) Permutation(i uint) error {
	permutations := e.Permutations()

	if i < 1 || i > permutations {
		return &token.PermutationError{
			Type: token.PermutationErrorIndexOutOfBound,
		}
	}
	// do nothing

	return nil
}

// Permutations returns the number of permutations for this token
func (e *SubArithmetic) Permutations() uint {
	return 1
}

// PermutationsAll returns the number of all possible permutations for this token including its children
func (e *SubArithmetic) PermutationsAll() uint {
	return e.a.PermutationsAll() * e.b.PermutationsAll()
}

func (e *SubArithmetic) String() string {
	as := e.a.String()
	bs := e.b.String()

	if as == "" || bs == "" || as == "TODO" || bs == "TODO" {
		return "TODO"
	}

	a, err := strconv.Atoi(as)
	if err != nil {
		panic(err)
	}
	b, err := strconv.Atoi(bs)
	if err != nil {
		panic(err)
	}

	return strconv.Itoa(a - b)
}

// List interface methods

// Get returns the current referenced token at the given index. The error return argument is not nil, if the index is out of bound.
func (e *SubArithmetic) Get(i int) (token.Token, error) {
	switch i {
	case 0:
		return e.a, nil
	case 1:
		return e.b, nil
	default:
		return nil, &lists.ListError{
			Type: lists.ListErrorOutOfBound,
		}
	}
}

// Len returns the number of the current referenced tokens
func (e *SubArithmetic) Len() int {
	return 2
}

// InternalGet returns the current referenced internal token at the given index. The error return argument is not nil, if the index is out of bound.
func (e *SubArithmetic) InternalGet(i int) (token.Token, error) {
	return e.Get(i)
}

// InternalLen returns the number of referenced internal tokens
func (e *SubArithmetic) InternalLen() int {
	return e.Len()
}

// InternalLogicalRemove removes the referenced internal token and returns the replacement for the current token or nil if the current token should be removed.
func (e *SubArithmetic) InternalLogicalRemove(tok token.Token) token.Token {
	if tok == e.a || tok == e.b {
		return nil
	}

	return e
}

// InternalReplace replaces an old with a new internal token if it is referenced by this token. The error return argument is not nil, if the replacement is not suitable.
func (e *SubArithmetic) InternalReplace(oldToken, newToken token.Token) error {
	if oldToken == e.a {
		e.a = newToken
	}
	if oldToken == e.b {
		e.b = newToken
	}

	return nil
}

// MulArithmetic implements an arithmetic token multiplying the values of two tokens
type MulArithmetic struct {
	a token.Token
	b token.Token
}

// NewMulArithmetic returns a new instance of a MulArithmetic token
func NewMulArithmetic(a, b token.Token) *MulArithmetic {
	return &MulArithmetic{
		a: a,
		b: b,
	}
}

// Clone returns a copy of the token and all its children
func (e *MulArithmetic) Clone() token.Token {
	return &MulArithmetic{
		a: e.a.Clone(),
		b: e.b.Clone(),
	}
}

// Parse tries to parse the token beginning from the current position in the parser data.
// If the parsing is successful the error argument is nil and the next current position after the token is returned.
func (e *MulArithmetic) Parse(pars *token.InternalParser, cur int) (int, []error) {
	panic("TODO implement")
}

// Permutation sets a specific permutation for this token
func (e *MulArithmetic) Permutation(i uint) error {
	permutations := e.Permutations()

	if i < 1 || i > permutations {
		return &token.PermutationError{
			Type: token.PermutationErrorIndexOutOfBound,
		}
	}
	// do nothing

	return nil
}

// Permutations returns the number of permutations for this token
func (e *MulArithmetic) Permutations() uint {
	return 1
}

// PermutationsAll returns the number of all possible permutations for this token including its children
func (e *MulArithmetic) PermutationsAll() uint {
	return e.a.PermutationsAll() * e.b.PermutationsAll()
}

func (e *MulArithmetic) String() string {
	as := e.a.String()
	bs := e.b.String()

	if as == "" || bs == "" || as == "TODO" || bs == "TODO" {
		return "TODO"
	}

	a, err := strconv.Atoi(as)
	if err != nil {
		panic(err)
	}
	b, err := strconv.Atoi(bs)
	if err != nil {
		panic(err)
	}

	return strconv.Itoa(a * b)
}

// List interface methods

// Get returns the current referenced token at the given index. The error return argument is not nil, if the index is out of bound.
func (e *MulArithmetic) Get(i int) (token.Token, error) {
	switch i {
	case 0:
		return e.a, nil
	case 1:
		return e.b, nil
	default:
		return nil, &lists.ListError{
			Type: lists.ListErrorOutOfBound,
		}
	}
}

// Len returns the number of the current referenced tokens
func (e *MulArithmetic) Len() int {
	return 2
}

// InternalGet returns the current referenced internal token at the given index. The error return argument is not nil, if the index is out of bound.
func (e *MulArithmetic) InternalGet(i int) (token.Token, error) {
	return e.Get(i)
}

// InternalLen returns the number of referenced internal tokens
func (e *MulArithmetic) InternalLen() int {
	return e.Len()
}

// InternalLogicalRemove removes the referenced internal token and returns the replacement for the current token or nil if the current token should be removed.
func (e *MulArithmetic) InternalLogicalRemove(tok token.Token) token.Token {
	if tok == e.a || tok == e.b {
		return nil
	}

	return e
}

// InternalReplace replaces an old with a new internal token if it is referenced by this token. The error return argument is not nil, if the replacement is not suitable.
func (e *MulArithmetic) InternalReplace(oldToken, newToken token.Token) error {
	if oldToken == e.a {
		e.a = newToken
	}
	if oldToken == e.b {
		e.b = newToken
	}

	return nil
}

// DivArithmetic implements an arithmetic token dividing the values of two tokens
type DivArithmetic struct {
	a token.Token
	b token.Token
}

// NewDivArithmetic returns a new instance of a DivArithmetic token
func NewDivArithmetic(a, b token.Token) *DivArithmetic {
	return &DivArithmetic{
		a: a,
		b: b,
	}
}

// Clone returns a copy of the token and all its children
func (e *DivArithmetic) Clone() token.Token {
	return &DivArithmetic{
		a: e.a.Clone(),
		b: e.b.Clone(),
	}
}

// Parse tries to parse the token beginning from the current position in the parser data.
// If the parsing is successful the error argument is nil and the next current position after the token is returned.
func (e *DivArithmetic) Parse(pars *token.InternalParser, cur int) (int, []error) {
	panic("TODO implement")
}

// Permutation sets a specific permutation for this token
func (e *DivArithmetic) Permutation(i uint) error {
	permutations := e.Permutations()

	if i < 1 || i > permutations {
		return &token.PermutationError{
			Type: token.PermutationErrorIndexOutOfBound,
		}
	}
	// do nothing

	return nil
}

// Permutations returns the number of permutations for this token
func (e *DivArithmetic) Permutations() uint {
	return 1
}

// PermutationsAll returns the number of all possible permutations for this token including its children
func (e *DivArithmetic) PermutationsAll() uint {
	return e.a.PermutationsAll() * e.b.PermutationsAll()
}

func (e *DivArithmetic) String() string {
	as := e.a.String()
	bs := e.b.String()

	if as == "" || bs == "" || as == "TODO" || bs == "TODO" {
		return "TODO"
	}

	a, err := strconv.Atoi(as)
	if err != nil {
		panic(err)
	}
	b, err := strconv.Atoi(bs)
	if err != nil {
		panic(err)
	}

	return strconv.Itoa(a / b)
}

// List interface methods

// Get returns the current referenced token at the given index. The error return argument is not nil, if the index is out of bound.
func (e *DivArithmetic) Get(i int) (token.Token, error) {
	switch i {
	case 0:
		return e.a, nil
	case 1:
		return e.b, nil
	default:
		return nil, &lists.ListError{
			Type: lists.ListErrorOutOfBound,
		}
	}
}

// Len returns the number of the current referenced tokens
func (e *DivArithmetic) Len() int {
	return 2
}

// InternalGet returns the current referenced internal token at the given index. The error return argument is not nil, if the index is out of bound.
func (e *DivArithmetic) InternalGet(i int) (token.Token, error) {
	return e.Get(i)
}

// InternalLen returns the number of referenced internal tokens
func (e *DivArithmetic) InternalLen() int {
	return e.Len()
}

// InternalLogicalRemove removes the referenced internal token and returns the replacement for the current token or nil if the current token should be removed.
func (e *DivArithmetic) InternalLogicalRemove(tok token.Token) token.Token {
	if tok == e.a || tok == e.b {
		return nil
	}

	return e
}

// InternalReplace replaces an old with a new internal token if it is referenced by this token. The error return argument is not nil, if the replacement is not suitable.
func (e *DivArithmetic) InternalReplace(oldToken, newToken token.Token) error {
	if oldToken == e.a {
		e.a = newToken
	}
	if oldToken == e.b {
		e.b = newToken
	}

	return nil
}
