package main

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	usage()
	result := m.Run()
	fmt.Println("Tests complete....")
	os.Exit(result)
}

func TestGetPrice(t *testing.T) {
	res, err := getPrice("£1.80/unit")
	if err != nil {
		t.Error("Got error message ", err)
	}
	if res != 1.8 {

	}
}
func TestRound(t *testing.T) {
	t.Parallel()
	testRound := round(1.3654)
	if testRound != 1 {
		t.Error("Invalid round, round should round down to 1, but got ", testRound)
	}

	testRound = round(1.999)
	if testRound != 2 {
		t.Error("Invalid round, round should round up to 2, but got ", testRound)
	}

}

func TestToFixed(t *testing.T) {
	t.Parallel()
	testRound := toFixed(1.3654, 2)
	if testRound != 1.37 {
		t.Error("Invalid toFixed, round should round up to 1.37, but got ", testRound)
	}
	testRound = toFixed(1.3645, 2)
	if testRound != 1.36 {
		t.Error("Invalid toFixed, round should round down to 1.36, but got ", testRound)
	}
	testRound = toFixed(1.544444444444444445, 0)
	if testRound != 2 {
		t.Error("Invalid toFixed, round should round up to 2, but got ", testRound)
	}
	testRound = toFixed(1.444444444444444445, 0)
	if testRound != 1 {
		t.Error("Invalid toFixed, round should round down to 1, but got ", testRound)
	}

}
