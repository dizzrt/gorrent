package btype

var _ BencodeDict = (*dict)(nil)

type BencodeDict interface {
	BencodeValue

	Len() int
	Has(key BencodeString) bool
	Set(key BencodeString, value BencodeValue)
	Get(key BencodeString) BencodeValue
	Index(index int) BencodeValue
	Keys() []BencodeString
	Remove(key BencodeString)
	Clear()
}

type dict struct {
	keys []BencodeString
	data map[BencodeString]BencodeValue
}

func NewBencodeDict() BencodeDict {
	return &dict{
		keys: make([]BencodeString, 0),
		data: make(map[BencodeString]BencodeValue),
	}
}

// region BencodeValue interface implementation

func (d *dict) Type() BencodeType {
	return BencodeTypeDict
}

func (d *dict) String() string {
	return ""
}

func (d *dict) Int() int64 {
	return 0
}

func (d *dict) List() BencodeList {
	return nil
}

func (d *dict) Dict() BencodeDict {
	return d
}

func (d *dict) Marshal() []byte {
	res := []byte{'d'}
	for _, k := range d.keys {
		res = append(res, k.Marshal()...)
		res = append(res, d.data[k].Marshal()...)
	}

	res = append(res, 'e')
	return res
}

// endregion

// region BencodeDict interface implementation

func (d *dict) Len() int {
	return len(d.keys)
}

func (d *dict) Has(key BencodeString) bool {
	_, ok := d.data[key]
	return ok
}

func (d *dict) Set(key BencodeString, value BencodeValue) {
	if !d.Has(key) {
		d.keys = append(d.keys, key)
	}

	d.data[key] = value
}

func (d *dict) Get(key BencodeString) BencodeValue {
	return d.data[key]
}

func (d *dict) Index(index int) BencodeValue {
	return d.data[d.keys[index]]
}

func (d *dict) Keys() []BencodeString {
	return d.keys
}

func (d *dict) Remove(key BencodeString) {
	delete(d.data, key)
	for i := range d.keys {
		if d.keys[i] == key {
			d.keys = append(d.keys[:i], d.keys[i+1:]...)
			break
		}
	}
}

func (d *dict) Clear() {
	d.keys = make([]BencodeString, 0)
	d.data = make(map[BencodeString]BencodeValue)
}

// endregion
