package btype

var _ BencodeList = (*blist)(nil)

type BencodeList interface {
	BencodeValue

	Len() int
	PushBack(v BencodeValue)
	PushFront(v BencodeValue)
	Insert(index int, v BencodeValue)
	Front() BencodeValue
	Back() BencodeValue
	Index(index int) BencodeValue
	PopFront()
	PopBack()
	Pop(index int)
	Clear()
}

type blist struct {
	values []BencodeValue
}

func NewBencodeList(values ...BencodeValue) BencodeList {
	ls := &blist{
		values: make([]BencodeValue, 0, len(values)),
	}

	if len(values) > 0 {
		for _, v := range values {
			ls.values = append(ls.values, v)
		}
	}

	return ls
}

// region BencodeValue interface implementation

func (ls *blist) Type() BencodeType {
	return BencodeTypeList
}

func (ls *blist) String() string {
	return ""
}

func (ls *blist) Int() int64 {
	return 0
}

func (ls *blist) List() BencodeList {
	return ls
}

func (ls *blist) Dict() BencodeDict {
	return nil
}

func (ls *blist) Marshal() []byte {
	res := []byte{'l'}
	for _, v := range ls.values {
		res = append(res, v.Marshal()...)
	}

	res = append(res, 'e')
	return res
}

// endregion

// region BencodeList interface implementation

func (ls *blist) Len() int {
	return len(ls.values)
}

func (ls *blist) Empty() bool {
	return ls.Len() == 0
}

func (ls *blist) PushBack(v BencodeValue) {
	ls.values = append(ls.values, v)
}

func (ls *blist) PushFront(v BencodeValue) {
	ls.values = append([]BencodeValue{v}, ls.values...)
}

func (ls *blist) Insert(index int, v BencodeValue) {
	ls.values = append(ls.values[:index], append([]BencodeValue{v}, ls.values[index:]...)...)
}

func (ls *blist) Front() BencodeValue {
	return ls.values[0]
}

func (ls *blist) Back() BencodeValue {
	return ls.values[ls.Len()-1]
}

func (ls *blist) Index(index int) BencodeValue {
	return ls.values[index]
}

func (ls *blist) PopFront() {
	ls.values = ls.values[1:]
}

func (ls *blist) PopBack() {
	ls.values = ls.values[:ls.Len()-1]
}

func (ls *blist) Pop(index int) {
	ls.values = append(ls.values[:index], ls.values[index+1:]...)
}

func (ls *blist) Clear() {
	ls.values = make([]BencodeValue, 0)
}

// endregion
