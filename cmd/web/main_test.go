package main

import "testing"

func TestRun(t *testing.T) {
	_, err := run()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}
