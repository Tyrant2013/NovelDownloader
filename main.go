package main

import (
	"io"
	"fmt"
	"github.com/NovelDownloader/downloader"
	"strings"
	"os"
	"bufio"
	"github.com/NovelDownloader/config"
)

func loadConifg(path string)(configs []config.ConfigType) {
	file, err := os.Open(path)
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
			configs = append(configs, config.ConfigType{
				Filename: strs[0],
				Url: strs[1],
				Lasturl: strs[2],
				Lastname: strs[3],
			})
		}
	}
	return
}

func updateConfig(path string, configs []config.ConfigType) {
	file, err := os.OpenFile(path, os.O_WRONLY, 0644)
	if err == nil {
		defer file.Close()
		var str string
		for i := 0; i < len(configs); i++ {
			conf := configs[i]
			str += conf.Filename + "," + conf.Url + "," + conf.Lasturl + "," + conf.Lastname
			if i < len(configs) - 1 {
				str += "\n"
			}
		}
		file.WriteString(str)
	} else {
		println("更新配置文件出错，", err.Error())
	}
}

func main() {
	configs := loadConifg("config.ini")
	dl := downloader.NewDownloader()
	for i := 0; i < len(configs); i++ {
		conf := &configs[i]
		fmt.Println("开始下载:", conf.Filename, "上次最新章节:", conf.Lastname, " ......")
		lastUrl, lastName := dl.StartWithConfig(conf)
		if lastUrl != "" {
			conf.Lasturl = lastUrl
			conf.Lastname = lastName
			updateConfig("config.ini", configs)
		}
	}
	fmt.Println("下载完成")
}