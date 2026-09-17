package codec

import (
	"encoding/base64"
	"testing"
)

// The four RFC 4648 variants must all decode. Before the fix only StdEncoding
// (padded, `+/`) and RawURLEncoding (unpadded, `-_`) were tried, so a value
// whose encoding used the URL alphabet *and* carried `=` padding -- or the
// standard alphabet with no padding -- was handed back undecoded and every
// rule ran against the ciphertext.
func TestDecodeBase64AllRFC4648Variants(t *testing.T) {
	// `?` at an offset of 2 mod 3 is what puts a 62/63 alphabet character in
	// the output: printable ASCII can only reach index 62/63 in the last
	// character of a group.
	const plaintext = "GET /v1/x?token=ghp_0123456789abcdefghijklmnopqrstuvwxyz"

	for name, encoding := range map[string]*base64.Encoding{
		"std":     base64.StdEncoding,
		"std raw": base64.RawStdEncoding,
		"url":     base64.URLEncoding,
		"url raw": base64.RawURLEncoding,
	} {
		t.Run(name, func(t *testing.T) {
			encoded := encoding.EncodeToString([]byte(plaintext))
			if got := decodeBase64(encoded); got != plaintext {
				t.Fatalf("decodeBase64(%q) = %q, want %q", encoded, got, plaintext)
			}
		})
	}
}
