package invite

import "github.com/speps/go-hashids/v2"

func initHashID(salt string) (*hashids.HashID, error) {
	hd := hashids.NewData()
	hd.Alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567" // base32 alphabet
	hd.Salt = salt                                   // random salt
	hd.MinLength = 6                                 // length should be 6

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
