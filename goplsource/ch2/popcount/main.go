// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 45.

// (Package doc comment intentionally malformed to demonstrate golint.)
//!+
package popcount

import "fmt"

// pc[i] is the population count of i.
var pc [256]byte

func init() {
	for i := range pc {
		pc[i] = pc[i/2] + byte(i&1)
		//fmt.printf("pc[i]: "+pc[i] + "  pc[i/2]: " + pc[i/2] + "  byte(i&1):" + byte(i&1) )
		fmt.Printf("pc[%d]: %d pc[%d/2]: %d byte(i&1): %d\n", i, pc[i], i, pc[i/2], byte(i&1))
		//fmt.Printf("%g°F = %g°C\n", freezingF, fToC(freezingF))
	}
}

// PopCount returns the population count (number of set bits) of x.
func PopCount(x uint64) int {
	return int(pc[byte(x>>(0*8))] +
		pc[byte(x>>(1*8))] +
		pc[byte(x>>(2*8))] +
		pc[byte(x>>(3*8))] +
		pc[byte(x>>(4*8))] +
		pc[byte(x>>(5*8))] +
		pc[byte(x>>(6*8))] +
		pc[byte(x>>(7*8))])
}

//!-
