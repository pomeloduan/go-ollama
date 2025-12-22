package rule

import "testing"

func TestGetMathRule(t *testing.T) {
	rule := UseRule("math")
	if rule == nil {
		t.Fatal("expected get math rule")
	}
	if rule.SystemMessage() == "" {
		t.Fatal("expected get math rule system message")
	}
}

func TestGoodAnswer(t *testing.T) {
	rule := UseRule("math")
	if rule == nil {
		t.Fatal("expected get math rule")
	}
	const text = `isMath: true
resolvation: good answer`
	answer := rule.ParseAnswer(text)
	if answer != "good answer" {
		t.Fatal("expected get answer")
	}
}

func TestBadAnswerNotMath(t *testing.T) {
	rule := UseRule("math")
	if rule == nil {
		t.Fatal("expected get math rule")
	}
	// case 0
	{
		const text = `isMath: false
resolvation: text`
		answer := rule.ParseAnswer(text)
		if answer != "不是数学问题" {
			t.Fatal("expected get 不是数学问题")
		}
	}
	// case 1
	{
		const text = `isMath: any
resolvation: text`
		answer := rule.ParseAnswer(text)
		if answer != "不是数学问题" {
			t.Fatal("expected get 不是数学问题")
		}
	}
	// case 2
	{
		const text = `resolvation: text`
		answer := rule.ParseAnswer(text)
		if answer != "不是数学问题" {
			t.Fatal("expected get 不是数学问题")
		}
	}
}

func TestBadAnswerNoResolvation(t *testing.T) {
	rule := UseRule("math")
	if rule == nil {
		t.Fatal("expected get math rule")
	}

	const text = `isMath: true`
	answer := rule.ParseAnswer(text)
	if answer != "不是数学问题" {
		t.Fatal("expected get 不是数学问题")
	}
}

func TestBadAnswerNotFormated(t *testing.T) {
	rule := UseRule("math")
	if rule == nil {
		t.Fatal("expected get math rule")
	}
	// case 0
	{
		const text = `a: b
c: d`
		answer := rule.ParseAnswer(text)
		if answer != "不是数学问题" {
			t.Fatal("expected get 不是数学问题")
		}
	}
	// case 1
	{
		const text = `any`
		answer := rule.ParseAnswer(text)
		if answer != "不是数学问题" {
			t.Fatal("expected get 不是数学问题")
		}
	}
}
