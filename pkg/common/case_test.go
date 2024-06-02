package common

import "testing"

func TestSnakeCaseToPascalCase(t *testing.T) {
	tests := []struct {
		name string
		text string
		want string
	}{
		{
			name: "empty",
			text: "",
			want: "",
		},
		{
			name: "basic",
			text: "hello_world",
			want: "HelloWorld",
		},
		{
			name: "with_number",
			text: "hello_123_world",
			want: "Hello_123World",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SnakeCaseToPascalCase(tt.text); got != tt.want {
				t.Errorf("SnakeCaseToPascalCase(%q) != %q", got, tt.want)
			}
		})
	}
}
