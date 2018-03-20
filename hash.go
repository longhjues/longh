package longh

import (
	"crypto"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

// Hash 散列函数 返回空值则为错误
func Hash(ht crypto.Hash, b []byte) []byte {
	switch ht {
	case crypto.MD5:
		h := md5.New()
		h.Write(b)
		return h.Sum(nil)
	case crypto.SHA1:
		h := sha1.New()
		h.Write(b)
		return h.Sum(nil)
	case crypto.SHA256:
		h := sha256.New()
		h.Write(b)
		return h.Sum(nil)
	default:
		panic(errors.New("the hash is not support"))
	}
	return nil
}

// HashHex can hash and encode to string
func HashHex(ht crypto.Hash, b []byte) string {
	return hex.EncodeToString(Hash(ht, b))
}

// HMACSHA1Hash hmac-sha1
func HMACSHA1Hash(src, key []byte) []byte {
	h := hmac.New(sha1.New, key)
	h.Write(src)
	return h.Sum(nil)
}

// HMACSHA1HexHash hmac-sha1 and encode to string
func HMACSHA1HexHash(src, key []byte) string {
	h := hmac.New(sha1.New, key)
	h.Write(src)
	return hex.EncodeToString(h.Sum(nil))
}
