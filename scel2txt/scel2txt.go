package main

import (
	"fmt"
	"log"
	"os"
	"unicode/utf8"
)

const (
	startPy = 0x1540 //拼音表偏移
	//startChinese = 0x2628 //汉语词组表偏移
)

func main() {
	scel2txt("../data/scel/1001-41.scel")
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

	dicName := make([]byte, 65)
	dicType := make([]byte, 65)
	dicDesc := make([]byte, 65)
	dicDemo := make([]byte, 65)

	f.ReadAt(dicName, 0x130)
	f.ReadAt(dicType, 0x338)
	f.ReadAt(dicDesc, 0x540)
	f.ReadAt(dicDemo, 0xd40)
	if string(dicName[0:8]) != "\xA0\x5B\x69\x72\xCD\x8B\x47\x6C" {
		fmt.Println("词库名错误")
	}

	fmt.Println("词库名：" + string(unicode2utf8str(dicName)))
	fmt.Println("词库类型：" + string(unicode2utf8str(dicType)))
	fmt.Println("描述信息：" + string(unicode2utf8str(dicDesc)))
	fmt.Println("词库示例：" + string(unicode2utf8str(dicDemo)))

	//getPyTable(data[startPy:startChinese])
	//getChinese(data[startChinese:])
	return
}

func getPyTable(data []byte) {
	return
}

func getChinese(data []byte) {
	pos := 0
	length := len(data)
	for pos < length {

	}
	return
}

func unicode2utf8str(input []byte) (outstr []byte) {
	size := len(input)
	for i := 0; i < size-1; i += 2 {
		var b uint16                //input是byte数组，即uint8类型，要组合成两个字节长的16进制数，就得用uinit16
		b = uint16(input[i+1]) << 8 //作为高字节，向左移动8位，即一个字节
		b |= uint16(input[i])       //作为低字节
		out := unicode2utf8char(b)
		for j := 0; j < len(out); j++ {
			outstr = append(outstr, out[j])
		}
	}
	return outstr
}

func unicode2utf8char(in uint16) (out []byte) {
	n := 0
	if in >= 0x0000 && in <= 0x007f {
		n = 1
	} else if in >= 0x0080 && in <= 0x07ff {
		n = 2
	} else if in >= 0x0800 && in <= 0xffff {
		n = 3
	}
	if n == 0 {
		fmt.Println("输入的不是short吧,解析有问题\n")
	}

	buf := make([]byte, n)
	size := utf8.EncodeRune(buf, rune(in))
	out = buf[:size]
	return
}
