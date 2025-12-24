package rule

import "testing"

func TestRules(t *testing.T) {
	manager, err := StartRuleManager()
	if err != nil {
		t.Fatal("start rule manager error")
	}
	{ // case math
		rule, ok := manager.ruleMap["math"]
		if !ok || rule.name != "math" {
			t.Fatal("expected get math rule")
		}
	}
	{ // case poet
		rule, ok := manager.ruleMap["poet"]
		if !ok || rule.name != "poet" {
			t.Fatal("expected get poet rule")
		}
	}
	{ // case hp
		rule, ok := manager.ruleMap["hp"]
		if !ok || rule.name != "hp" {
			t.Fatal("expected get hp rule")
		}
	}
}
