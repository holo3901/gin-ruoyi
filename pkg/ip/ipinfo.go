package ipinfo

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
)

func GetRemoteClientIp(r *http.Request) string {
	remoteIp := r.RemoteAddr

	if ip := r.Header.Get("x-Real-IP"); ip != "" {
		remoteIp = ip
	} else if ip = r.Header.Get("X-Forwarded-For"); ip != "" {
		remoteIp = ip
	} else {
		remoteIp, _, _ = net.SplitHostPort(remoteIp)
	}

	//本地ip
	if remoteIp == "::1" {
		remoteIp = "127.0.0.1"
	}
	return remoteIp
}

type Tunit struct {
	Pro  string `json:"pro"`
	City string `json:"city"`
}

func GetRealAddressByIP(ip string) string {
	url := "http://whois.pconline.com.cn/ipJson.jsp?ip=" + ip + "&json=true"
	resp, err := http.Get(url)
	var result = "内网ip"
	if err != nil {
		result = "内网ip"
	} else {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			result = "内网ip"
		} else {
			dws := new(Tunit)
			json.Unmarshal(body, &dws)
			result = dws.Pro + " " + dws.City
		}
	}
	return result
}
