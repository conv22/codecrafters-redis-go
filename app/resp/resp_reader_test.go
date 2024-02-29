package resp

import (
	"reflect"
	"strings"
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
				{ValueType: RESP_ENCODING_CONSTANTS.BULK_STRING, Value: "foo"},
				{ValueType: RESP_ENCODING_CONSTANTS.BULK_STRING, Value: "bar"},
				{ValueType: RESP_ENCODING_CONSTANTS.BULK_STRING, Value: "baz"},
			},
			expectError: false,
		},
		{
			name:           "Valid RESP v2 response with single element",
			input:          "*1\r\n$3\r\nfoo\r\n",
			expectedOutput: []ParsedCmd{{ValueType: RESP_ENCODING_CONSTANTS.BULK_STRING, Value: "foo"}},
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
			expectedOutput: []ParsedCmd{{ValueType: RESP_ENCODING_CONSTANTS.BULK_STRING, Value: "ping"}},
			expectError:    false,
		},
		{
			name:  "Valid RESP v2 response",
			input: "*5\r\n$3\r\nSET\r\n$6\r\nmangos\r\n$11\r\nwatermelons\r\n$2\r\nPX\r\n$3\r\n100\r\n",
			expectedOutput: []ParsedCmd{
				{ValueType: RESP_ENCODING_CONSTANTS.BULK_STRING, Value: "SET"},
				{ValueType: RESP_ENCODING_CONSTANTS.BULK_STRING, Value: "mangos"},
				{ValueType: RESP_ENCODING_CONSTANTS.BULK_STRING, Value: "watermelons"},
				{ValueType: RESP_ENCODING_CONSTANTS.BULK_STRING, Value: "PX"},
				{ValueType: RESP_ENCODING_CONSTANTS.BULK_STRING, Value: "100"},
			},
			expectError: false,
		},
		{
			name:  "Valid RESP v2 response",
			input: "*5\r\n$3\r\nSET\r\n$6\r\nmangos\r\n$11\r\nwatermelons\r\n$2\r\nPX\r\n$3\r\n100\r\n",
			expectedOutput: []ParsedCmd{
				{ValueType: RESP_ENCODING_CONSTANTS.BULK_STRING, Value: "SET"},
				{ValueType: RESP_ENCODING_CONSTANTS.BULK_STRING, Value: "mangos"},
				{ValueType: RESP_ENCODING_CONSTANTS.BULK_STRING, Value: "watermelons"},
				{ValueType: RESP_ENCODING_CONSTANTS.BULK_STRING, Value: "PX"},
				{ValueType: RESP_ENCODING_CONSTANTS.BULK_STRING, Value: "100"},
			},
			expectError: false,
		},
		{
			name:  "Valid RESP v2 response with different value types",
			input: "*5\r\n$3\r\nSET\r\n+mangos\r\n:+125\r\n-Error message\r\n$3\r\n100\r\n",
			expectedOutput: []ParsedCmd{
				{ValueType: RESP_ENCODING_CONSTANTS.BULK_STRING, Value: "SET"},
				{ValueType: RESP_ENCODING_CONSTANTS.STRING, Value: "mangos"},
				{ValueType: RESP_ENCODING_CONSTANTS.INTEGER, Value: "+125"},
				{ValueType: RESP_ENCODING_CONSTANTS.ERROR, Value: "Error message"},
				{ValueType: RESP_ENCODING_CONSTANTS.BULK_STRING, Value: "100"},
			},
			expectError: false,
		},
		{
			name:  "Valid RESP v2 response with different value types",
			input: "*5\r\n$3\r\nSET\r\n+mangos\r\n:+125\r\n-Error message\r\n$3\r\n100\r\n*5\r\n$3\r\nSET\r\n+mangos\r\n:+125\r\n-Error message\r\n$3\r\n100\r\n",
			expectedOutput: []ParsedCmd{
				{ValueType: RESP_ENCODING_CONSTANTS.BULK_STRING, Value: "SET"},
				{ValueType: RESP_ENCODING_CONSTANTS.STRING, Value: "mangos"},
				{ValueType: RESP_ENCODING_CONSTANTS.INTEGER, Value: "+125"},
				{ValueType: RESP_ENCODING_CONSTANTS.ERROR, Value: "Error message"},
				{ValueType: RESP_ENCODING_CONSTANTS.BULK_STRING, Value: "100"},
			},
			expectError: false,
		},
		{
			name:  "Valid RESP v2 response with string",
			input: "+Test\r\n",
			expectedOutput: []ParsedCmd{
				{
					ValueType: RESP_ENCODING_CONSTANTS.STRING, Value: "Test",
				},
			},
			expectError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			reader := strings.NewReader(test.input)
			respReader := NewReader(reader)

			output, _, err := respReader.HandleRead()

			if test.expectError {
				if err == nil {
					t.Errorf("expected an error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if !reflect.DeepEqual(output, test.expectedOutput) {
				t.Errorf("output mismatch\nexpected: %v\nactual: %v", test.expectedOutput, output)
			}
		})
	}
}
