package btype

import "strconv"

var _ BencodeValue = (*BencodeInt)(nil)

type BencodeInt int64

func NewBencodeInt(i int64) BencodeInt {
	return BencodeInt(i)
}

// region BencodeValue interface implementation

func (i BencodeInt) Type() BencodeType {
	return BencodeTypeInt
}

func (i BencodeInt) String() string {
	return ""
}

func (i BencodeInt) Int() int64 {
	return int64(i)
}

func (i BencodeInt) List() BencodeList {
	return nil
}

func (i BencodeInt) Dict() BencodeDict {
	return nil
}

func (i BencodeInt) Marshal() []byte {
	res := "i" + strconv.FormatInt(int64(i), 10) + "e"
	return []byte(res)
}

// endregion
