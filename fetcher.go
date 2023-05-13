package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"old/github.com/tidwall/gjson"
	"strings"
)

func fetchRaw(u string) (result []string) {
	client := &http.Client{
		Transport: &http.Transport{},
	}
	if (!strings.Contains(u, "github.") && Config.Proxy.Useproxy) || Config.Proxy.ForceUse {
		proxy, _ := url.Parse(Config.Proxy.URL)
		client = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxy),
			},
		}
	}

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	response := string(body)
	var strs []string
	if len(response) > 5 {
		strs = strings.Split(response, "\n")
		for n, s := range strs {
			strs[n] = strings.ReplaceAll(s, "\r", "")
		}
	}
	for _, proxy := range strs {
		result = append(result, fmt.Sprintf("http://%s", proxy))
	}
	return result
}

func fetchJson(d URLDATA) (result []string) {
	client := &http.Client{
		Transport: &http.Transport{},
	}
	if (!strings.Contains(d.URL, "github.") && Config.Proxy.Useproxy) || Config.Proxy.ForceUse {
		proxy, _ := url.Parse(Config.Proxy.URL)
		client = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxy),
			},
		}
	}
	req, err := http.NewRequest("GET", d.URL, nil)
	if err != nil {
		return
	}
	req.Header.Set(`sec-ch-ua`, `"Google Chrome";v="111", "Not(A:Brand";v="8", "Chromium";v="111"`)
	req.Header.Set(`sec-ch-ua-mobile`, `?0`)
	req.Header.Set(`sec-ch-ua-platform`, `"Windows"`)
	req.Header.Set(`Upgrade-Insecure-Requests`, `1`)
	req.Header.Set(`accept`, `text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7`)
	req.Header.Set(`accept-language`, `en-US,en;q=0.9,ar;q=0.8`)
	req.Header.Set(`cache-control`, `max-age=0`)
	req.Header.Set(`sec-ch-ua`, `"Google Chrome";v="111", "Not(A:Brand";v="8", "Chromium";v="111"`)
	req.Header.Set(`sec-ch-ua-mobile`, `?0`)
	req.Header.Set(`sec-ch-ua-platform`, `"Windows"`)
	req.Header.Set(`sec-fetch-dest`, `document`)
	req.Header.Set(`sec-fetch-mode`, `navigate`)
	req.Header.Set(`sec-fetch-site`, `none`)
	req.Header.Set(`sec-fetch-user`, `?1`)
	req.Header.Set(`upgrade-insecure-requests`, `1`)
	req.Header.Set(`user-agent`, `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36`)
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	ar := []gjson.Result{}
	if d.ToArray != "" {
		ar = gjson.Get(string(body), d.ToArray).Array()
		for _, a := range ar {
			if d.IPPortJSON != "" {
				result = append(result, fmt.Sprintf("http://%s", gjson.Get(a.String(), d.IPPortJSON)))
			} else {
				result = append(result, fmt.Sprintf("http://%s:%s", gjson.Get(a.String(), d.IPJSON), gjson.Get(a.String(), d.PortJSON)))
			}
		}
	}
	return result
}

func fetchAutoRegex(u string) (result []string) {
	client := &http.Client{
		Transport: &http.Transport{},
	}
	if (!strings.Contains(u, "github.") && Config.Proxy.Useproxy) || Config.Proxy.ForceUse {
		proxy, _ := url.Parse(Config.Proxy.URL)
		client = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxy),
			},
		}
	}

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	response := string(body)
	list := autoRegex(response)
	for _, l := range list {
		if isValidIPPort(l) {
			result = append(result, fmt.Sprintf("http://%s", l))
		}
	}
	return result
}
func fetchIPPortRegex(d URLDATA) (result []string) {
	client := &http.Client{
		Transport: &http.Transport{},
	}
	if (!strings.Contains(d.URL, "github.") && Config.Proxy.Useproxy) || Config.Proxy.ForceUse {
		proxy, _ := url.Parse(Config.Proxy.URL)
		client = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxy),
			},
		}
	}

	req, err := http.NewRequest("GET", d.URL, nil)
	if err != nil {
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	response := string(body)
	list := regexIPPort(response, d.IPPortRg)
	for _, l := range list {
		if isValidIPPort(l) {
			result = append(result, fmt.Sprintf("http://%s", l))
		}
	}
	return result
}
