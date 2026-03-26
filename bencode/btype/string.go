package btype

import "strconv"

var _ BencodeValue = (*BencodeString)(nil)

type BencodeString string

func NewBencodeString(s string) BencodeString {
	return BencodeString(s)
}

// region BencodeValue interface implementation

func (s BencodeString) Type() BencodeType {
	return BencodeTypeString
}

func (s BencodeString) String() string {
	return string(s)
}

func (s BencodeString) Int() int64 {
	return 0
}

func (s BencodeString) List() BencodeList {
	return nil
}

func (s BencodeString) Dict() BencodeDict {
	return nil
}

func (s BencodeString) Marshal() []byte {
	l := len(string(s))
	res := strconv.Itoa(l) + ":" + string(s)

	return []byte(res)
}

// endregion
