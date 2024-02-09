package parser

import (
	"reflect"
	"testing"
)

func TestParseRESPV2(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedOutput []string
		expectError    bool
	}{
		{
			name:           "Valid RESP v2 response",
			input:          "*3\r\n$3\r\nfoo\r\n$3\r\nbar\r\n$3\r\nbaz\r\n",
			expectedOutput: []string{"foo", "bar", "baz"},
			expectError:    false,
		},
		{
			name:           "Valid RESP v2 response with single element",
			input:          "*1\r\n$3\r\nfoo\r\n",
			expectedOutput: []string{"foo"},
			expectError:    false,
		},
		{
			name:           "Empty input",
			input:          "",
			expectedOutput: []string{},
			expectError:    false,
		},
		{
			name:           "Invalid input",
			input:          "PING",
			expectedOutput: nil,
			expectError:    true,
		},
		{
			name:           "Valid RESP v2 response with 'ping'",
			input:          "*1\r\n$4\r\nping\r\n",
			expectedOutput: []string{"ping"},
			expectError:    false,
		},
	}

	p := &RespParser{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := p.HandleParse(test.input)

			// Check for error
			if test.expectError && err == nil {
				t.Error("Expected error, but got nil")
			} else if !test.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			// Check output
			if !reflect.DeepEqual(result, test.expectedOutput) {
				t.Errorf("Unexpected output. Expected: %v, Got: %v", test.expectedOutput, result)
			}
		})
	}
}
