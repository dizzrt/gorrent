package btype

type BencodeType int8

const (
	BencodeTypeString BencodeType = iota
	BencodeTypeInt
	BencodeTypeList
	BencodeTypeDict
)

type BencodeValue interface {
	Type() BencodeType
	String() string
	Int() int64
	List() BencodeList
	Dict() BencodeDict
	Marshal() []byte
}
