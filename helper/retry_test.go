package helper_test

import (
	"errors"
	"testing"

	"github.com/lwinmgmg/outbox/helper"
)

func TestRetry(t *testing.T) {
	count := 5
	actualCount := []int{0}
	helper.Retry(count, func(t ...int) error {
		t[0]++
		return errors.New("Testing")
	}, actualCount...)
	if actualCount[0] != 5 {
		t.Errorf("Expecting : %v and Getting : %v", count, actualCount[0])
	}
}
