package parser

import (
	"reflect"
	"testing"
)

func TestParseRESPV2(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedOutput []ParsedCmd
		expectError    bool
	}{
		{
			name:  "Valid RESP v2 response",
			input: "*3\r\n$3\r\nfoo\r\n$3\r\nbar\r\n$3\r\nbaz\r\n",
			expectedOutput: []ParsedCmd{
				{ValueType: RespEncodingConstants.BulkString, Value: "foo"},
				{ValueType: RespEncodingConstants.BulkString, Value: "bar"},
				{ValueType: RespEncodingConstants.BulkString, Value: "baz"},
			},
			expectError: false,
		},
		{
			name:           "Valid RESP v2 response with single element",
			input:          "*1\r\n$3\r\nfoo\r\n",
			expectedOutput: []ParsedCmd{{ValueType: RespEncodingConstants.BulkString, Value: "foo"}},
			expectError:    false,
		},
		{
			name:           "Empty input",
			input:          "",
			expectedOutput: []ParsedCmd{},
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
			expectedOutput: []ParsedCmd{{ValueType: RespEncodingConstants.BulkString, Value: "ping"}},
			expectError:    false,
		},
		{
			name:           "Invalid string format",
			input:          "*1\r\n$4\r\nping\r\\",
			expectedOutput: nil,
			expectError:    true,
		},
		{
			name:  "Valid RESP v2 response",
			input: "*5\r\n$3\r\nSET\r\n$6\r\nmangos\r\n$11\r\nwatermelons\r\n$2\r\nPX\r\n$3\r\n100\r\n",
			expectedOutput: []ParsedCmd{
				{ValueType: RespEncodingConstants.BulkString, Value: "SET"},
				{ValueType: RespEncodingConstants.BulkString, Value: "mangos"},
				{ValueType: RespEncodingConstants.BulkString, Value: "watermelons"},
				{ValueType: RespEncodingConstants.BulkString, Value: "PX"},
				{ValueType: RespEncodingConstants.BulkString, Value: "100"},
			},
			expectError: false,
		},
		{
			name:  "Valid RESP v2 response",
			input: "*5\r\n$3\r\nSET\r\n$6\r\nmangos\r\n$11\r\nwatermelons\r\n$2\r\nPX\r\n$3\r\n100\r\n",
			expectedOutput: []ParsedCmd{
				{ValueType: RespEncodingConstants.BulkString, Value: "SET"},
				{ValueType: RespEncodingConstants.BulkString, Value: "mangos"},
				{ValueType: RespEncodingConstants.BulkString, Value: "watermelons"},
				{ValueType: RespEncodingConstants.BulkString, Value: "PX"},
				{ValueType: RespEncodingConstants.BulkString, Value: "100"},
			},
			expectError: false,
		},
		{
			name:  "Valid RESP v2 response with different value types",
			input: "*5\r\n$3\r\nSET\r\n+mangos\r\n:+125\r\n-Error message\r\n$3\r\n100\r\n",
			expectedOutput: []ParsedCmd{
				{ValueType: RespEncodingConstants.BulkString, Value: "SET"},
				{ValueType: RespEncodingConstants.String, Value: "mangos"},
				{ValueType: RespEncodingConstants.Integer, Value: "+125"},
				{ValueType: RespEncodingConstants.Error, Value: "Error message"},
				{ValueType: RespEncodingConstants.BulkString, Value: "100"},
			},
			expectError: false,
		},
		{
			name:  "Valid RESP v2 response with string",
			input: "+Test\r\n",
			expectedOutput: []ParsedCmd{
				{
					ValueType: RespEncodingConstants.String, Value: "Test",
				},
			},
			expectError: false,
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
