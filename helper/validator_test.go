package helper_test

import (
	"testing"

	"github.com/lwinmgmg/outbox/helper"
)

func TestPanicEmptyString(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	helper.PanicEmptyString("", "panic")
}

func TestDefaultString(t *testing.T) {
	data := helper.DefaultString("", "ABC")
	if data != "ABC" {
		t.Errorf("Expected ABC, Getting %v", data)
	}
	data1 := helper.DefaultString("Original", "ABC")
	if data1 != "Original" {
		t.Errorf("Expected Original, Getting %v", data1)
	}
}

func TestDefaultInt(t *testing.T) {
	data := helper.DefaultInt(0, 10)
	if data != 10 {
		t.Errorf("Expected 10, Getting %v", data)
	}
	data1 := helper.DefaultInt(5, 10)
	if data1 != 5 {
		t.Errorf("Expected 5, Getting %v", data1)
	}
}

func TestDefaultBoolean(t *testing.T) {
	data := helper.DefaultBoolean(false, true)
	if !data {
		t.Errorf("Expected true, Getting %v", data)
	}
	data1 := helper.DefaultBoolean(true, false)
	if !data1 {
		t.Errorf("Expected true, Getting %v", data1)
	}
}
