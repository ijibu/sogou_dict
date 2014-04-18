package main

import (
	"fmt"
	"log"
	"os"
)

const (
	startPy = 0x1540 //拼音表偏移
	//startChinese = 0x2628 //汉语词组表偏移
)

func main() {
	scel2txt("./scel/1001-41.scel")
}

func scel2txt(fileName string) {
	//读取整个文件的内容
	f, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	var b []byte = make([]byte, 12)
	f.ReadAt(b, 0)
	if string(b) != "\x40\x15\x00\x00\x44\x43\x53\x01\x01\x00\x00\x00" {
		log.Fatalf("确认你选择的文件 \"%s\"是搜狗(.scel)词库? \n", fileName)
	}

	len1 := 0x338 - 0x130
	len2 := 0x540 - 0x338
	len3 := 0xd40 - 0x540
	len4 := startPy - 0xd40

	dicName := make([]byte, len1)
	dicType := make([]byte, len2)
	dicDesc := make([]byte, len3)
	dicDemo := make([]byte, len4)

	f.ReadAt(dicName, 0x130)
	f.ReadAt(dicType, 0x338)
	f.ReadAt(dicDesc, 0x540)
	f.ReadAt(dicDemo, 0xd40)

	fmt.Println("词库名：" + unicode2utf8str(dicName, len1))
	fmt.Println("词库类型：" + unicode2utf8str(dicType, len2))
	fmt.Println("描述信息：" + unicode2utf8str(dicDesc, len3))
	fmt.Println("词库示例：" + unicode2utf8str(dicDemo, len4))

	//getPyTable(data[startPy:startChinese])
	//getChinese(data[startChinese:])
	return
}

func getPyTable(data []byte) {
	return
}

func getChinese(data []byte) {
	return
}

func unicode2utf8str(input []byte, insize int) (outstr string) {
	outstr = "\\0"

	for i := 0; i < insize/2; i++ {
		out := unicode2utf8char(int(input[i]))
		outstr += string(out)
	}
	return
}

func unicode2utf8char(in int) (out int) {
	if in >= 0x0000 && in <= 0x007f {
		out = in
		return
	} else if in >= 0x0080 && in <= 0x07ff {
		out = 0xc0 | (in >> 6)
		out++
		out = 0x80 | (in & (0xff >> 2))
		return
	} else if in >= 0x0800 && in <= 0xffff {
		out = 0xe0 | (in >> 12)
		out++
		out = 0x80 | (in >> 6 & 0x003f)
		out++
		out = 0x80 | (in & (0xff >> 2))
		return
	}
	fmt.Println("输入的不是short吧,解析有问题\n")
	return
}
