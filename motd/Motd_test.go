package motd

import (
	"encoding/json"
	"testing"

	motdjava "github.com/xrcuo/api_boot/motdJava"
)

func TestBE(t *testing.T) {
	Host := ""

	Data, err := MotdBE(Host)
	if err != nil {
		t.Error(err)
		return
	}

	DataJson, err := json.Marshal(Data)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(string(DataJson))
}

func TestJava(t *testing.T) {
	Host := ""

	Data, err := motdjava.MotdJava(Host)
	if err != nil {
		t.Error(err)
		return
	}

	DataJson, err := json.Marshal(Data)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(string(DataJson))
}
