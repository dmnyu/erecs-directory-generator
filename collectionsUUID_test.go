package main

import "testing"

func TestContainsCol(t *testing.T) {

	col := "oh66"
	want := "bc5ded73-8a38-43f3-80a3-29ed3775d869"
	got := getColUUID(col)

	if got != want {
		t.Log("Dumping Cols")
		printCols()
		t.Errorf("Wanted %s, Got %s", want, got)
	}
}
