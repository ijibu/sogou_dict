//批量解析抓取下来的html文件，获取下载地址
/*
go build parse.go
./parse > parse.log		//重定向结果到parse.log文件中去
*/
package main

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	parse()
}

//处理html文件
func parse() {
	path := "./data/"
	filepath.Walk(path, func(path string, f os.FileInfo, e error) error {
		if f == nil {
			return e
		}
		if f.IsDir() {
			return nil
		}

		parseHtmlFile(path)

		return nil
	})
}

func parseHtmlFile(path string) {
	str := strings.Split(path, "\\")
	//解析词库ID
	dictId := strings.Replace(str[1], ".html", "", -1)
	fmt.Print(dictId)
	fmt.Print(",")
	//读取整个文件的内容
	content, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	content1 := string(content)

	//解析词库条数
	re, _ := regexp.Compile("\\<td width=\"50\\%\"\\>词条：([\\S\\s]+?)个\\</td\\>")
	ret := re.FindAllStringSubmatch(content1, -1)
	dictNum := ret[0][1]
	fmt.Print(dictNum)
	fmt.Print(",")

	//解析词库下载地址
	re, _ = regexp.Compile("\\<a onclick[\\S\\s]+? href=\"([\\S\\s]+?)\"[\\S\\s]+?\\</a\\>")
	ret = re.FindAllStringSubmatch(content1, -1)
	downloadUrl := ret[0][1]
	fmt.Print(downloadUrl)
	fmt.Print(",")

	//解析词库名字
	re, _ = regexp.Compile("name=(.*)")
	dowUrl, _ := url.QueryUnescape(downloadUrl)
	//fmt.Print(dowUrl)
	ret = re.FindAllStringSubmatch(dowUrl, -1)
	dictName := ret[0][1]
	fmt.Print(dictName)

	fmt.Print("\n")
}
