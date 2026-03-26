package btype

import (
	"testing"

	"github.com/bytedance/mockey"
	"github.com/smartystreets/goconvey/convey"
)

func TestBencodeType(t *testing.T) {
	mockey.PatchConvey("BencodeType", t, func() {
		TestBencodeString(t)
		TestBencodeInt(t)
		TestBencodeList(t)
		TestBencodeDict(t)
		TestNested(t)
	})
}

func TestBencodeString(t *testing.T) {
	mockey.PatchConvey("BencodeString", func() {
		s := NewBencodeString("hello")
		convey.So(s.Type(), convey.ShouldEqual, BencodeTypeString)
		convey.So(s.String(), convey.ShouldEqual, "hello")
		convey.So(s.Int(), convey.ShouldEqual, 0)
		convey.So(s.List(), convey.ShouldEqual, nil)
		convey.So(s.Dict(), convey.ShouldEqual, nil)
		convey.So(s.Marshal(), convey.ShouldResemble, []byte("5:hello"))
	})
}

func TestBencodeInt(t *testing.T) {
	mockey.PatchConvey("BencodeInt", func() {
		i := NewBencodeInt(123)
		convey.So(i.Type(), convey.ShouldEqual, BencodeTypeInt)
		convey.So(i.String(), convey.ShouldEqual, "")
		convey.So(i.Int(), convey.ShouldEqual, 123)
		convey.So(i.List(), convey.ShouldEqual, nil)
		convey.So(i.Dict(), convey.ShouldEqual, nil)
		convey.So(i.Marshal(), convey.ShouldResemble, []byte("i123e"))
	})
}

func TestBencodeList(t *testing.T) {
	mockey.PatchConvey("BencodeList", func() {
		ls := NewBencodeList()
		convey.So(ls.Type(), convey.ShouldEqual, BencodeTypeList)
		convey.So(ls.String(), convey.ShouldEqual, "")
		convey.So(ls.Int(), convey.ShouldEqual, 0)
		convey.So(ls.List(), convey.ShouldEqual, ls)
		convey.So(ls.Dict(), convey.ShouldEqual, nil)
		convey.So(ls.Marshal(), convey.ShouldResemble, []byte("le"))

		convey.So(ls.Len(), convey.ShouldEqual, 0)

		ls.PushBack(NewBencodeString("world"))
		ls.PushFront(NewBencodeString("hello"))
		ls.Insert(1, NewBencodeString(","))
		convey.So(ls.Marshal(), convey.ShouldResemble, []byte("l5:hello1:,5:worlde"))

		convey.So(ls.Len(), convey.ShouldEqual, 3)
		convey.So(ls.Front(), convey.ShouldEqual, NewBencodeString("hello"))
		convey.So(ls.Index(1), convey.ShouldEqual, NewBencodeString(","))
		convey.So(ls.Back(), convey.ShouldEqual, NewBencodeString("world"))

		ls.PopBack()
		ls.PopFront()
		convey.So(ls.Len(), convey.ShouldEqual, 1)
		convey.So(ls.Index(0), convey.ShouldEqual, NewBencodeString(","))

		ls.Pop(0)
		convey.So(ls.Len(), convey.ShouldEqual, 0)
	})
}

func TestBencodeDict(t *testing.T) {
	mockey.PatchConvey("BencodeDict", func() {
		d := NewBencodeDict()
		convey.So(d.Type(), convey.ShouldEqual, BencodeTypeDict)
		convey.So(d.String(), convey.ShouldEqual, "")
		convey.So(d.Int(), convey.ShouldEqual, 0)
		convey.So(d.List(), convey.ShouldEqual, nil)
		convey.So(d.Dict(), convey.ShouldEqual, d)
		convey.So(d.Marshal(), convey.ShouldResemble, []byte("de"))

		convey.So(d.Len(), convey.ShouldEqual, 0)

		d.Set(NewBencodeString("a"), NewBencodeString("hello"))
		d.Set(NewBencodeString("b"), NewBencodeString(","))
		d.Set(NewBencodeString("c"), NewBencodeString("world"))
		convey.So(d.Marshal(), convey.ShouldResemble, []byte("d1:a5:hello1:b1:,1:c5:worlde"))
		convey.So(d.Len(), convey.ShouldEqual, 3)
		convey.So(d.Get(NewBencodeString("a")), convey.ShouldEqual, NewBencodeString("hello"))
		convey.So(d.Index(1), convey.ShouldEqual, NewBencodeString(","))
		convey.So(d.Get(NewBencodeString("c")), convey.ShouldEqual, NewBencodeString("world"))

		convey.So(d.Has(NewBencodeString("a")), convey.ShouldEqual, true)
		convey.So(d.Has(NewBencodeString("world")), convey.ShouldEqual, false)

		keys := d.Keys()
		convey.So(keys, convey.ShouldResemble, []BencodeString{"a", "b", "c"})

		d.Remove(NewBencodeString("a"))
		convey.So(d.Len(), convey.ShouldEqual, 2)
		convey.So(d.Has(NewBencodeString("a")), convey.ShouldEqual, false)

		keys = d.Keys()
		convey.So(keys, convey.ShouldResemble, []BencodeString{"b", "c"})

		d.Clear()
		convey.So(d.Len(), convey.ShouldEqual, 0)

		d.Set(NewBencodeString("a"), NewBencodeString("hello"))
		convey.So(d.Len(), convey.ShouldEqual, 1)
		convey.So(d.Has(NewBencodeString("a")), convey.ShouldEqual, true)
		convey.So(d.Get(NewBencodeString("a")), convey.ShouldEqual, NewBencodeString("hello"))
	})
}

func TestNested(t *testing.T) {
	mockey.PatchConvey("Nested", func() {
		hello := NewBencodeString("hello")
		comma := NewBencodeString(",")
		world := NewBencodeString("world")

		ls := NewBencodeList(hello, comma, world)
		convey.So(ls.Len(), convey.ShouldEqual, 3)

		d := NewBencodeDict()
		d.Set(NewBencodeString("ls"), ls)
		d.Set(NewBencodeString("num"), NewBencodeInt(123))

		d2 := NewBencodeDict()
		d2.Set(NewBencodeString("dic"), d)
		convey.So(d2.Marshal(), convey.ShouldResemble, []byte("d3:dicd2:lsl5:hello1:,5:worlde3:numi123eee"))

		ss := d2.Get(NewBencodeString("dic"))
		convey.So(ss, convey.ShouldEqual, d)

		convey.So(ss.Type(), convey.ShouldEqual, BencodeTypeDict)
		convey.So(ss.Dict(), convey.ShouldEqual, d)
		convey.So(ss.List(), convey.ShouldEqual, nil)
		convey.So(ss.Int(), convey.ShouldEqual, 0)
		convey.So(ss.String(), convey.ShouldEqual, "")

		sd := ss.Dict()
		convey.So(sd.Len(), convey.ShouldEqual, 2)
		convey.So(sd.Has(NewBencodeString("ls")), convey.ShouldEqual, true)
		convey.So(sd.Has(NewBencodeString("num")), convey.ShouldEqual, true)

		lss := sd.Get(NewBencodeString("ls")).List()
		convey.So(lss.Type(), convey.ShouldEqual, BencodeTypeList)

		lsd := lss.List()
		convey.So(lsd, convey.ShouldEqual, ls)
		convey.So(lsd.Len(), convey.ShouldEqual, 3)
		convey.So(lsd.Front(), convey.ShouldEqual, hello)
		convey.So(lsd.Index(1), convey.ShouldEqual, comma)
		convey.So(lsd.Back(), convey.ShouldEqual, world)
	})
}
