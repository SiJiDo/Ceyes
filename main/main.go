package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

type Yaml struct {
	Fofa_email string `yaml:"fofa_email"`
	Fofa_api   string `yaml:"Fofa_api"`
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func GetFofaAuth() (string, string) {
	conf := new(Yaml)
	yamlFile, err := ioutil.ReadFile("config.yaml")

	err = yaml.Unmarshal(yamlFile, conf)

	checkError(err)
	fofa_email := conf.Fofa_email
	fofa_api := conf.Fofa_api
	return fofa_email, fofa_api
}

func setFofaAuthfile(src string) {
	stu := &Yaml{
		Fofa_email: "",
		Fofa_api:   "",
	}
	data, err := yaml.Marshal(stu)
	checkError(err)
	err = ioutil.WriteFile(src, data, 0777)
	checkError(err)
}

func main() {

	var filename string
	var fofaSearch string
	var sortCount bool
	var cloud bool
	var fofaDomain string
	var fofaDork []string
	fofa_email, fofa_api := GetFofaAuth()

	flag.StringVar(&fofaDomain, "d", "", "domain deafult use dork like (domain= xxx || host= xxx)")
	flag.StringVar(&filename, "f", "", "domain text")
	flag.StringVar(&fofaSearch, "s", "", "fofa search dork")
	flag.BoolVar(&sortCount, "sc", false, "sort result by count, deafult ip sort")
	flag.BoolVar(&cloud, "cloud", false, "check cloud")

	flag.Parse()

	f, err := os.Open("config.yaml")
	if err != nil && os.IsNotExist(err) {
		setFofaAuthfile("config.yaml")
	}
	f.Close()

	if filename == "" && fofaDomain == "" {
		fofaDork = append(fofaDork, fofaSearch)

	} else if fofaDomain != "" {
		tmp := "domain=\"" + fofaDomain + "\" || host=\"" + fofaDomain + "\""
		fofaDork = append(fofaDork, tmp)
	} else { //read domain info to make fofa dork and search
		file, err := os.Open(filename)

		if err != nil {
			fmt.Println("file error:", err)
		}
		reader := bufio.NewReader(file)
		for {
			domain, err := reader.ReadString('\n') // end flag
			domain = strings.Replace(domain, " ", "", -1)
			domain = strings.Replace(domain, "\n", "", -1)
			domain = strings.Replace(domain, "\r", "", -1)
			if domain != "" {
				tmp := "domain=\"" + domain + "\" || host=\"" + domain + "\""
				fofaDork = append(fofaDork, tmp)
			}
			if err == io.EOF {
				break
			}
		}
	}

	//strat search
	final_result := make(map[string]int)
	final_result_domain := make(map[string]string)
	final_cloud_result := make(map[string]string)
	for i := range fofaDork {
		fmt.Println("[+]now fofa dork is: [ " + fofaDork[i] + " ]")
		result, cloud_result := fofac(fofa_email, fofa_api, fofaDork[i], cloud)
		//sort result
		r := make([]ip_count, 0)
		if sortCount == false {
			r = sortbycount(result)
		} else {
			r = sortbyip(result)
		}
		for _, pair := range r {
			if len(pair.Key) == 15 {
				fmt.Printf("[+]ipc:%v   count: %v	%v \n", pair.Key, pair.Val, cloud_result[pair.Key])
			} else if len(pair.Key) == 16 {
				fmt.Printf("[+]ipc:%v  count: %v	%v \n", pair.Key, pair.Val, cloud_result[pair.Key])
			} else {
				fmt.Printf("[+]ipc:%v\t count: %v	%v \n", pair.Key, pair.Val, cloud_result[pair.Key])
			}
		}
		if filename == "" {
			break
		} else {
			for k, v := range result {
				_, status := final_result[k]
				if cloud_result[k] != "" {
					final_cloud_result[k] = "[" + cloud_result[k] + "]\t"
					if strings.Contains(final_cloud_result[k], "(maybe)") == false {
						final_cloud_result[k] = final_cloud_result[k] + "\t"
					}
				} else {
					final_cloud_result[k] = "\t\t\t"
				}
				if status == true {
					final_result[k] = final_result[k] + result[k]
					final_result_domain[k] = final_result_domain[k] + ", " + strings.Split(fofaDork[i], "host=")[1]
				} else {
					final_result[k] = v
					final_result_domain[k] = strings.Split(fofaDork[i], "host=")[1]

				}
			}
		}
	}

	if filename != "" {
		fmt.Println("====================all domains in file result=========================")

		r := make([]ip_count, 0)
		if sortCount == false {
			r = sortbycount(final_result)
		} else {
			r = sortbyip(final_result)
		}
		for _, pair := range r {
			domain := final_result_domain[pair.Key]
			if len(pair.Key) == 15 {
				fmt.Printf("[+]ipc:%v   count: %v \t%v domain: %v\n", pair.Key, pair.Val, final_cloud_result[pair.Key], domain)
			} else if len(pair.Key) == 16 {
				fmt.Printf("[+]ipc:%v  count: %v \t%v domain: %v\n", pair.Key, pair.Val, final_cloud_result[pair.Key], domain)
			} else {
				fmt.Printf("[+]ipc:%v\t count: %v \t%v domain: %v\n", pair.Key, pair.Val, final_cloud_result[pair.Key], domain)
			}
		}

	}
}
