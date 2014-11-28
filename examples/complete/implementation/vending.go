package implementation

import (
	"errors"
)

const (
	coin25 = 25
	coin50 = 50
)

var (
	// ErrUnknownCoin states that a given coin is unknown to the vending machine
	ErrUnknownCoin = errors.New("Unknown coin")
)

// VendingMachine holds the state of a vending machine
type VendingMachine struct {
	credit int
}

// NewVendingMachine returns a instantiated state of a vending machine
func NewVendingMachine() *VendingMachine {
	return &VendingMachine{
		credit: 0,
	}
}

// Credit returns the current credit of the vending machine
func (v VendingMachine) Credit() int {
	return v.credit
}

// Coin inserts a coin into the vending machine. On success the credit of the machine will be increased by the coin.
func (v *VendingMachine) Coin(credit int) error {
	switch credit {
	case coin25:
		v.credit += credit
	case coin50:
		v.credit += credit
	default:
		return ErrUnknownCoin
	}

	return nil
}

// Vend executes a vend of the machine if enough credit (100) has been put in and returns true.
func (v *VendingMachine) Vend() bool {
	if v.credit < 100 {
		return false
	}

	v.credit -= 100

	return true
}
