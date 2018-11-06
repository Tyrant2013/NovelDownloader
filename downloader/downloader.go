package downloader

import (
	"os"
	"strings"
	"fmt"
	"github.com/NovelDownloader/config"
	"github.com/PuerkitoBio/goquery"
)

type (
	Downloader struct {

	}
	Chapter struct {
		name string
		url string
	}
)

func NewDownloader() (ins *Downloader) {
	ins = &Downloader{
		
	}
	return
}

func (_this Downloader)StartWithConfig(conf *config.ConfigType) (lastDownloadUrl, lastChapterName string) {
	chapters := findAllChaptersFromUrl(conf.Url, conf.ListReg)
	validChapters := filterValidChapters(chapters, conf.Lasturl)
	if len(validChapters) <= 0 {
		fmt.Println(conf.Filename + "没有更新")
	} else {
		lastDownloadUrl, lastChapterName = downloadAllValidChapters(validChapters, *conf)
	}
	return
}

func findAllChaptersFromUrl(url, listReg string)(chapters []Chapter) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		fmt.Println("Fatal error: ", err.Error())
	}
	// "div.listmain"
	doc.Find(listReg).Each(func(index int, s *goquery.Selection) {
		s.Find("a").Each(func(i int, s *goquery.Selection) {
			src, finded := s.Attr("href")
			name := s.Text()
			if finded {
				chapters = append(chapters, Chapter{
					name,
					src,
				})
			}
		})
	})
	return
}

func filterValidChapters(chapters []Chapter, lastOne string) []Chapter {
	if lastOne == "" {
		if len(chapters) > 12 {
			return chapters[12:]
		} else {
			return chapters[(len(chapters) / 2) :]
		}
	}
	for i := 0; i < len(chapters); i++ {
		if  chapters[i].url == lastOne {
			return reverse(chapters[:i])
		}
	}
	return chapters
}

func reverse(chapters []Chapter) []Chapter {
	for i, j := 0, len(chapters) - 1; i < j; i, j = i + 1, j - 1 {
		chapters[i], chapters[j] = chapters[j], chapters[i]
	}
	return chapters
}

func downloadAllValidChapters(chapters []Chapter, conf config.ConfigType) (lastDownload, lastName string) {
	for _, chap := range chapters {
		fmt.Printf("下载：%s ......", chap.name)
		content := downloadChapterContent(conf, chap.url)
		if content == "" {
			return
		}
		saveToFile(chap.name, content, conf.Filename + ".txt")
		lastDownload = chap.url
		lastName = chap.name
		fmt.Println("完成。")
	}
	return
}

func downloadChapterContent(conf config.ConfigType, chapterUrl string)(content string) {
	uri := conf.Url + chapterUrl
	doc, err := goquery.NewDocument(uri)
	if err == nil {
		// "#content"
		content = strings.Replace(doc.Find(conf.ContentReg).Text(), "请记住本书首发域名：wwww.4xiaoshuo.com。4小说网手机版阅读网址：m.4xiaoshuo.com", "", -1)
	} else {
		fmt.Println("error:", err.Error())
	}
	return
}

func saveToFile(title, content, filePath string) {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil {
		defer file.Close()
		file.WriteString(title + "\n" + content)
	} else {
		fmt.Println("打开文件错误:", err.Error())
	}
}

