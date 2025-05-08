package invite

import (
	"unicode"

	"github.com/speps/go-hashids/v2"
)

const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567"
const codeLength = 6

func initHashID(salt string) (*hashids.HashID, error) {
	hd := hashids.NewData()
	hd.Alphabet = alphabet    // base32 alphabet
	hd.Salt = salt            // random salt
	hd.MinLength = codeLength // length should be 6

	return hashids.NewWithData(hd)
}

func (c *Client) generateCode(number int) (string, error) {
	code, err := c.hashID.Encode([]int{number})
	if err != nil {
		return "", err
	}
	if len(code) != 6 {
		return "", ErrInvalidCodeLength
	}
	return code, nil
}

func isAlphanumeric(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}
