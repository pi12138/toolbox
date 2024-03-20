package hs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type Item struct {
	Cost      int    `json:"cost"`
	DeptName  string `json:"deptName"`
	ItemName  string `json:"itemName"`
	ItemPrice string `json:"itemPrice"`
	ItemQty   string `json:"itemQty"`
	ItemSpecs string `json:"itemSpecs"`
	ItemUnits string `json:"itemUnits"`
	TradeTime string `json:"tradeTime"`
	VisitId   string `json:"visitId"`
}

type Data struct {
	DailyCost int64  `json:"dailyCost"`
	Date      string `json:"date"`
	Items     []Item `json:"items"`
}

type CrawlRespBody struct {
	Msg        string `json:"msg"`
	ResultCode string `json:"resultCode"`
	StartTime  int64  `json:"startTime"`
	Success    bool   `json:"success"`
	TimeConsum int64  `json:"timeConsum"`
	TraceId    string `json:"traceId"`
	Data       Data   `json:"data"`
}

type queryArgs struct {
	CorpId            string `json:"corpId"`
	T                 string `json:"t"`
	InvokerChannel    string `json:"invokerChannel"`
	InvokerDeviceType string `json:"invokerDeviceType"`
	InvokerAppVersion string `json:"invokerAppVersion"`
	UnionId           string `json:"unionId"`
	Callback          string `json:"callback"`
	HosId             string `json:"hosId"`
	Date              string `json:"date"`
}

func readCookies() ([]http.Cookie, error) {
	b, err := os.ReadFile("tmp/hs/cookies.json")
	if err != nil {
		return nil, fmt.Errorf(`[readCookies] os.ReadFile("tmp/hs/cookies.json") error. %w`, err)
	}
	cookies := []http.Cookie{}
	if err := json.Unmarshal(b, &cookies); err != nil {
		return nil, fmt.Errorf("[readCookies] json.Unmarshal error. %w", err)
	}
	return cookies, nil
}

func addCookie(req *http.Request) error {
	var retry bool
	now := time.Now()
	for {
		var isExpired bool
		cookies, err := readCookies()
		if err != nil {
			return fmt.Errorf("readCookies error. %w", err)
		}
		if len(cookies) == 0 {
			isExpired = true
		}

		for _, c := range cookies {
			if !c.Expires.IsZero() && c.Expires.Before(now) {
				isExpired = true
				break
			}
		}

		if isExpired && retry {
			return fmt.Errorf("addCookie error. need check Login")
		}

		if isExpired {
			Login()
			retry = true
			continue
		}

		for _, c := range cookies {
			req.AddCookie(&c)
		}
		return nil
	}
}

func Crawl(date, hosId string) (*CrawlRespBody, error) {
	var args queryArgs
	b, err := os.ReadFile("tmp/hs/crawl.json")
	if err != nil {
		return nil, fmt.Errorf(`os.ReadFile("tmp/hs/crawl.json"). %w`, err)
	}
	if err := json.Unmarshal(b, &args); err != nil {
		return nil, fmt.Errorf(`crawl.json json.Unmarshal error. %w`, err)
	}
	args.Date = date
	args.HosId = hosId
	params := url.Values{}
	params.Add("unionId", args.UnionId)
	params.Add("corpId", args.CorpId)
	params.Add("date", args.Date)
	params.Add("hosId", args.HosId)
	params.Add("t", args.T)
	params.Add("invokerChannel", args.InvokerChannel)
	params.Add("invokerDeviceType", args.InvokerDeviceType)
	params.Add("invokerAppVersion", args.InvokerChannel)
	params.Add("callback", args.Callback)
	Url := url.URL{
		Host:     "user-xinyang.yuantutech.com",
		Path:     "/user-web/restapi/inhos/inhosbilldetailByHosId",
		Scheme:   "https",
		RawQuery: params.Encode(),
	}

	client := http.Client{
		Timeout: 3 * time.Second,
	}
	req, err := newRequest("GET", Url.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("newRequest error. %w", err)
	}
	if err := addCookie(req); err != nil {
		return nil, fmt.Errorf("addCookie error. %w", err)
	}

	// print cookie
	// for _, i := range req.Cookies() {
	// 	fmt.Println(`cookie:`, i)
	// }
	// fmt.Println(req.Header.Values("Cookie"))

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request error. %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("StatusCode is not 200. StatusCode: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll body error. %w", err)
	}

	switch getContentType(resp.Header) {
	case "application/javascript":
		data := extractJSONPResponse(body)
		if len(data) == 0 {
			return nil, fmt.Errorf("extractJSONPResponse error. no data. Body: %s", body)
		}
		var Body CrawlRespBody
		if err := json.Unmarshal(data, &Body); err != nil {
			return nil, fmt.Errorf("body json.Unmarshal error. %w", err)
		}
		if !Body.Success {
			return nil, fmt.Errorf("%s", Body.Msg)
		}

		return &Body, nil
	default:
		return nil, fmt.Errorf("unhandler content-type. %s", getContentType(resp.Header))
	}
}

func extractJSONPResponse(jsonpResponse []byte) []byte {
	startIndex := bytes.IndexByte(jsonpResponse, '(')
	endIndex := bytes.LastIndexByte(jsonpResponse, ')')
	if startIndex == -1 || endIndex == -1 {
		return nil
	}
	return jsonpResponse[startIndex+1 : endIndex]
}

func getContentType(h http.Header) string {
	s := h.Get("Content-Type")
	return strings.Split(s, ";")[0]
}
