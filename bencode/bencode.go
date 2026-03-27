package bencode

import (
	"errors"
	"strconv"

	"github.com/dizzrt/gorrent/bencode/btype"
)

var (
	ErrInvalidBencodeFormat = errors.New("invalid bencode format")
)

func Unmarshal(data []byte) (btype.BencodeValue, error) {
	v, _, err := unmarshal(data, 0, len(data))
	return v, err
}

func unmarshal(data []byte, i, size int) (bvalue btype.BencodeValue, count int, err error) {
	// unmarshal number
	if data[i] == 'i' {
		bvalue, count, err = unmarshalInt(data, i)
		return
	}

	// unmarshal string
	if data[i] >= '0' && data[i] <= '9' {
		bvalue, count, err = unmarshalString(data, i)
		return
	}

	// unmarshal list
	if data[i] == 'l' {
		ls := btype.NewBencodeList()
		for j := i + 1; j < size; j++ {
			count++
			if data[j] == 'e' {
				break
			}

			v, c, e := unmarshal(data, j, size)
			if e != nil {
				err = e
				return
			}

			j += c
			count += c
			ls.PushBack(v)
		}

		bvalue = ls
		return
	}

	if data[i] == 'd' {
		dic := btype.NewBencodeDict()
		for j := i + 1; j < size; j++ {
			count++
			if data[j] == 'e' {
				break
			}

			key, c, e := unmarshalString(data, j)
			if e != nil {
				err = e
				return
			}

			j += c + 1
			count += c + 1

			v, c, e := unmarshal(data, j, size)
			if e != nil {
				err = e
				return
			}

			j += c
			count += c
			dic.Set(key, v)
		}

		bvalue = dic
		return
	}

	err = ErrInvalidBencodeFormat
	return
}

func unmarshalInt(data []byte, i int) (res btype.BencodeInt, count int, err error) {
	isFirst := true
	isNegative := false

	size := len(data)
	var temp int64 = 0
	for j := i + 1; j < size; j++ {
		count++

		v := data[j]
		if v == '-' {
			if isFirst {
				isNegative = true
				continue
			}

			err = ErrInvalidBencodeFormat
			return
		}

		if v == 'e' {
			break
		}

		if v < '0' || v > '9' {
			err = ErrInvalidBencodeFormat
			return
		}

		temp = temp*10 + int64(v-'0')
		isFirst = false
	}

	if isNegative {
		temp = -temp
	}

	res = btype.NewBencodeInt(temp)
	return
}

func unmarshalString(data []byte, i int) (res btype.BencodeString, count int, err error) {
	lengthBytes := make([]byte, 0, 1)

	j := i
	endFlag := false
	size := len(data)

	for j < size {
		count++

		v := data[j]
		j++

		if v == ':' {
			endFlag = true
			break
		}

		if v < '0' || v > '9' {
			err = ErrInvalidBencodeFormat
			return
		}

		lengthBytes = append(lengthBytes, v)
	}

	if !endFlag {
		err = ErrInvalidBencodeFormat
		return
	}

	length, err := strconv.ParseInt(string(lengthBytes), 10, 64)
	if err != nil {
		return
	}

	end := j + int(length)
	if length <= 0 || end > size {
		err = ErrInvalidBencodeFormat
		return
	}

	s := make([]byte, 0, int(length))
	for j < end {
		count++
		s = append(s, data[j])
		j++
	}

	count--
	res = btype.NewBencodeString(string(s))
	return
}
