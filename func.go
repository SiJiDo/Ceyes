package main

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"os"
	"sort"
	"strings"

	"github.com/kirinlabs/HttpRequest"
	"github.com/tidwall/gjson"
)

func in(target string, str_array []string) bool {
	for _, element := range str_array {
		if target == element {
			return true
		}
	}
	return false
}

func check_cloud(cloud bool, info string) string {
	result := ""
	if cloud == false {
		return result
	}

	org_cloudlist := make(map[string]string)
	org_cloudlist["HuaWeiCloud"] = "Huawei Cloud Service data center"
	org_cloudlist["AliyunCloud"] = "Hangzhou Alibaba Advertising Co.,Ltd. & Alibaba US Technology Co., Ltd."
	org_cloudlist["TencentCloud"] = "Shenzhen Tencent Computer Systems Company Limited & Tencent Building, Kejizhongyi Avenue"
	org_cloudlist["BaiduCloud"] = "Beijing Baidu Netcom Science and Technology Co., Ltd."
	org_cloudlist["AmazonCloud"] = "AMAZON-02 & AMAZON-AES"
	org_cloudlist["GoogleCloud"] = "GOOGLE-IT & GOOGLE"
	org_cloudlist["AzureCloud"] = " MICROSOFT-CORP-MSN-AS-BLOCK"
	org_cloudlist["Cloudflare"] = "CLOUDFLARENET"

	cloudname_cloudlist := make(map[string]string)
	cloudname_cloudlist["HuaWeiCloud"] = "HuaWeiCloud"
	cloudname_cloudlist["AliyunCloud"] = "aliyun"
	cloudname_cloudlist["TencentCloud"] = "tencent"
	cloudname_cloudlist["BaiduCloud"] = "baidu"
	cloudname_cloudlist["AmazonCloud"] = "Amazon"
	cloudname_cloudlist["GoogleCloud"] = "google"
	cloudname_cloudlist["AzureCloud"] = "azure"
	cloudname_cloudlist["Cloudflare"] = "Cloudflare"

	tag := strings.Split(info, "+")
	org := tag[1]
	cloud_name := tag[2]
	for k, v := range org_cloudlist {
		if strings.Contains(v, org) {
			result = k + "(maybe)"
		}
	}
	for k, v := range cloudname_cloudlist {
		if cloud_name == v {
			result = k
		}
	}

	return result
}

func fofac(fofa_email string, fofa_api string, fofa_dock string, cloud bool) (map[string]int, map[string]string) {

	result := make(map[string]int)
	cloud_result := make(map[string]string)
	dorkbase64 := base64.StdEncoding.EncodeToString([]byte(fofa_dock))

	url1 := "https://fofa.info/api/v1/search/all?fields=host,ip,port,as_organization&size=10000&email=" + url.QueryEscape(fofa_email) + "&key=" + url.QueryEscape(fofa_api) + "&qbase64=" + url.QueryEscape(dorkbase64)
	req := HttpRequest.NewRequest()
	req.SetCookies(map[string]string{
		"auth_token": "123",
	})

	resp, err := req.Get(url1)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	body, err := resp.Body()
	jsondata := gjson.Get(string(body), "results")

	datalist := jsondata.Array()
	var list []string
	for i := range datalist {
		data := datalist[i]
		ip := data.Array()[1].String()
		org := data.Array()[3].String()
		cloud_name := data.Array()[4].String()
		if in(ip, list) == false {
			list = append(list, ip+"+"+org+"+"+cloud_name)
		}
	}

	var ip_list []string

	for i := range list {
		if strings.Contains(list[i], ":") {
			continue
		}
		ip := strings.Split(list[i], "+")[0]
		if in(ip, ip_list) {
			continue
		}
		ip_list = append(ip_list, ip)
		ipc := strings.Split(list[i], ".")[0:3]
		a := ipc[0] + "." + ipc[1] + "." + ipc[2] + ".0/24"
		_, status := result[a]
		if status == true {
			result[a] = result[a] + 1
			if strings.Contains(cloud_result[a], "(maybe)") {
				cloud_result[a] = check_cloud(cloud, list[i])
			}
		} else {
			result[a] = 1
			cloud_result[a] = check_cloud(cloud, list[i])
		}
	}

	return result, cloud_result
}

func Contains(s1, s2 string) {
	panic("unimplemented")
}

type ip_count struct {
	Key string
	Val int
}

func sortbyip(result map[string]int) []ip_count {
	newresult := make([]ip_count, 0)
	for k, v := range result {
		newresult = append(newresult, ip_count{Key: k, Val: v})
	}

	sort.Slice(newresult, func(i, j int) bool {
		return newresult[i].Val > newresult[j].Val
	})
	return newresult
}

func sortbycount(result map[string]int) []ip_count {
	newresult := make([]ip_count, 0)
	for k, v := range result {
		newresult = append(newresult, ip_count{Key: k, Val: v})
	}

	sort.Slice(newresult, func(i, j int) bool {
		return strings.Compare(newresult[i].Key, newresult[j].Key) < 0
	})
	return newresult
}
