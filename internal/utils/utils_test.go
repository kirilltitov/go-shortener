package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidLink(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{
			name:  "valid http",
			value: "http://ya.ru",
			want:  true,
		},
		{
			name:  "valid https",
			value: "https://ya.ru",
			want:  true,
		},
		{
			name:  "invalid",
			value: "ya.ru",
			want:  false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, IsValidLink(test.value))
		})
	}
}
