package hashs

import "github.com/speps/go-hashids"

type HashID struct {
	salt      string
	minLength int
	hd        *hashids.HashIDData
}

func NewHashID(salt string, minLength int) *HashID {
	hd := hashids.NewData()
	hd.Salt = salt
	hd.MinLength = minLength
	return &HashID{
		salt:      salt,
		minLength: minLength,
		hd:        hd,
	}
}
func (h *HashID) Encode(id int64) string {
	hd, _ := hashids.NewWithData(h.hd)
	r, _ := hd.EncodeInt64([]int64{id})
	return r
}
func (h *HashID) Decode(id string) (int64, error) {
	hd, _ := hashids.NewWithData(h.hd)
	r, err := hd.DecodeInt64WithError(id)
	return r[0], err
}
