package goasana

import (
	"testing"
)

func TestGetMe(t *testing.T) {
	me, err := GetMe()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if me.Name != "Amodio Pesce" {
		t.Log("Name not match")
		t.FailNow()
	}
}
