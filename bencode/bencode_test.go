package bencode

import (
	"testing"

	"github.com/bytedance/mockey"
	"github.com/dizzrt/gorrent/bencode/btype"
	"github.com/smartystreets/goconvey/convey"
)

func TestUnmarshal(t *testing.T) {
	mockey.PatchConvey("Unmarshal", t, func() {
		TestUnmarshalInt(t)
		TestUnmarshalString(t)
		TestUnmarshalList(t)
		TestUnmarshalDict(t)
		TestUnmarshalNested(t)
	})
}

func TestUnmarshalInt(t *testing.T) {
	mockey.PatchConvey("Number", func() {
		data := []byte("i123e")
		v, err := Unmarshal(data)
		convey.So(err, convey.ShouldBeNil)
		convey.So(v, convey.ShouldResemble, btype.NewBencodeInt(123))

		data = []byte("i-123e")
		v, err = Unmarshal(data)
		convey.So(err, convey.ShouldBeNil)
		convey.So(v, convey.ShouldResemble, btype.NewBencodeInt(-123))

		data = []byte("i0e")
		v, err = Unmarshal(data)
		convey.So(err, convey.ShouldBeNil)
		convey.So(v, convey.ShouldResemble, btype.NewBencodeInt(0))
	})
}

func TestUnmarshalString(t *testing.T) {
	mockey.PatchConvey("String", func() {
		data := []byte("5:hello")
		v, err := Unmarshal(data)
		convey.So(err, convey.ShouldBeNil)
		convey.So(v, convey.ShouldResemble, btype.NewBencodeString("hello"))

		data = []byte("0:")
		v, err = Unmarshal(data)
		convey.So(err, convey.ShouldEqual, ErrInvalidBencodeFormat)

		data = []byte("1:")
		v, err = Unmarshal(data)
		convey.So(err, convey.ShouldEqual, ErrInvalidBencodeFormat)

		data = []byte("2:a")
		v, err = Unmarshal(data)
		convey.So(err, convey.ShouldEqual, ErrInvalidBencodeFormat)

		data = []byte("2:abc")
		v, err = Unmarshal(data)
		convey.So(err, convey.ShouldBeNil)
		convey.So(v, convey.ShouldEqual, btype.NewBencodeString("ab"))
	})
}

func TestUnmarshalList(t *testing.T) {
	mockey.PatchConvey("List", func() {
		ls := btype.NewBencodeList()
		ls.PushBack(btype.BencodeInt(123))
		ls.PushBack(btype.BencodeString("abc"))

		data := ls.Marshal()
		v, err := Unmarshal(data)
		convey.So(err, convey.ShouldBeNil)
		convey.So(v, convey.ShouldResemble, ls)

		ls = btype.NewBencodeList()
		ls.PushBack(btype.BencodeInt(1))
		ls.PushBack(btype.BencodeInt(2))
		ls.PushBack(btype.BencodeInt(3))

		data = ls.Marshal()
		v, err = Unmarshal(data)
		convey.So(err, convey.ShouldBeNil)
		convey.So(v, convey.ShouldResemble, ls)

		ls = btype.NewBencodeList()
		ls.PushBack(btype.BencodeString("hello"))
		ls.PushBack(btype.BencodeString(","))
		ls.PushBack(btype.BencodeString("world"))

		data = ls.Marshal()
		v, err = Unmarshal(data)
		convey.So(err, convey.ShouldBeNil)
		convey.So(v, convey.ShouldResemble, ls)

		ls2 := btype.NewBencodeList()
		ls2.PushBack(ls)

		data = ls2.Marshal()
		v, err = Unmarshal(data)
		convey.So(err, convey.ShouldBeNil)
		convey.So(v, convey.ShouldResemble, ls2)

		ls2.PushBack(btype.NewBencodeInt(123))
		ls2.PushBack(btype.NewBencodeString("abc"))
		data = ls2.Marshal()
		v, err = Unmarshal(data)
		convey.So(err, convey.ShouldBeNil)
		convey.So(v, convey.ShouldResemble, ls2)
	})
}

func TestUnmarshalDict(t *testing.T) {
	mockey.PatchConvey("Dict", func() {
		dic := btype.NewBencodeDict()
		dic.Set(btype.BencodeString("a"), btype.BencodeString("hello"))
		dic.Set(btype.BencodeString("b"), btype.BencodeString(","))
		dic.Set(btype.BencodeString("c"), btype.BencodeString("world"))

		data := dic.Marshal()
		v, err := Unmarshal(data)
		convey.So(err, convey.ShouldBeNil)
		convey.So(v, convey.ShouldResemble, dic)

		dic = btype.NewBencodeDict()
		dic.Set(btype.BencodeString("a"), btype.BencodeInt(1))
		dic.Set(btype.BencodeString("a"), btype.BencodeInt(2))
		dic.Set(btype.BencodeString("a"), btype.BencodeInt(3))

		data = dic.Marshal()
		v, err = Unmarshal(data)
		convey.So(err, convey.ShouldBeNil)
		convey.So(v, convey.ShouldResemble, dic)

		dic = btype.NewBencodeDict()
		dic.Set(btype.BencodeString("a"), btype.BencodeString("hello"))
		dic.Set(btype.BencodeString("num"), btype.BencodeInt(123))
		dic.Set(btype.BencodeString("c"), btype.BencodeString("world"))
		dic.Set(btype.BencodeString("num2"), btype.BencodeInt(13))
		dic.Set(btype.BencodeString("num"), btype.BencodeInt(45))

		data = dic.Marshal()
		v, err = Unmarshal(data)
		convey.So(err, convey.ShouldBeNil)
		convey.So(v, convey.ShouldResemble, dic)

		dic2 := btype.NewBencodeDict()
		dic2.Set(btype.BencodeString("alpha"), btype.BencodeString("hello"))
		dic2.Set(btype.BencodeString("dic"), dic)
		dic2.Set(btype.BencodeString("num"), btype.BencodeInt(123))

		data = dic2.Marshal()
		v, err = Unmarshal(data)
		convey.So(err, convey.ShouldBeNil)
		convey.So(v, convey.ShouldResemble, dic2)
	})
}

func TestUnmarshalNested(t *testing.T) {
	mockey.PatchConvey("Nested", func() {
		hello := btype.NewBencodeString("hello")
		comma := btype.NewBencodeString(",")
		world := btype.NewBencodeString("world")
		ls := btype.NewBencodeList(hello, comma, world, btype.BencodeInt(111))

		num := btype.NewBencodeInt(123)
		ls2 := btype.NewBencodeList(ls, num)

		ls3 := btype.NewBencodeList(btype.NewBencodeInt(1), btype.NewBencodeString("xxx"))

		dic := btype.NewBencodeDict()
		dic.Set(btype.NewBencodeString("ls2"), ls2)

		dic3 := btype.NewBencodeDict()
		dic3.Set(btype.NewBencodeString("ls3"), ls3)
		dic3.Set(btype.NewBencodeString("num"), btype.NewBencodeInt(666))
		dic3.Set(btype.NewBencodeString("dic"), dic)
		dic3.Set(btype.NewBencodeString("zzz"), btype.NewBencodeString("zzz"))

		data := dic3.Marshal()
		v, err := Unmarshal(data)
		convey.So(err, convey.ShouldBeNil)
		convey.So(v, convey.ShouldResemble, dic3)
	})
}
