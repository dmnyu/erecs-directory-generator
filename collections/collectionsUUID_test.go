package collections

import (
	"testing"
)

func TestGetRepos(t *testing.T) {
	want := 3
	got := len(repoCodes)
	if want != got {
		t.Errorf("Wanted %d, got %d", want, got)
	}
}

func TestGetEntryMapForRepo(t *testing.T) {
	want := "a9d0a537-9afd-455f-bb26-ddfef5e3cdb8"
	got := GetEntryMapForRepo("tamwag")["oh066"]
	if want != got {
		t.Errorf("Wanted %s, got %s", want, got)
	}
}

func TestGetUUID(t *testing.T) {
	want := "a9d0a537-9afd-455f-bb26-ddfef5e3cdb8"
	got := GetUUID("tamwag", "oh066")
	if want != got {
		t.Errorf("Wanted %s, got %s", want, got)
	}
}
