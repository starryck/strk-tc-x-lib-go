package gbradix

import (
	"bytes"
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"math"
	"strings"
)

func Base16EncodeBtob(src []byte) []byte {
	dst := make([]byte, len(src)*2)
	hex.Encode(dst, src)
	return dst
}

func Base16DecodeBtob(src []byte) ([]byte, error) {
	dst := make([]byte, len(src)/2)
	_, err := hex.Decode(dst, src)
	return dst, err
}

func Base16EncodeBtoa(src []byte) string {
	dst := Base16EncodeBtob(src)
	return string(dst)
}

func Base16DecodeAtob(src string) ([]byte, error) {
	dst, err := Base16DecodeBtob([]byte(src))
	return dst, err
}

func Base16EncodeAtoa(src string) string {
	dst := Base16EncodeBtoa([]byte(src))
	return dst
}

func Base16DecodeAtoa(src string) (string, error) {
	dst, err := Base16DecodeAtob(src)
	return string(dst), err
}

func Base32EncodeBtob(src []byte) []byte {
	dst := make([]byte, int(math.Ceil(float64(len(src))/5)*8))
	base32.StdEncoding.Encode(dst, src)
	return dst
}

func Base32DecodeBtob(src []byte) ([]byte, error) {
	src = bytes.TrimRight(src, string(base32.StdPadding))
	dst := make([]byte, int(math.Floor(float64(len(src))/8*5)))
	_, err := base32.StdEncoding.WithPadding(base32.NoPadding).Decode(dst, src)
	return dst, err
}

func Base32EncodeBtoa(src []byte) string {
	dst := Base32EncodeBtob(src)
	return string(dst)
}

func Base32DecodeAtob(src string) ([]byte, error) {
	dst, err := Base32DecodeBtob([]byte(src))
	return dst, err
}

func Base32EncodeAtoa(src string) string {
	dst := Base32EncodeBtoa([]byte(src))
	return dst
}

func Base32DecodeAtoa(src string) (string, error) {
	dst, err := Base32DecodeAtob(src)
	return string(dst), err
}

func Base32URLEncodeBtob(src []byte) []byte {
	dst := make([]byte, int(math.Ceil(float64(len(src))/5*8)))
	base32.StdEncoding.WithPadding(base32.NoPadding).Encode(dst, src)
	return dst
}

func Base32URLDecodeBtob(src []byte) ([]byte, error) {
	dst := make([]byte, int(math.Floor(float64(len(src))/8*5)))
	_, err := base32.StdEncoding.WithPadding(base32.NoPadding).Decode(dst, src)
	return dst, err
}

func Base32URLEncodeBtoa(src []byte) string {
	dst := Base32URLEncodeBtob(src)
	return string(dst)
}

func Base32URLDecodeAtob(src string) ([]byte, error) {
	dst, err := Base32URLDecodeBtob([]byte(src))
	return dst, err
}

func Base32URLEncodeAtoa(src string) string {
	dst := Base32URLEncodeBtoa([]byte(src))
	return dst
}

func Base32URLDecodeAtoa(src string) (string, error) {
	dst, err := Base32URLDecodeAtob(src)
	return string(dst), err
}

func Base64EncodeBtob(src []byte) []byte {
	dst := make([]byte, int(math.Ceil(float64(len(src))/3)*4))
	base64.StdEncoding.Encode(dst, src)
	return dst
}

func Base64DecodeBtob(src []byte) ([]byte, error) {
	src = bytes.TrimRight(src, string(base64.StdPadding))
	dst := make([]byte, int(math.Floor(float64(len(src))/4*3)))
	_, err := base64.RawStdEncoding.Decode(dst, src)
	return dst, err
}

func Base64EncodeBtoa(src []byte) string {
	dst := Base64EncodeBtob(src)
	return string(dst)
}

func Base64DecodeAtob(src string) ([]byte, error) {
	dst, err := Base64DecodeBtob([]byte(src))
	return dst, err
}

func Base64EncodeAtoa(src string) string {
	dst := Base64EncodeBtoa([]byte(src))
	return dst
}

func Base64DecodeAtoa(src string) (string, error) {
	dst, err := Base64DecodeAtob(src)
	return string(dst), err
}

func Base64URLEncodeBtob(src []byte) []byte {
	dst := make([]byte, int(math.Ceil(float64(len(src))/3*4)))
	base64.RawURLEncoding.Encode(dst, src)
	return dst
}

func Base64URLDecodeBtob(src []byte) ([]byte, error) {
	dst := make([]byte, int(math.Floor(float64(len(src))/4*3)))
	_, err := base64.RawURLEncoding.Decode(dst, src)
	return dst, err
}

func Base64URLEncodeBtoa(src []byte) string {
	dst := Base64URLEncodeBtob(src)
	return string(dst)
}

func Base64URLDecodeAtob(src string) ([]byte, error) {
	dst, err := Base64URLDecodeBtob([]byte(src))
	return dst, err
}

func Base64URLEncodeAtoa(src string) string {
	dst := Base64URLEncodeBtoa([]byte(src))
	return dst
}

func Base64URLDecodeAtoa(src string) (string, error) {
	dst, err := Base64URLDecodeAtob(src)
	return string(dst), err
}

var (
	Base10Bytes    = []byte("0123456789")
	Base16Bytes    = []byte("0123456789ABCDEF")
	Base36Bytes    = []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	Base62Bytes    = []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
	Base10IndexMap = map[byte]int64{}
	Base16IndexMap = map[byte]int64{}
	Base36IndexMap = map[byte]int64{}
	Base62IndexMap = map[byte]int64{}
)

func init() {
	initBaseNIndexes(Base10Bytes, Base10IndexMap)
	initBaseNIndexes(Base16Bytes, Base16IndexMap)
	initBaseNIndexes(Base36Bytes, Base36IndexMap)
	initBaseNIndexes(Base62Bytes, Base62IndexMap)
}

func initBaseNIndexes(radixes []byte, indexes map[byte]int64) {
	for idx, byt := range radixes {
		indexes[byt] = int64(idx)
	}
}

func BaseNEncodeItoa(num int64, radixes []byte) string {
	if num == 0 {
		return string(radixes[0])
	}
	base := uint64(len(radixes))
	unum := uint64(math.Abs(float64(num)))
	size := int(math.Log10(float64(unum))/math.Log10(float64(base))) + 1
	rdxs := make([]string, size)
	for i := size - 1; i >= 0; i-- {
		rdxs[i] = string(radixes[unum%base])
		unum = unum / base
	}
	rstr := strings.Join(rdxs, "")
	if num < 0 {
		rstr = "-" + rstr
	}
	return rstr
}

func BaseNDecodeAtoi(rstr string, indexes map[byte]int64) int64 {
	base := int64(len(indexes))
	inum := int64(0)
	sign := int64(1)
	if strings.HasPrefix(rstr, "-") {
		sign = int64(-1)
	}
	for _, rdx := range rstr {
		inum = inum*base + sign*indexes[byte(rdx)]
	}
	return inum
}

func Base10EncodeItoa(num int64) string {
	return BaseNEncodeItoa(num, Base10Bytes)
}

func Base10DecodeAtoi(rstr string) int64 {
	return BaseNDecodeAtoi(rstr, Base10IndexMap)
}

func Base16EncodeItoa(num int64) string {
	return BaseNEncodeItoa(num, Base16Bytes)
}

func Base16DecodeAtoi(rstr string) int64 {
	return BaseNDecodeAtoi(rstr, Base16IndexMap)
}

func Base36EncodeItoa(num int64) string {
	return BaseNEncodeItoa(num, Base36Bytes)
}

func Base36DecodeAtoi(rstr string) int64 {
	return BaseNDecodeAtoi(rstr, Base36IndexMap)
}

func Base62EncodeItoa(num int64) string {
	return BaseNEncodeItoa(num, Base62Bytes)
}

func Base62DecodeAtoi(rstr string) int64 {
	return BaseNDecodeAtoi(rstr, Base62IndexMap)
}
