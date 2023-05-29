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

func fofac(fofa_email string, fofa_api string, fofa_dock string) map[string]int {
	result := make(map[string]int)
	dorkbase64 := base64.StdEncoding.EncodeToString([]byte(fofa_dock))

	url1 := "https://fofa.info/api/v1/search/all?size=10000&email=" + url.QueryEscape(fofa_email) + "&key=" + url.QueryEscape(fofa_api) + "&qbase64=" + url.QueryEscape(dorkbase64)
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
		if in(ip, list) == false {
			list = append(list, ip)
		}
	}

	for i := range list {
		if strings.Contains(list[i], ":") {
			continue
		}
		ipc := strings.Split(list[i], ".")[0:3]
		a := ipc[0] + "." + ipc[1] + "." + ipc[2] + ".0/24"
		_, status := result[a]
		if status == true {

			result[a] = result[a] + 1
		} else {
			result[a] = 1
		}
	}

	return result
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
