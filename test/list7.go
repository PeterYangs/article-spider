package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

func main() {
	proxyAddr := "socks5://127.0.0.1:4781"
	//proxyAddr := "http://127.0.0.1:4780"
	url := "https://store.shopping.yahoo.co.jp"
	cli := NewHttpClient(proxyAddr)
	data, _ := HttpGET(cli, url)
	fmt.Println(string(data))
}

func NewHttpClient(proxyAddr string) *http.Client {

	//proxy := func(_ *http.Request) (*url.URL, error) {
	//	return url.Parse(proxyAddr)
	//}

	//dialer, err := proxy.SOCKS5("tcp", proxyAddr, nil, proxy.Direct)
	//if err != nil {
	//	fmt.Fprintln(os.Stderr, "can't connect to the proxy:", err)
	//	os.Exit(1)
	//}

	netTransport := &http.Transport{
		Proxy: func(r *http.Request) (*url.URL, error) {
			return url.Parse(proxyAddr)
		},
		//MaxIdleConnsPerHost:   10,                             //每个host最大空闲连接
		//ResponseHeaderTimeout: time.Second * time.Duration(5), //数据收发5秒超时
	}

	return &http.Client{
		Timeout:   time.Second * 10,
		Transport: netTransport,
	}
}

func HttpGET(client *http.Client, url string) (body []byte, err error) {
	rsp, err := client.Get(url)
	if err != nil {
		return
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK || err != nil {
		err = fmt.Errorf("HTTP GET Code=%v, URI=%v, err=%v", rsp.StatusCode, url, err)
		return
	}

	return ioutil.ReadAll(rsp.Body)
}
