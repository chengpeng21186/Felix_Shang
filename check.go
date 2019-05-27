package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	KUBE = "KUBE"
	HTTP = "HTTP"
	TCP  = "TCP"
)

func main() {

	if enable, checkstr := CheckEnable(); enable {
		//fmt.Println(checkstr)
		fmt.Println("Get PRECHECK OK!", checkstr)
		ss := strings.Split(checkstr, "//")
		//fmt.Println(ss)
		if strings.EqualFold(ss[0], KUBE) {
			checkKube(ss[1])
			// EqualFold判断两个字符串在完全小写的情况下是否相等
		} else if strings.EqualFold(ss[0], TCP) {
			checkTcp(ss[1])
		} else if strings.EqualFold(ss[0], HTTP) {
			checkHttp(ss[1])
		} else {
			fmt.Println("check conditon is invalid")
		}
	}
}

//export PRECHECK=KUBE//xh:xhht
func checkKube(nssvc string) {
	nss := strings.Split(nssvc, ":")
	for i := 0; ; i++ {
		if sendKube(nss[0], nss[1]) {
			return
		} else {
			time.Sleep(5 * time.Second)
			if i > 60 {
				time.Sleep(time.Second * 60)
			}
		}
	}
}

func sendKube(namespace, deploy string) bool {
	//fmt.Printf(GetKuberneteAddr())
	url := "https://" + GetKuberneteAddr() + "/apis/apps/v1beta2/namespaces/" + namespace + "/deployments/" + deploy
	fmt.Printf(url)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Bearer dbeb0da3abc866256468947936af87b9")
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("error request kube %v \n", err)
		return false
	}
	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)

	res := make(map[string]interface{})
	err1 := json.Unmarshal(result, &res)
	if err != nil || err1 != nil {
		fmt.Printf("get kube body  %v \n", err)
		return false
	}

	k, _ := res["kind"]
	kind := k.(string)
	if strings.EqualFold(kind, "Status") {
		fmt.Printf("%q status err or not found", deploy)
		return false
	} else if strings.EqualFold(kind, "Deployment") {
		s, _ := res["status"]
		smap := s.(map[string]interface{})
		avalible, ok := smap["availableReplicas"].(float64)
		if ok && avalible > 0 {
			fmt.Printf("check %q success \n", deploy)
			return true
		} else {
			fmt.Printf("wait %q success ...\n", deploy)
			return false
		}
	}
	fmt.Println("unknow  status ...")
	return false

}

func GetKuberneteAddr() string {
	host := os.Getenv("KUBERNETES_SERVICE_HOST")
	port := os.Getenv("KUBERNETES_SERVICE_PORT")
	if len(host) > 0 && len(port) > 0 {
		return host + ":" + port
	}
	return "10.254.0.1:443"
}

func CheckEnable() (bool, string) {
	checkstr := os.Getenv("PRECHECK")
	if len(checkstr) > 0 {
		//fmt.Println(checkstr)
		return true, checkstr
	}
	fmt.Println("Not found PRECHECK")
	return false, ""
}

//PRECHECK=HTTP//liuduo.caiwu.corp@10
func checkHttp(nssvc string) {
	nss := strings.Split(nssvc, "@")
	for i := 0; ; i++ {
		if sendHttp(nss[0], nss[1]) {
			fmt.Println("check HTTP health!")
			return
		} else {
			time.Sleep(5 * time.Second)
			if i > 60 {
				time.Sleep(time.Second * 60)
			}
		}
	}
}

//PRECHECK=TCP//10.143.135.210:8080@10
func checkTcp(nssvc string) {
	nss := strings.Split(nssvc, "@")
	//fmt.Println(nss)
	for i := 0; ; i++ {
		//fmt.Println(time.Second) //1s
		if sendTcp(nss[0], nss[1]) {
			fmt.Println("check TCP health!")
			return
		} else {
			time.Sleep(5 * time.Second)
			if i > 60 {
				time.Sleep(time.Second * 60)
			}
		}
	}
}

func sendTcp(addr, t string) bool {
	//fmt.Println(addr)
	//fmt.Println(t)
	//string到int
	tt, err := strconv.Atoi(t)
	if err != nil {
		fmt.Println("err parse time")
		return false
	}
	//time.Duration（时间长度，消耗时间）
	//time.Time（时间点）
	//time.C（放时间的channel通道）（注：Time.C:=make(chan time.Time)）
	timeout := time.Duration(int64(tt) * int64(time.Second))
	//建立连接，设置超时时间
	resp, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		fmt.Printf("get addr %v err timeoout \n", addr)
		return false
	}
	//resp.SetDeadline(time.Now().Add(rwTimeout))
	//判断err，没有err则defer resp.Close()，不然不关闭.
	defer resp.Close()
	return true
}

func sendHttp(url, t string) bool {
	tt, err := strconv.Atoi(t)
	if err != nil {
		fmt.Println("err parse time")
		return false
	}
	timeout := time.Duration(int64(tt) * int64(time.Second))
	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Get("http://" + url)
	//fmt.Println(resp)
	if err != nil {
		fmt.Printf("get url %v err timeoout \n", url)
		return false
	}
	defer resp.Body.Close()
	fmt.Println(resp.StatusCode)
	if 200 > resp.StatusCode || resp.StatusCode > 399 {
		fmt.Printf("http url %v code not 200 ~ 400 \n", url)
		return false
	}
	return true
}
