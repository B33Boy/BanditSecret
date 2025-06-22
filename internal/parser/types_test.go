package parser

import (
	"testing"
)

func TestUnmarshalJSON(t *testing.T) {

	tests := []struct {
		name    string
		input   string
		want    TimeMs
		wantErr bool
	}{
		{
			name:  "Valid time",
			input: `"01:02:03.004"`,
			want: 1*(60*60*1000) +
				2*(60*1000) +
				3*1000 +
				4,
			wantErr: false,
		},
		{
			name:    "Invalid time",
			input:   `invalid`,
			want:    0,
			wantErr: true,
		},
		{
			name:    "Invalid time, unquoted",
			input:   `1234`,
			want:    1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got TimeMs
			err := got.UnmarshalJSON([]byte(tt.input))

			// If we have an error, check for this test if we wantErr => if not, then fail test
			if (err != nil) != tt.wantErr {
				t.Fatalf("expected error = %v, got %v", tt.wantErr, err)
			}

			// In the case that we don't have an error
			// AND we don't want an error
			// AND what we received is not equal to what we want
			if !tt.wantErr && got != tt.want {
				t.Errorf("expected %d, got %d", tt.want, got)
			}
		})

	}

}
