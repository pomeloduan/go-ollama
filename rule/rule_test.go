package rule

import "testing"

func TestRules(t *testing.T) {
	config, err := readConfig("./config.yml")
	if err != nil {
		t.Fatal("read config file error")
	}
	{ // case math
		_, ok := config.Rules["math"]
		if !ok {
			t.Fatal("expected get math rule")
		}
	}
	{ // case poet
		_, ok := config.Rules["poet"]
		if !ok {
			t.Fatal("expected get poet rule")
		}
	}
	{ // case hp
		_, ok := config.Rules["hp"]
		if !ok {
			t.Fatal("expected get hp rule")
		}
	}
}
