package codec

import (
	"encoding/base64"
)

// likelyBase64Chars is a set of characters that you would expect to find at
// least one of in base64 encoded data. This risks missing about 1% of
// base64 encoded data that doesn't contain these characters, but gives you
// the performance gain of not trying to decode a lot of long symbols in code.
var likelyBase64Chars = make([]bool, 256)

func init() {
	for _, c := range `0123456789+/-_` {
		likelyBase64Chars[c] = true
	}
}

// base64Encodings are the four RFC 4648 variants: both alphabets (standard
// `+/` and URL-safe `-_`), each padded and unpadded. The encodings that share
// an alphabet are mutually exclusive on padding -- the padded form rejects a
// length that is not a multiple of four, the raw form rejects `=` -- so at most
// one member of each alphabet pair can succeed and the order only decides which
// alphabet wins on input that is valid under both (which decodes identically).
var base64Encodings = [...]*base64.Encoding{
	base64.StdEncoding,
	base64.RawStdEncoding,
	base64.URLEncoding,
	base64.RawURLEncoding,
}

// decodeBase64 decodes base64 encoded printable ASCII characters
func decodeBase64(encodedValue string) string {
	// Exit early if it doesn't seem like base64
	if !hasByte(encodedValue, likelyBase64Chars) {
		return ""
	}

	for _, encoding := range base64Encodings {
		decodedValue, err := encoding.DecodeString(encodedValue)
		if err == nil && isPrintableASCII(decodedValue) {
			return string(decodedValue)
		}
	}

	return ""
}
