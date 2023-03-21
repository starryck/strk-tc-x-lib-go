package gbradix

import (
	"bytes"
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"math"
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
	dst := make([]byte, int(math.Ceil(float64(len(src))/8)*5))
	_, err := base32.StdEncoding.Decode(dst, src)
	return bytes.TrimRight(dst, "\x00"), err
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
	return bytes.TrimRight(dst, "\x00"), err
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
	dst := make([]byte, int(math.Floor(float64(len(src))/4)*3))
	_, err := base64.StdEncoding.Decode(dst, src)
	return bytes.TrimRight(dst, "\x00"), err
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
	return bytes.TrimRight(dst, "\x00"), err
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
