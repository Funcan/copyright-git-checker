package main

import (
	"regexp"
	"testing"
)

func TestPattern(t *testing.T) {
	tests := []struct {
		text        string
		should_fail bool
	}{
		{
			text:        "# (c) Copyright 2023 Hewlett Packard Enterprise Development LP",
			should_fail: true,
		},
		{
			text:        "# (c) Copyright 2023-2024 Hewlett Packard Enterprise Development LP",
			should_fail: true,
		},
		{
			text:        "# (c) Copyright 2023-24 Hewlett Packard Enterprise Development LP",
			should_fail: true,
		},
		{
			text:        "# (c) Copyright 2023, 2025 Hewlett Packard Enterprise Development LP",
			should_fail: false,
		},
		{
			text:        "# (c) Copyright 2023,2025 Hewlett Packard Enterprise Development LP",
			should_fail: false,
		},
		{
			text:        "# (c) Copyright 2023,25 Hewlett Packard Enterprise Development LP",
			should_fail: false,
		},
		{
			text:        "Some random text",
			should_fail: false,
		},
	}

	re := regexp.MustCompile(regex)

	for _, test := range tests {
		match := re.MatchString(test.text)

		if fails(test.text, "2025") != test.should_fail {
			t.Errorf("Expected match to be %v, got %v for text: %s", test.should_fail, match, test.text)
		}
	}
}
