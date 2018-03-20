package longh

import (
	"crypto"
	"crypto/rsa"
	"encoding/base64"
	"math/big"
	"strconv"
)

// SignRSA 生成RSA签名
func SignRSA(ht crypto.Hash, priv *rsa.PrivateKey, src []byte) ([]byte, error) {
	hashed := Hash(ht, src)
	return rsa.SignPKCS1v15(nil, priv, ht, hashed)
}

// VerifyRSA 验证RSA签名
func VerifyRSA(ht crypto.Hash, pub *rsa.PublicKey, src, sign []byte) error {
	hashed := Hash(ht, src)
	return rsa.VerifyPKCS1v15(pub, ht, hashed, sign)
}

// SignRSAFromBase64 验证RSA签名
// Key 与 返回值Sign 为base64后的值
func SignRSAFromBase64(ht crypto.Hash, base64Key, src string) (string, error) {
	priv, err := ParsePrivateKey(base64Key)
	if err != nil {
		return "", err
	}

	sign, err := SignRSA(ht, priv, []byte(src))
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(sign), nil
}

// PublicKeyMakeFromEN 根据E,N生成一个公钥
func PublicKeyMakeFromEN(e, n []byte) (*rsa.PublicKey, error) {
	var pub rsa.PublicKey
	pub.N = big.NewInt(0)
	pub.N.SetBytes(n)
	var err error
	pub.E, err = strconv.Atoi(string(e))
	if err != nil {
		return nil, err
	}
	return &pub, nil
}
