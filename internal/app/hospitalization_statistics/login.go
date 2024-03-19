package hs

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	UA = "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Mobile Safari/537.36"
)

func newRequest(method string, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", UA)
	return req, nil
}

type LoginArgs struct {
	CorpId            string `json:"corpId"`
	Password          string `json:"password"`
	PhoneNum          string `json:"phoneNum"`
	T                 string `json:"t"`
	InvokerChannel    string `json:"invokerChannel"`
	InvokerDeviceType string `json:"invokerDeviceType"`
	InvokerAppVersion string `json:"invokerAppVersion"`
	UnionId           string `json:"unionId"`
	Callback          string `json:"callback"`
}

type LoginRespBody struct {
	Msg        string `json:"msg"`
	ResultCode string `json:"result_code"`
	StartTime  uint64 `json:"start_time"`
	Success    bool   `json:"success"`
	TimeConsum uint   `json:"time_consum"`
	TraceId    string `json:"trace_id"`
}

func Login() error {
	var args LoginArgs
	data, err := os.ReadFile("tmp/hs/login.json")
	if err != nil {
		return fmt.Errorf("os.ReadFile(\"tmp/hs/login.json\") error. %w", err)
	}
	if err := json.Unmarshal(data, &args); err != nil {
		return fmt.Errorf("json.Unmarshal error. %w", err)
	}
	params := url.Values{
		"corpId":            []string{args.CorpId},
		"password":          []string{args.Password},
		"phoneNum":          []string{args.PhoneNum},
		"t":                 []string{args.T},
		"invokerChannel":    []string{args.InvokerChannel},
		"invokerDeviceType": []string{args.InvokerDeviceType},
		"invokerAppVersion": []string{args.InvokerAppVersion},
		"unionId":           []string{args.UnionId},
		"callback":          []string{args.Callback},
	}
	Url := url.URL{
		Scheme:   "https",
		Host:     "user-xinyang.yuantutech.com",
		Path:     "/user-web/restapi/common/ytUsers/login",
		RawQuery: params.Encode(),
	}
	req, err := newRequest("GET", Url.String(), nil)
	if err != nil {
		return fmt.Errorf("newRequest error. %w", err)
	}

	c := http.Client{
		Timeout: 3 * time.Second,
	}
	resp, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("c.Do error. %w", err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("login error. StatusCode: %d", resp.StatusCode)
	}

	bodyB, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("io.ReadAll error. %w", err)
	}
	var respBody LoginRespBody
	if err := json.Unmarshal(extractJSONPResponse(bodyB), &respBody); err != nil {
		return fmt.Errorf("json.Unmarshal error. %w", err)
	}
	if !respBody.Success {
		return fmt.Errorf("%s", respBody.Msg)
	}

	var cookies = []*http.Cookie{}
	for _, i := range resp.Header.Values("Set-Cookie") {
		c, err := parseCookie(i)
		if err != nil {
			fmt.Printf("parseCookie error. %s. c: %s\n", err, c)
			continue
		}
		if c.Domain != "" && c.Domain != ".yuantutech.com" {
			continue
		}
		cookies = append(cookies, c)
	}

	b, err := json.Marshal(cookies)
	if err != nil {
		return fmt.Errorf("json.Marshal error. %w", err)
	}
	if err := os.WriteFile("tmp/hs/cookies.json", b, 0666); err != nil {
		return fmt.Errorf("os.WriteFile(\"tmp/hs/cookies.json\") error. %w", err)
	}
	return nil
}

func parseCookie(c string) (*http.Cookie, error) {
	// 拆分Set-Cookie字符串
	parts := strings.Split(c, "; ")

	// 解析第一个部分为键值对
	cookieParts := strings.SplitN(parts[0], "=", 2)
	if len(cookieParts) != 2 {
		return nil, fmt.Errorf("invalid Set-Cookie format")
	}

	// 创建http.Cookie对象
	cookie := &http.Cookie{
		Name:  strings.TrimSpace(cookieParts[0]),
		Value: strings.TrimSpace(cookieParts[1]),
	}

	// 解析其他部分的属性
	for i := 1; i < len(parts); i++ {
		attrParts := strings.SplitN(parts[i], "=", 2)
		if len(attrParts) == 2 {
			attrKey := strings.TrimSpace(attrParts[0])
			attrValue := strings.TrimSpace(attrParts[1])

			switch attrKey {
			case "Domain":
				cookie.Domain = attrValue
			case "Path":
				cookie.Path = attrValue
			case "Expires":
				expiration, err := time.Parse("Mon, 02-Jan-2006 15:04:05 MST", attrValue)
				if err == nil {
					cookie.Expires = expiration
				} else {
					fmt.Printf("http.ParseTime error. %s. attrValue: %s\n", err, attrValue)
				}
			case "Max-Age":
				maxAge, err := strconv.Atoi(attrValue)
				if err == nil {
					cookie.MaxAge = maxAge
				}
			case "Secure":
				cookie.Secure = true
			case "HttpOnly":
				cookie.HttpOnly = true
			}
		}
	}

	return cookie, nil
}
