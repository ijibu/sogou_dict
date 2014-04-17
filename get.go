//批量获取搜狗拼音词库数据。
package main

import (
	//"bufio"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	//"time"
)

var num *int = flag.Int("n", 40, "please input a num like 1024")
var cupNum int = runtime.NumCPU()
var c = make(chan int, cupNum) //设置最大并发

func main() {
	flag.Usage = show_usage
	flag.Parse()

	if *num == 0 {
		show_usage()
		return
	}

	runtime.GOMAXPROCS(cupNum) //设置cpu的核的数量，从而实现高并发

	var (
		logFileDir string
		downDir    string
	)

	logFileDir = "./log/"
	downDir = "./data/"
	if !isDirExists(logFileDir) { //目录不存在，则进行创建。
		err := os.MkdirAll(logFileDir, 777) //递归创建目录，linux下面还要考虑目录的权限设置。
		if err != nil {
			panic(err)
		}
	}
	if !isDirExists(downDir) { //目录不存在，则进行创建。
		err := os.MkdirAll(downDir, 777) //递归创建目录，linux下面还要考虑目录的权限设置。
		if err != nil {
			panic(err)
		}
	}

	logfile, _ := os.OpenFile(logFileDir+"ijibu.log", os.O_RDWR|os.O_CREATE, 0)
	logger := log.New(logfile, "\r\n", log.Ldate|log.Ltime|log.Llongfile)

	defer logfile.Close()

	for i := 0; i < cupNum; i++ {
		go work(logger, logfile, downDir)
	}
	start()
}

func start() {
	id := 1
	for {
		c <- id
		id++
	}
}

func work(logger *log.Logger, logfile *os.File, downDir string) {
	for {
		id := <-c
		code := strconv.Itoa(id)
		getDict(logger, logfile, downDir, code)
	}
}

func getDict(logger *log.Logger, logfile *os.File, downDir string, code string) {
	getUrl := "http://pinyin.sogou.com/dict/cell.php?id=" + code
	client := &http.Client{CheckRedirect: myRedirect}
	resp, err := client.Get(getUrl)
	if err == nil {
		if resp.StatusCode == 200 {
			logger.Println(logfile, code+":sucess"+strconv.Itoa(resp.StatusCode))
			fmt.Println(code + ":sucess")

			fileName := downDir + code + ".html"
			//不加os.O_RDWR的话，在linux下面无法写入文件。
			f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0666) //其实这里的 O_RDWR应该是 O_RDWR|O_CREATE，也就是文件不存在的情况下就建一个空文件，但是因为windows下还有BUG，如果使用这个O_CREATE，就会直接清空文件，所以这里就不用了这个标志，你自己事先建立好文件。
			if err != nil {
				panic(err)
			}

			defer f.Close()

			//处理返回的html内容
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("http read error")
			}

			src := string(body)

			//将HTML标签全转换成小写
			re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
			src = re.ReplaceAllStringFunc(src, strings.ToLower)

			/*
				//去除STYLE
				re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
				src = re.ReplaceAllString(src, "")

				//去除SCRIPT
				re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
				src = re.ReplaceAllString(src, "")
			*/

			//<div class="dlinfobox">
			re, _ = regexp.Compile("\\<div class=\"dlinfobox\"\\>([\\S\\s]+?)\\</div\\>")
			ret := re.FindString(src)

			buf := []byte(ret)
			f.Write(buf)
		} else {
			logger.Println(logfile, code+":http get StatusCode"+strconv.Itoa(resp.StatusCode))
			fmt.Println(code + ":" + strconv.Itoa(resp.StatusCode))
		}
		defer resp.Body.Close()
	} else {
		logger.Println(logfile, code+":http get error"+code)
		fmt.Println(code + ":error")
	}
}

func isDirExists(path string) bool {
	fi, err := os.Stat(path)

	if err != nil {
		return os.IsExist(err)
	} else {
		return fi.IsDir()
	}
}

func show_usage() {
	fmt.Fprintf(os.Stderr,
		"Usage: %s [-n=<num>] \n"+
			"       <command> [<args>]\n\n",
		os.Args[0])
	fmt.Fprintf(os.Stderr,
		"Flags:\n")
	flag.PrintDefaults()
}

func myRedirect(req *http.Request, via []*http.Request) (e error) {
	return errors.New(req.URL.String())
}
