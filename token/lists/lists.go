package lists

import (
	"github.com/zimmski/tavor/token"
)

type List interface {
	token.Token

	Len() int
}
