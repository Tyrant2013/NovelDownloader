package main

import (
	"io"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"os"
	"bufio"
)

var (
	confFilePath = "config.ini"
	novelInfos = map[string]string{
		"tianxingzhilu" : "http://www.4xiaoshuo.com/52/52722/",
	}
)

type (
	ConfigType struct {
		filename string
		url string
		lastOne string
	}
	Chapter struct {
		name string
		url string
	}
)

func loadConifg()(config []ConfigType) {
	file, err := os.Open(confFilePath)
	if err == nil {
		defer file.Close()
		br := bufio.NewReader(file)
		for {
			a, _, c := br.ReadLine()
			if c == io.EOF {
				break
			}
			str := string(a)
			strs :=strings.Split(str, ",")
			config = append(config, ConfigType{
				strs[0],
				strs[1],
				strs[2],
			})
		}
	}
	return
}

func updateConfig(configs []ConfigType) {
	file, err := os.OpenFile(confFilePath, os.O_WRONLY, 0644)
	if err == nil {
		defer file.Close()
		var str string
		for i := 0; i < len(configs); i++ {
			conf := configs[i]
			if i < len(configs) - 1 {
				str += conf.filename + "," + conf.url + "," + conf.lastOne + "\n"
			} else {
				str += conf.filename + "," + conf.url + "," + conf.lastOne
			}
		}
		println("config.ini :", str)
		file.WriteString(str)
	} else {
		println("更新配置文件出错，", err.Error())
	}
}

func checkIsNewNovel(lastOne string) bool {
	return lastOne == ""
}

func findUrls(url string)(urls []Chapter) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		fmt.Println("Fatal error: ", err.Error())
	}
	doc.Find("div.listmain").Each(func(index int, s *goquery.Selection) {
		s.Find("a").Each(func(i int, s *goquery.Selection) {
			src, finded := s.Attr("href")
			name := s.Text()
			if finded {
				urls = append(urls, Chapter{
					name,
					src,
				})
			}
		})
	})
	return
}

func filterUrls(urls []Chapter, lastOne string) []Chapter {
	if lastOne == "" {
		if len(urls) > 12 {
			return urls[12:]
		} else {
			return urls[(len(urls) / 2) :]
		}
	}
	for i := 0; i < len(urls); i++ {
		if  urls[i].url == lastOne {
			return reverse(urls[:i])
		}
	}
	return urls
}

func reverse(urls []Chapter) []Chapter {
	for i, j := 0, len(urls) - 1; i < j; i, j = i + 1, j - 1 {
		urls[i], urls[j] = urls[j], urls[i]
	}
	return urls
}

func downloadWithUrls(conf ConfigType, urls []Chapter) (lastOne string, finished bool) {
	for _, chap := range urls {
		fmt.Printf("下载：%s ......", chap.name)
		content := downloadText(conf, chap.url)
		if content == "" {
			finished = false
			return
		}
		saveToFile(chap.name, content, conf.filename)
		fmt.Println(" 完成。")
		lastOne = chap.url
	}
	finished = true
	return
}

func downloadText(conf ConfigType, url string)(content string) {
	uri := conf.url + url
	doc, err := goquery.NewDocument(uri)
	if err == nil {
		// title = doc.Find("div.content>h1:first-child").Text()
		content = strings.Replace(doc.Find("#content").Text(), "请记住本书首发域名：wwww.4xiaoshuo.com。4小说网手机版阅读网址：m.4xiaoshuo.com", "", -1)
	} else {
		println("error:", err.Error())
	}
	return
}

func saveToFile(title string, content string, filePath string) {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil {
		defer file.Close()
		file.WriteString(title + "\n" + content)
	}
}

func main() {
	configs := loadConifg()
	for i := 0; i < len(configs); i++ {
		conf := &configs[i]
		urls := findUrls(conf.url)
		validedUrls := filterUrls(urls, conf.lastOne)
		lastOne, finished := downloadWithUrls(*conf, validedUrls)
		if !finished {
			fmt.Println("下载过程出错了, 最后下载的是:", lastOne)
		}
		conf.lastOne = lastOne
		updateConfig(configs)
	}
}