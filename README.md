# Proxy-Wizard
 Proxy-Wizard is a Go script that fetches and tests proxies, and provides an API to access the working proxies.
 Free proxy list in output.txt updated every minute.
## Installation
 - Ensure that Go is installed on your system. If it isn't, you can download it from [here](https://golang.org/dl/).
 - Clone this repository: $ git clone `https://github.com/d7emy/Proxy-Wizard.git`
 - Navigate to the project directory: `$ cd Proxy-Wizard`
## Usage
 - Open the `config.yaml` file to specify your desired proxy sources. You can add or remove sources as needed.
 - Start the program: `go run .`
 - The program will begin fetching and testing proxies. Working proxies will be saved to the `output.txt` file.
 - To retrieve the Working proxies, send a GET request to http://localhost:707/get-proxies. The proxies will be returned in raw.

## Configuration
- ### Threads
    ```yaml
    threadsCount: 125
    urllerThreadsCount: 15
    ```
     -  `threadsCount` specifies the number of threads to be used for proxies checker.
     -  `urllerThreadsCount ` specifies the number of threads to be used for fetching proxies.
- ### Proxy Configuration
    ```yaml
    proxy:
      useproxy: false
      forceuse: false
      url: http://username:password@ip:port #or http://ip:port
    ```
    - `url` this should be set to the URL of the proxy server to be used
    - `useproxy` if it's true it will use the proxy when fetching from blocked proxies websites.
    - `forceuse` if it's true it will use the proxy when fetching from github urls of another non proxy websites.
- ### URLs Configuration
    ```yaml
    urls:
      - url: https://proxylist.geonode.com/api/proxy-list?limit=500&page=1&sort_by=lastChecked&sort_type=desc&protocols=http,https
        to_array: data
        ip_json: ip
        port_json: port
      - url: https://openproxy.space/list/http
        autoRegex: true
      - url: https://api.proxyscrape.com/v2/?request=getproxies&protocol=http&timeout=10000&country=all&ssl=all&anonymity=all
    ```

    - `url` this should be the url for the proxies. fetching keys (leave the key empty if the response in raw):
        - `to_array` is the path for the json array in the response body, `ip_json` the key for the ip, `port_json` the key for the port, or u can use `ip_port_json` if the ip:port in the key's value.
        - `autoRegex` if it's true it will Automatically regex all the ip:port in the response body.
        - `ip_port_rg` this should be the regex string to get the proxies from for example html response.

## Notice
 - This Repo is for educational purposes. i'm not response for any illegal usage.
## License

[![MIT License](https://img.shields.io/badge/License-MIT-yellow.svg)](https://choosealicense.com/licenses/mit/)

**Free HQ Software, Hell Yeah!**
