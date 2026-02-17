package generator

import (
	"testing"
)

func TestGenerateCharset(t *testing.T) {
    s := Generate()
    for _, r := range s {
        if !(('a' <= r && r <= 'z') || ('A' <= r && r <= 'Z') || ('0' <= r && r <= '9') || r == '_') {
            t.Errorf("invalid character %c in generated string", r)
        }
    }
}