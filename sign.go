package longh

import (
	"crypto"
	"crypto/md5"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

// ValidateSignForMD5 验证签名
// 签名规则为所有参数字典序追加":appkey"
// 最后md5
// 时间暂设为可有可无 缓慢过渡
func ValidateSignForMD5(r *http.Request, appKey string) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	t := r.FormValue("Time")
	if t != "" {
		unix, err := strconv.ParseInt(t, 10, 64)
		if err != nil {
			return err
		}
		if time.Now().Unix()-unix > 60 {
			return errors.New("overtime")
		}
	}

	var baseStr = SortAndAppendMap(r.Form, "", true, "&", "=", ":"+appKey, "sign", "Sign", "act", "Act")
	baseStr = strings.Replace(baseStr, " ", "+", -1)
	var sign = r.FormValue("Sign")
	if sign == "" {
		sign = r.FormValue("sign")
	}
	if !strings.EqualFold(HashHex(crypto.MD5, []byte(baseStr)), sign) {
		return fmt.Errorf("check sign failed:%s, %s", baseStr, HashHex(crypto.MD5, []byte(baseStr)))
	}
	return nil
}

// ProduceSignForMD5 生成签名 自家用
func ProduceSignForMD5(form map[string][]string, appKey string) string {
	if form == nil {
		return ""
	}
	form["Time"] = []string{strconv.FormatInt(time.Now().Unix(), 10)}
	var baseStr = SortAndAppendMap(form, "", true, "&", "=", ":"+appKey, "sign", "Sign", "act", "Act")
	baseStr = strings.Replace(baseStr, " ", "+", -1)
	return HashHex(crypto.MD5, []byte(baseStr))
}

// SortAndAppendMap 简化拼接字符串
// 举个栗子
// baseStr := SortAndAppendMap(r.PostForm, "", true, "&", "=", ":appid", "sign")
// strings.EqualFold(MD5HexHash([]byte(baseStr)), r.FormValue("sign"))
func SortAndAppendMap(value url.Values, strHd string, haveKey bool, strAnd, strEq, strAp string, excepts ...string) string {
	if value == nil {
		return ""
	}
	var s = make([]string, 0, len(value))
	for k := range value {
		var b = true
		for _, v := range excepts {
			if k == v {
				b = false
				break
			}
		}
		if b {
			s = append(s, k)
		}
	}
	sort.Strings(s)
	var buf strings.Builder
	var b = true
	buf.WriteString(strHd)
	for _, v := range s {
		if value[v] == nil {
			continue
		}
		if b {
			if haveKey {
				buf.WriteString(v + strEq)
			}
			buf.WriteString(value[v][0])
			b = false
			continue
		}
		buf.WriteString(strAnd)
		if haveKey {
			buf.WriteString(v + strEq)
		}
		buf.WriteString(value[v][0])
	}
	buf.WriteString(strAp)
	return buf.String()
}

// VerifySHA1WithRSA 验证SHA1WithRSA签名
// Key 与 Sign 为base64后的值
func VerifySHA1WithRSA(base64Key, src, base64Sign string) error {
	key, err := base64.StdEncoding.DecodeString(base64Key)
	if err != nil {
		return err
	}
	sign, err := base64.StdEncoding.DecodeString(base64Sign)
	if err != nil {
		return err
	}

	pubInterface, err := x509.ParsePKIXPublicKey(key)
	if err != nil {
		return err
	}
	pub := pubInterface.(*rsa.PublicKey)
	h := sha1.New()
	h.Write([]byte(src))
	hashed := h.Sum(nil)

	return rsa.VerifyPKCS1v15(pub, crypto.SHA1, hashed, sign)
}

// VerifySHA256WithRSA 验证SHA256WithRSA签名
// Key 与 Sign 为base64后的值
func VerifySHA256WithRSA(base64Key, src, base64Sign string) error {
	key, err := base64.StdEncoding.DecodeString(base64Key)
	if err != nil {
		return err
	}
	sign, err := base64.StdEncoding.DecodeString(base64Sign)
	if err != nil {
		return err
	}

	pubInterface, err := x509.ParsePKIXPublicKey(key)
	if err != nil {
		return err
	}
	pub := pubInterface.(*rsa.PublicKey)
	h := sha256.New()
	h.Write([]byte(src))
	hashed := h.Sum(nil)

	return rsa.VerifyPKCS1v15(pub, crypto.SHA256, hashed, sign)
}

// SignSHA1WithRSA 生成SHA1WithRSA签名
// Key 与 Sign 为base64后的值
func SignSHA1WithRSA(base64Key, src string) (string, error) {
	priv, err := ParsePrivateKey(base64Key)
	if err != nil {
		return "", err
	}

	h := sha1.New()
	h.Write([]byte(src))
	hashed := h.Sum(nil)

	b, err := rsa.SignPKCS1v15(nil, priv, crypto.SHA1, hashed)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(b), nil
}

// ParsePrivateKey 解析私钥
// 自带base64
// 先解析是否为PKCS1格式的
// 再解析是否为PKCS8格式的
func ParsePrivateKey(base64Key string) (*rsa.PrivateKey, error) {
	key, err := base64.StdEncoding.DecodeString(base64Key)
	if err != nil {
		return nil, err
	}

	var priv *rsa.PrivateKey
	priv, err = x509.ParsePKCS1PrivateKey(key)
	if err == nil {
		return priv, nil
	}

	pubInterface, err := x509.ParsePKCS8PrivateKey(key)
	if err != nil {
		return nil, err
	}
	var ok bool
	priv, ok = pubInterface.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("pubInterface is not *rsa.PrivateKey")
	}
	return priv, err
}

// SignSHA256WithRSA 生成SHA256WithRSA签名
// Key 与 Sign 为base64后的值
func SignSHA256WithRSA(base64Key, src string) (string, error) {
	priv, err := ParsePrivateKey(base64Key)
	if err != nil {
		return "", err
	}

	h := sha256.New()
	h.Write([]byte(src))
	hashed := h.Sum(nil)

	b, err := rsa.SignPKCS1v15(nil, priv, crypto.SHA256, hashed)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(b), nil
}

// SignMD5WithRSA 验证MD5WithRSA签名
// Key 与 返回值Sign 为base64后的值
func SignMD5WithRSA(base64Key, src string) (string, error) {
	priv, err := ParsePrivateKey(base64Key)
	if err != nil {
		return "", err
	}
	h := md5.New()
	h.Write([]byte(src))
	hashed := h.Sum(nil)

	sign, err := rsa.SignPKCS1v15(nil, priv, crypto.MD5, hashed)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(sign), nil
}

// VerifyMD5WithRSA 验证MD5WithRSA签名
// Key 与 Sign 为base64后的值
func VerifyMD5WithRSA(base64Key, src, base64Sign string) error {
	key, err := base64.StdEncoding.DecodeString(base64Key)
	if err != nil {
		return err
	}
	sign, err := base64.StdEncoding.DecodeString(base64Sign)
	if err != nil {
		return err
	}

	pubInterface, err := x509.ParsePKIXPublicKey(key)
	if err != nil {
		return err
	}
	pub := pubInterface.(*rsa.PublicKey)
	h := md5.New()
	h.Write([]byte(src))
	hashed := h.Sum(nil)

	return rsa.VerifyPKCS1v15(pub, crypto.MD5, hashed, sign)
}
