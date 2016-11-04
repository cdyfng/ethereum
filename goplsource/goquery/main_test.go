package main

import "testing"

func Test(t *testing.T) {
	p := getEtherPrice()
	if p <= 0 {
		t.Errorf("p(%f ) < 0 ",	p)
	}
}
