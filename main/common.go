package main

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	uuid "github.com/odeke-em/go-uuid"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func GeneSignID(stime time.Time, uuidStr string) string {
	nano := stime.UnixNano()
	nanostr := strconv.FormatInt(nano, 10)

	if len(uuidStr) > 32 {
		uuidStr = uuidStr[32:]
	}

	didby, _ := hex.DecodeString(uuidStr)
	didsum := GetSum(didby)
	rand.Seed(nano)

	productid := fmt.Sprintf("%s%s%v%v%v", stime.Format("20060102"), nanostr[len(nanostr)-4:], rand.Intn(9), rand.Intn(9), didsum)
	Debug(productid)
	return productid
	// return fmt.Sprintf("%v", stime.UnixNano())
}

func GetSum(buf []byte) uint16 {
	var a uint16 = 0xbeaf
	for _, v := range buf {
		a += uint16(v)
	}
	return a
}

func GetReplaceUUID() string {
	uuidString := uuid.UUID4().String()
	return strings.Replace(uuidString, "-", "", -1)
}

func Base64Encode(src []byte) string {
	return base64.StdEncoding.EncodeToString(src)
}

func Base64Decode(src string) []byte {
	dst, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return nil
	}

	return dst
}

func GetTimeNow() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func ParseRequest(r *http.Request, data interface{}) (ret int) {
	reqBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		Error("ReadAll fail: ", err)
		ret = ErrInnerFault
		return
	}

	Debug(string(reqBytes))

	if err := json.Unmarshal(reqBytes, &data); nil != err {
		Error(err)
		ret = ErrInnerFault
		return
	}

	return Success
}

func StructToJsonString(v interface{}) (string, int) {
	jsonBody, err := json.Marshal(v)
	if err != nil {
		Error("json Marshal error:", err)
		return "", ErrInnerFault
	}
	return string(jsonBody), Success
}

func dialTimeout(network, addr string) (net.Conn, error) {
	c, err := net.DialTimeout(network, addr, time.Second*50) //设置建立连接超时
	if err != nil {
		return nil, err
	}
	err = c.SetDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		return nil, err
	} //设置发送接收数据超时
	return c, nil
}

func GetHttpClient() (client *http.Client) {
	transport := http.Transport{
		Dial:              dialTimeout,
		DisableKeepAlives: true,
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
	}
	client = &http.Client{
		Transport: &transport,
	}
	return
}

func HttpRequest(method, url string, body []byte, header map[string]string) (resp *http.Response, err error) {
	newReq, _ := http.NewRequest(method, url, bytes.NewReader(body))

	for k, v := range header {
		newReq.Header.Set(k, v)
	}

	client := GetHttpClient()
	resp, err = client.Do(newReq)
	if err != nil {
		Error("new http req err.", err.Error())
		return
	}
	return
}

func ParseResponseString(response *http.Response) (string, error) {
	//var result map[string]interface{}
	body, err := ioutil.ReadAll(response.Body) // response.Body 是一个数据流
	return string(body), err                   // 将 io数据流转换为string类型返回！
}

func signalHandle() {
	t1 := time.NewTimer(1 * time.Hour)
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGUSR1, syscall.SIGUSR2)
	for {
		select {
		case sig := <-ch:
			switch sig {
			case syscall.SIGUSR1:
				setDebugLevel("DEBUG")
				t1.Reset(1 * time.Hour)
			case syscall.SIGUSR2:
				setDebugLevel("ERROR")
				t1.Reset(1 * time.Hour)
				t1.Stop()
			default:
			}
		case <-t1.C:
			setDebugLevel("ERROR")
			t1.Stop()
		}
	}
}
