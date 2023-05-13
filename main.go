package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gookit/color"
)

type Proxy struct {
	Data []struct {
		ID   string `json:"_id"`
		IP   string `json:"ip"`
		Port string `json:"port"`
	} `json:"data"`
	Total int `json:"total"`
	Page  int `json:"page"`
	Limit int `json:"limit"`
}
type CLIENT struct {
	client  *http.Client
	address string
}

const (
	url1, url2 = "http://example.com/", "https://example.com/"
)

var (
	lastCheck  string
	output     string
	outputTemp string
	textChan   = make(chan string)
	count      int
	checked    int
)

func println(str string, color func(a ...interface{}) string) {
	fmt.Printf("[%s] %s\n", color(now()), str)
}
func print(str string, color func(a ...interface{}) string) {
	fmt.Printf("[%s] %s", color(now()), str)
}
func printr(str string, i int, color func(a ...interface{}) string) {
	fmt.Printf("\r%s\r[%s] %s", strings.Repeat(" ", i), color(now()), str)
}
func MustReadBody(resp *http.Response, err error) []byte {
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return body
}

func textAppender() {
	//f, err := os.OpenFile("output.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	f, err := os.Create("output.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer f.Close()
	for text := range textChan {
		_, err := f.WriteString(text + "\n")
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
		outputTemp += text + "\r\n"
	}
}
func init() {
	go func() {

		for _, l := range ReadAllLines("output.txt") {
			output += l + "\r\n"
		}

		http.HandleFunc("/get-proxies", p)

		log.Fatal(http.ListenAndServe(":707", nil))
	}()
}

func p(w http.ResponseWriter, r *http.Request) {
	result := ""
	if r.URL.Query()["withdate"] != nil && r.URL.Query()["withdate"][0] != "" {
		result += lastCheck + "\r\n"
	}
	result += output
	fmt.Fprintf(w, result)
}

func main() {

	print("initializing benchmark checkers...", Cyan)
	expectedSum1, expectedSum2 := hashIT(MustReadBody(http.Get(url1))), hashIT(MustReadBody(http.Get(url2)))
	printr(fmt.Sprintf("Initialized, expectedSum1: %s, expectedSum2: %s\n", expectedSum1, expectedSum2), 55, Cyan)

	for {

		textChan = make(chan string)
		outputTemp = ""
		println("Fetching...", Cyan)

		ar := fetchAll()
		a := len(ar)
		removeDuplicates(&ar)
		//f, _ := os.Create("o.txt")
		//for _, a := range ar {
		//	f.WriteString(a + "\r\n")
		//}
		//f.Close()
		shuffle(&ar)
		count = len(ar)
		checked = 0
		CLIENTs := readProxies(ar)
		println(fmt.Sprintf("%d Proxy has been Scrapped, %d is duplicated", len(CLIENTs), a-len(CLIENTs)), Cyan)
		CLIENTsChan := make(chan CLIENT)
		go textAppender()
		threadsCount := Config.ThreadsCount
		println(fmt.Sprintf("Threads Count: %d", threadsCount), Cyan)
		var wg = new(sync.WaitGroup)
		for i := 0; i < threadsCount; i++ {
			wg.Add(1)
			go worker(expectedSum1, expectedSum2, CLIENTsChan, wg)
		}
		for _, c := range CLIENTs {
			CLIENTsChan <- c
		}
		close(CLIENTsChan)
		wg.Wait()
		close(textChan)
		output = outputTemp
		lastCheck = now()
		fmt.Println("\rDone.            ")
	}
}

func worker(expectedSum1, expectedSum2 string, c chan CLIENT, wg *sync.WaitGroup) {
	for c := range c {
		if checkProxy(expectedSum1, expectedSum2, c.client) {
			textChan <- c.address
			checked++
			fmt.Printf("\r[%d/%d]", checked, count)
		} else {
			checked++
			//fmt.Println("Doesnt work", c.address)
		}

	}
	wg.Done()
}

func readProxies(ar []string) []CLIENT {
	CLIENTs := []CLIENT{}
	for _, p := range ar {
		u, err := url.Parse(p)
		if err != nil {
			continue
		}
		tr := &http.Transport{
			Proxy: http.ProxyURL(u),
			DialContext: (&net.Dialer{
				Timeout:   2 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		}

		client := &http.Client{
			Transport: tr,
			Timeout:   5 * time.Second,
		}
		CLIENTs = append(CLIENTs, CLIENT{client, p})
	}
	return CLIENTs
}

func checkProxy(expectedSum1, expectedSum2 string, proxy *http.Client) bool {
	result1, err := md5sum(url1, proxy)
	if err != nil || result1 != expectedSum1 {
		return false
	}

	result2, err := md5sum(url2, proxy)
	if err != nil || result2 != expectedSum2 {
		return false
	}

	return true
}

func fetchAll() []string {
	var (
		ProxiesStrs     = []string{}
		ProxiesStrsChan = make(chan string)
	)
	go func() {
		for val := range ProxiesStrsChan {
			if _, err := url.Parse(val); err == nil {
				ProxiesStrs = append(ProxiesStrs, val)
			}
		}
	}()
	var wg sync.WaitGroup
	wg.Add(1)
	var urls = Config.Urls

	urlChan := make(chan URLDATA)
	go func() {
		for _, u := range urls {
			urlChan <- u
		}
		close(urlChan)
	}()
	urllerThreads := 15
	for i := 0; i < urllerThreads; i++ {
		go func() {
			wg.Add(1)
			for u := range urlChan {
				a := 0
				if isRaw(u) {
					for _, p := range fetchRaw(u.URL) {
						a++
						ProxiesStrsChan <- p
					}
				} else if isJson(u) {
					for _, p := range fetchJson(u) {
						a++
						ProxiesStrsChan <- p
					}
				} else if isRegex(u) && isIPPortRegex(u) {

				} else if isRegex(u) && isIPRegexPortRegex(u) {

				} else if isAutoRegex(u) {
					for _, p := range fetchAutoRegex(u.URL) {
						a++
						ProxiesStrsChan <- p
					}
				} else {
					continue
				}

				if a > 0 {
					println(fmt.Sprintf("%s Feched, %d Proxy", u.URL, a), Green)
				} else {
					println(fmt.Sprintf("%s is Empty", u.URL), Red)
				}
			}
			wg.Add(-1)
		}()
	}
	time.Sleep(time.Second)
	wg.Add(-1)

	wg.Wait()
	close(ProxiesStrsChan)
	return ProxiesStrs
}

func now() string {
	loc, _ := time.LoadLocation("Asia/Riyadh")
	return time.Now().In(loc).Format("2006-01-02 15:04:05")
}

var (
	Black          = color.Black.Render
	Red            = color.Red.Render
	Green          = color.Green.Render
	Yellow         = color.Yellow.Render
	Blue           = color.Blue.Render
	Magenta        = color.Magenta.Render
	Cyan           = color.Cyan.Render
	White          = color.White.Render
	Gray           = color.Gray.Render
	LightRed       = color.LightRed.Render
	LightGreen     = color.LightGreen.Render
	LightYellow    = color.LightYellow.Render
	LightBlue      = color.LightBlue.Render
	LightMagenta   = color.LightMagenta.Render
	LightCyan      = color.LightCyan.Render
	LightWhite     = color.LightWhite.Render
	BgBlack        = color.BgBlack.Render
	BgRed          = color.BgRed.Render
	BgGreen        = color.BgGreen.Render
	BgYellow       = color.BgYellow.Render
	BgBlue         = color.BgBlue.Render
	BgMagenta      = color.BgMagenta.Render
	BgCyan         = color.BgCyan.Render
	BgWhite        = color.BgWhite.Render
	BgGray         = color.BgGray.Render
	BgDarkGray     = color.BgDarkGray.Render
	BgLightRed     = color.BgLightRed.Render
	BgLightGreen   = color.BgLightGreen.Render
	BgLightYellow  = color.BgLightYellow.Render
	BgLightBlue    = color.BgLightBlue.Render
	BgLightMagenta = color.BgLightMagenta.Render
	BgLightCyan    = color.BgLightCyan.Render
	BgLightWhite   = color.BgLightWhite.Render
	Bold           = color.Bold.Render
)
