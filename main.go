package main

import (
	// "io"
	"fmt"
	"github.com/NovelDownloader/downloader"
	// "strings"
	"os"
	// "bufio"
	"github.com/NovelDownloader/config"
	"encoding/json"
	"io/ioutil"
)

func loadConifg(path string)(configs []config.ConfigType) {
	// file, err := os.Open(path)
	// if err == nil {
	// 	defer file.Close()
	// 	br := bufio.NewReader(file)
	// 	for {
	// 		a, _, c := br.ReadLine()
	// 		if c == io.EOF {
	// 			break
	// 		}
	// 		str := string(a)
	// 		strs :=strings.Split(str, ",")
	// 		configs = append(configs, config.ConfigType{
	// 			Filename: strs[0],
	// 			Url: strs[1],
	// 			Lasturl: strs[2],
	// 			Lastname: strs[3],
	// 		})
	// 	}
	// }
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("打开文件错误:", err.Error())
	}
	json.Unmarshal(buf, &configs)
	return
}

func updateConfig(path string, configs []config.ConfigType) {
	// file, err := os.OpenFile(path, os.O_WRONLY, 0644)
	// if err == nil {
	// 	defer file.Close()
	// 	var str string
	// 	for i := 0; i < len(configs); i++ {
	// 		conf := configs[i]
	// 		str += conf.Filename + "," + conf.Url + "," + conf.Lasturl + "," + conf.Lastname
	// 		if i < len(configs) - 1 {
	// 			str += "\n"
	// 		}
	// 	}
	// 	file.WriteString(str)
	// } else {
	// 	println("更新配置文件出错，", err.Error())
	// }
	buf, err := json.MarshalIndent(configs, "", "\t")
	if err != nil {
		fmt.Println("更新配置文件出错:", err.Error())
		return
	}
	jsonStr := string(buf)
	file, fErr := os.OpenFile(path, os.O_WRONLY, 0222)
	if fErr != nil {
		fmt.Println("写入出错:", fErr.Error())
	}
	defer file.Close()
	file.WriteString(jsonStr)
}

func main() {
	configs := loadConifg("config.ini")
	dl := downloader.NewDownloader()
	needUpdateConfig := false
	for i := 0; i < len(configs); i++ {
		conf := &configs[i]
		if conf.End {
			continue
		}
		fmt.Println("开始下载:", conf.Filename, "上次最新章节:", conf.Lastname, " ......")
		lastUrl, lastName := dl.StartWithConfig(conf)
		if lastUrl != "" {
			conf.Lasturl = lastUrl
			conf.Lastname = lastName
			needUpdateConfig = true
		}
	}
	if needUpdateConfig {
		updateConfig("config.ini", configs)
	}
	fmt.Println("下载完成")
}