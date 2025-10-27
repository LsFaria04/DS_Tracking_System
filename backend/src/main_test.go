package main

import (
    "testing"
)


// Just a dummy test to see if everything is configured correctly
func TestHelloName(t *testing.T) {
    name := "Gladys"

    if name == "Test" {
		t.Errorf("Test error %s", name);
	}

}