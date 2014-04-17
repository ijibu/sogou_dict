//批量下载搜狗词库数据。
/*
该程序会使用经过parse.go处理生成的文件parse.log
*/
package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

const (
	UA       = "Mozilla/5.0 (Windows NT 6.1; rv:28.0) Gecko/20100101 Firefox/28.0"
	HOST     = "download.pinyin.sogou.com"
	REFERER  = "http://pinyin.sogou.com/dict/cell.php?id="
	FILE     = "parse.log"
	SAVEPATH = "scel/"
)

func main() {
	log.Printf("载入词库下载文件 %s", FILE)
	dictFile, err := os.Open(FILE)
	defer dictFile.Close()
	if err != nil {
		log.Fatalf("无法词库下载文件 \"%s\" \n", FILE)
	}

	reader := bufio.NewReader(dictFile)

	// 逐行读入分词
	for {
		dict, _ := reader.ReadString('\n')

		if len(dict) == 0 {
			// 文件结束
			break
		}
		dictInfo := strings.TrimSpace(dict)
		dictArr := strings.Split(dictInfo, ",")
		download(dictArr[0], dictArr[1], dictArr[2])
	}
}

func download(dictId string, dictNum string, downUrl string) {
	fileName := SAVEPATH + dictId + "-" + dictNum + ".scel"
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0666) //其实这里的 O_RDWR应该是 O_RDWR|O_CREATE，也就是文件不存在的情况下就建一个空文件，但是因为windows下还有BUG，如果使用这个O_CREATE，就会直接清空文件，所以这里就不用了这个标志，你自己事先建立好文件。
	if err != nil {
		panic(err)
	}

	defer f.Close()

	var req http.Request
	req.Method = "GET"
	req.Close = true
	req.URL, _ = url.Parse(downUrl)

	header := http.Header{}
	header.Set("User-Agent", UA)
	header.Set("Host", HOST)
	header.Set("Referer", REFERER+dictId)
	req.Header = header
	resp, err := http.DefaultClient.Do(&req)
	if err == nil {
		if resp.StatusCode == 200 {
			fmt.Println(dictId + ":sucess")
			_, err = io.Copy(f, resp.Body)
			if err != nil {
				panic(err)
			}
		} else {
			fmt.Println(dictId + ":" + strconv.Itoa(resp.StatusCode))
		}
		defer resp.Body.Close()
	} else {
		fmt.Println(dictId + ":error")
	}
}
