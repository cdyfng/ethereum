// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

package bank_test

import (
	"sync"
	"testing"
	"time"

	"github.com/cdyfng/ethereum/goplsource/ch9/bank3"
)

func TestBank(t *testing.T) {
	// Deposit [1..1000] concurrently.
	var n sync.WaitGroup
	for i := 1; i <= 1000; i++ {
		n.Add(1)
		go func(amount int) {
			bank.Deposit(amount)
			time.Sleep(time.Microsecond * 2)
			n.Done()
		}(i)
	}
	n.Wait()
	got, want := bank.Balance(), (1000+1)*1000/2
	if got != want {
		t.Errorf("Balance = %d, want %d", got, want)
	}
	t.Log("Balance = %d, want %d", got, want)
}
