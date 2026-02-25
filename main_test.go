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

func TestFixLine(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		year     string
		expected string
	}{
		{
			name:     "Single year to range",
			input:    "# (c) Copyright 2025 Hewlett Packard Enterprise Development LP",
			year:     "2026",
			expected: "# (c) Copyright 2025-2026 Hewlett Packard Enterprise Development LP",
		},
		{
			name:     "Update existing range end year",
			input:    "# (c) Copyright 2024-2025 Hewlett Packard Enterprise Development LP",
			year:     "2026",
			expected: "# (c) Copyright 2024-2026 Hewlett Packard Enterprise Development LP",
		},
		{
			name:     "Update abbreviated range",
			input:    "# (c) Copyright 2023-24 Hewlett Packard Enterprise Development LP",
			year:     "2026",
			expected: "# (c) Copyright 2023-2026 Hewlett Packard Enterprise Development LP",
		},
		{
			name:     "Old single year to range",
			input:    "# (c) Copyright 2023 Hewlett Packard Enterprise Development LP",
			year:     "2026",
			expected: "# (c) Copyright 2023-2026 Hewlett Packard Enterprise Development LP",
		},
		{
			name:     "Comma separated years - extend last to range",
			input:    "# (c) Copyright 2020, 2023 Hewlett Packard Enterprise Development LP",
			year:     "2026",
			expected: "# (c) Copyright 2020, 2023-2026 Hewlett Packard Enterprise Development LP",
		},
		{
			name:     "Comma separated years no space - extend last to range",
			input:    "# (C) Copyright 2023,2025 Hewlett Packard Enterprise Development LP",
			year:     "2026",
			expected: "# (C) Copyright 2023,2025-2026 Hewlett Packard Enterprise Development LP",
		},
		{
			name:     "No change needed - non-copyright line",
			input:    "Some random text",
			year:     "2026",
			expected: "Some random text",
		},
		{
			name:     "No change if already current year",
			input:    "# (c) Copyright 2026 Hewlett Packard Enterprise Development LP",
			year:     "2026",
			expected: "# (c) Copyright 2026 Hewlett Packard Enterprise Development LP",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := fixLine(test.input, test.year)
			if result != test.expected {
				t.Errorf("fixLine(%q, %q) = %q, expected %q", test.input, test.year, result, test.expected)
			}
		})
	}
}
