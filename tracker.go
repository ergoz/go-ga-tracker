package ga

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const HTTP_TIMEOUT = 2 * time.Second
const HTTP_MAX_IDLE_CONNS_PER_HOST = 1024
const GA_API_URL = "http://www.google-analytics.com/collect"
const GA_TRACKING_ID = "XX-XXXXXXX-X"

type (
	tracker struct {
		t *http.Transport
	}
)

var Tracker *tracker

func init() {
	transport := &http.Transport{}
	dial := &net.Dialer{Timeout: time.Duration(HTTP_TIMEOUT)}
	transport.Dial = dial.Dial
	transport.MaxIdleConnsPerHost = HTTP_MAX_IDLE_CONNS_PER_HOST
	Tracker = &tracker{t: transport}
}

func (s *tracker) Event(uid string, eventCategory string, eventAction string) {
	data := url.Values{
		"v":   {"1"},
		"tid": {GA_TRACKING_ID},
		"cid": {uid},
		"t":   {"event"},
		"ni":  {"1"},
		"ec":  {eventCategory},
		"ea":  {eventAction}}

	client := &http.Client{
		Transport: s.t,
	}

	req, err := http.NewRequest("POST", GA_API_URL, strings.NewReader(data.Encode()))
	if err != nil {
		fmt.Println("ga.TrackEvent err:", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	timer := time.AfterFunc(time.Duration(HTTP_TIMEOUT), func() {
		s.t.CancelRequest(req)
	})
	resp, err := client.Do(req)
	timer.Stop()
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Println("ga.TrackEvent err: Response status != 200")
	}
	ioutil.ReadAll(resp.Body)
}
