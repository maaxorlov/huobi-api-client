package huobiapi

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/url"
	"strings"
	"time"
)

func convertSymbol(symbol string) string {
	return strings.ToLower(strings.Split(symbol, "_")[1] + strings.Split(symbol, "_")[0])
}

func (A *HuobiApi) sign(method, path, parameters string) string {
	if method == "" || A.host == "" || path == "" || parameters == "" {
		return ""
	}
	hash := hmac.New(sha256.New, []byte(A.apiSecret))
	hash.Write(getHashBytes(method, A.host, path, parameters))
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

func getHashBytes(method, host, path, parameters string) []byte {
	var hb bytes.Buffer
	hb.WriteString(method)
	hb.WriteString("\n")
	hb.WriteString(host)
	hb.WriteString("\n")
	hb.WriteString(path)
	hb.WriteString("\n")
	hb.WriteString(parameters)

	return hb.Bytes()
}

func (A *HuobiApi) buildPrivateURI(method, path string, additionalURLParams ...map[string]string) string {
	var params = A.getPrivateRequestString(method, path, additionalURLParams...)
	var uri strings.Builder
	uri.Grow(len(A.root) + len(path) + len(params) + 1)
	uri.WriteString(A.root)
	uri.WriteString(path)
	uri.WriteRune('?')
	uri.WriteString(params)
	/*/
	 * можно было бы просто return A.root + path + "?" + A.getPrivateRequestString(method, path, additionalURLParams...),
	 * но по бенчмарк тестам это было бы медленнее
	/*/

	return uri.String()
}

func (A *HuobiApi) getPrivateRequestString(method, path string, params ...map[string]string) string {
	tValue := time.Now().UTC().Format("2006-01-02T15:04:05")
	req := url.Values{}
	req.Add(A.akKey, A.akValue)
	req.Add(A.smKey, A.smValue)
	req.Add(A.svKey, A.svValue)
	req.Add(A.tKey, tValue)
	if len(params) != 0 {
		for key, value := range params[0] {
			req.Add(key, value)
		}
	}
	req.Add(A.sKey, A.sign(method, path, req.Encode()))

	return req.Encode()
}

func (A *HuobiApi) buildPublicURI(path string, additionalURLParams ...map[string]string) string {
	var params = A.getPublicRequestString(additionalURLParams...)
	var uri strings.Builder
	if params == "" {
		uri.Grow(len(A.root) + len(path))
	} else {
		uri.Grow(len(A.root) + len(path) + len(params) + 1)
	}
	uri.WriteString(A.root)
	uri.WriteString(path)
	if params != "" {
		uri.WriteRune('?')
		uri.WriteString(params)
	}
	/*/
	 * можно было бы просто return A.root + path + "?" + A.getPublicRequestString(path, additionalURLParams...),
	 * но по бенчмарк тестам это было бы медленнее
	/*/

	return uri.String()
}

func (A *HuobiApi) getPublicRequestString(params ...map[string]string) string {
	req := url.Values{}
	if len(params) != 0 {
		for key, value := range params[0] {
			req.Add(key, value)
		}
	}

	return req.Encode()
}
