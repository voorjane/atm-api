package pkg

import (
	"errors"
	"sync"
)

type BankAccount interface {
	Deposit(amount float64) error
	Withdraw(amount float64) error
	GetBalance() float64
}

type Account struct {
	Id      int
	Balance float64
	mux     sync.Mutex `gorm:"-:all"`
}

func (a *Account) Deposit(amount float64) error {
	a.mux.Lock()
	defer a.mux.Unlock()
	if amount <= 0 {
		return errors.New("deposit amount must be positive")
	}
	a.Balance += amount
	return nil
}

func (a *Account) Withdraw(amount float64) error {
	a.mux.Lock()
	defer a.mux.Unlock()
	if amount <= 0 {
		return errors.New("withdraw amount must be positive")
	}
	if amount > a.Balance {
		return errors.New("not enough money for this operation")
	}
	a.Balance -= amount
	return nil
}

func (a *Account) GetBalance() float64 {
	a.mux.Lock()
	defer a.mux.Unlock()
	return a.Balance
}
