package main

import (
	"fmt"
	"girlImage/src/tool"
	"github.com/PuerkitoBio/goquery"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var imgList = map[string]string{
	"1": "https://www.mzitu.com/page/",
	"2": "https://www.mzitu.com/hot/page/",
	"3": "https://www.mzitu.com/best/page/",
	"4": "https://www.mzitu.com/xinggan/page/",
	"5": "https://www.mzitu.com/japan/page/",
	"6": "https://www.mzitu.com/taiwan/page/",
	"7": "https://www.mzitu.com/mm/page/",
}

var url string

var ImgPath = "/images/meizitu/"

var Header = map[string]string{
	"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.92 Safari/537.36",
	"referer":    "https://www.mzitu.com/page/",
}

var pageIndex = 1
var delay = 500 //每次请求延迟
var wg sync.WaitGroup

func main() {
	fmt.Println("1. 最新图片")
	fmt.Println("2. 最热图片")
	fmt.Println("3. 推荐图片")
	fmt.Println("4. 性感妹子")
	fmt.Println("5. 日本妹子")
	fmt.Println("6. 台湾妹子")
	fmt.Println("7. 清纯妹子")
	fmt.Print("请输入要下载的图片类型:")
	defer tool.End()
	var imgType string
	fmt.Scanln(&imgType)
	url = imgList[imgType]
	if url == "" {
		url = imgList["1"]
	}
	fmt.Println(url)
	dir, _ := os.Getwd()
	ImgPath = dir + ImgPath
	startDownload()
}

func startDownload() {
	currentUrl := url + strconv.Itoa(pageIndex)
	resp, err := tool.Get(currentUrl, Header)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	doc.Find("#pins").Find("span a").Each(func(i int, s *goquery.Selection) {
		detailUrl, _ := s.Attr("href")
		name := s.Text()
		getDetail(detailUrl, ImgPath+name)
		time.Sleep(time.Millisecond * 100)
	})
	pageIndex++
	startDownload()
}

func getDetail(url string, dirName string) {
	resp, err := tool.Get(url, Header)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	img := doc.Find(".main-image img")
	firstImg, _ := img.Attr("src")
	lastDoc := doc.Find(".pagenavi a")
	lastIndex := lastDoc.Length() - 2
	totalCount := 0
	lastDoc.Each(func(i int, s *goquery.Selection) {
		if i == lastIndex {
			lastUrl, _ := s.Attr("href")
			urlArr := strings.Split(lastUrl, "/")
			totalCount, _ = strconv.Atoi(urlArr[len(urlArr)-1])
		}
	})
	tool.CheckDir(dirName)
	for i := 1; i <= totalCount; i++ {
		imgUrl := strings.Replace(firstImg, "01.jpg", fmt.Sprintf("%02d", i)+".jpg", 1)
		wg.Add(1)
		tool.SaveFile(imgUrl, dirName, Header, delay, &wg)
	}
}
