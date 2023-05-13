package main

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func removeDuplicates(strings *[]string) {
	seen := make(map[string]bool)
	i := 0
	for _, s := range *strings {
		if !seen[s] {
			seen[s] = true
			(*strings)[i] = s
			i++
		}
	}
	*strings = (*strings)[:i]
}

func shuffle(s *[]string) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range *s {
		j := r.Intn(i + 1)
		(*s)[i], (*s)[j] = (*s)[j], (*s)[i]
	}
}
func ReadAllLines(path string) (lines []string) {
	file, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

func appendTheArray(array, newArray []string) []string {
	for _, line := range newArray {
		if _, err := url.Parse(line); err == nil {
			array = append(array, line)
		}
	}
	return array
}

/*
ToArray    string `yaml:"to_array,omitempty"`
IPJSON     string `yaml:"ip_json,omitempty"`
PortJSON   string `yaml:"port_json,omitempty"`
IPPortJSON string `yaml:"ip_port_json,omitempty"`
IPRg       string `yaml:"ip_rg,omitempty"`
PortRg     string `yaml:"port_rg,omitempty"`
IPPortRg   string `yaml:"ip_port_rg,omitempty"`
*/
func isRaw(d URLDATA) bool {
	return d.ToArray == "" && d.IPJSON == "" && d.PortJSON == "" && d.IPPortRg == "" && d.IPPortJSON == "" && !d.AutoRegex // && d.IPRg == "" && d.PortRg == ""
}

func isJson(d URLDATA) bool {
	return d.ToArray != "" && (d.IPPortJSON != "" || (d.IPJSON != "" && d.PortJSON != ""))
}

func isRegex(d URLDATA) bool {
	return d.IPPortRg != "" // || (d.IPRg != "" && d.PortRg != "")
}

func isIPPortRegex(d URLDATA) bool {
	return d.IPPortRg != "" //&& (d.IPRg == "" || d.PortRg == "")
}

func isIPRegexPortRegex(d URLDATA) bool {
	return d.IPPortRg == "" //&& (d.IPRg != "" && d.PortRg != "")
}

func isAutoRegex(d URLDATA) bool {
	return d.AutoRegex
}
func regexIPPort(resp, str string) []string {
	re := regexp.MustCompile(str)

	matches := re.FindAllStringSubmatch(resp, -1)

	result := []string{}
	for _, match := range matches {
		result = append(result, match[1])
	}
	return result
}
func autoRegex(str string) []string {
	re := regexp.MustCompile(`\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}:\d{1,5}`)
	ipPortList := re.FindAllString(str, -1)
	return ipPortList
}

func isValidIPPort(input string) bool {
	parts := strings.Split(input, ":")
	if len(parts) != 2 {
		return false
	}

	ip := net.ParseIP(parts[0])
	if ip == nil {
		return false
	}

	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return false
	}
	if port < 1 || port > 65535 {
		return false
	}

	_, err = net.ResolveTCPAddr("tcp", net.JoinHostPort(ip.String(), strconv.Itoa(port)))
	if err != nil {
		return false
	}

	if ip.IsLoopback() || ip.IsMulticast() {
		return false
	}

	if ip.IsPrivate() {
		return false
	}

	if ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
		return false
	}

	return true
}

func hashIT(s []byte) string {
	return fmt.Sprintf("%x", md5.Sum(s))
}

func md5sum(Url string, proxy *http.Client) (string, error) {
	resp, err := proxy.Get(Url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return hashIT(data), nil
}
