package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

func main() {
	proxyAddr := "http://127.0.0.1:4780/"
	url := "https://www.google.com/"
	cli := NewHttpClient(proxyAddr)
	data, _ := HttpGET(cli, url)
	fmt.Println(string(data))
}

func NewHttpClient(proxyAddr string) *http.Client {
	//proxy, err := url.Parse(proxyAddr)
	//if err != nil {
	//	return nil
	//}

	netTransport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		//Proxy: http.ProxyURL(proxy),
		Dial: func(netw, addr string) (net.Conn, error) {
			c, err := net.DialTimeout(netw, addr, time.Second*time.Duration(10))
			if err != nil {
				return nil, err
			}
			return c, nil
		},
		MaxIdleConnsPerHost:   10,                             //每个host最大空闲连接
		ResponseHeaderTimeout: time.Second * time.Duration(5), //数据收发5秒超时
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
