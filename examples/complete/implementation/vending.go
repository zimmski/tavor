package implementation

import (
	"errors"
)

const (
	Coin25 = 25
	Coin50 = 50
)

var (
	ErrUnknownCoin = errors.New("Unknown coin")
)

type VendingMachine struct {
	credit int
}

func NewVendingMachine() *VendingMachine {
	return &VendingMachine{
		credit: 0,
	}
}

func (v VendingMachine) Credit() int {
	return v.credit
}

func (v *VendingMachine) Coin(credit int) error {
	switch credit {
	case Coin25:
		v.credit += credit
	case Coin50:
		v.credit += credit
	default:
		return ErrUnknownCoin
	}

	return nil
}

func (v *VendingMachine) Vend() bool {
	if v.credit < 100 {
		return false
	}

	v.credit -= 100

	return true
}
