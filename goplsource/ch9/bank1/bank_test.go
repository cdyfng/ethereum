// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

package bank_test

import (
	"fmt"
	"testing"

	"github.com/cdyfng/ethereum/goplsource/ch9/bank1"
)

func TestBank(t *testing.T) {
	//done := make(chan struct{})
	done := make(chan int)

	// Alice
	go func() {
		bank.Deposit(200)
		fmt.Println("=", bank.Balance())
		done <- 1//struct{}{}
	}()

	// Bob
	go func() {
		bank.Deposit(100)
		done <- 1//struct{}{}
	}()

	// Wait for both transactions.
	<-done
	<-done

	if got, want := bank.Balance(), 300; got != want {
		t.Errorf("Balance = %d, want %d", got, want)
	}
}
