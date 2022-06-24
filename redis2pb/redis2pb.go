package main

import (
	"flag"
	"fmt"
	"github.com/liudong1994/goutil/redisconvert"
	"os"
	_ "github.com/liudong1994/goutil/redis2pb/pb"	// 编译所有pb(调用init), 支持后面根据pb name获取对应pb message
)

func main() {
	var pbmessage string
	flag.StringVar(&pbmessage , "pbm", "", "protobuf message name")

	var rediscliData string
	flag.StringVar(&rediscliData , "rcd", "", "redis-cli string data, WARN:use single quotes parameters")

	var jsonfile string
	flag.StringVar(&jsonfile , "json", "", "json file name")

	flag.Parse()

	if len(pbmessage) != 0 && len(rediscliData) != 0 {
		// 输入pb文件, redis-cli字符串数据, 自动解析pb 打印debugstring
		fmt.Println("redis-cli string CONVERT json")
		pbData, _ := redisconvert.Rediscli2pb2json(pbmessage, rediscliData)
		fmt.Println(pbData)

	} else if len(pbmessage) != 0 && len(jsonfile) != 0 {
		// 输入pb文件, json文件, 自动转换为redis-cli字符串数据, 直接set即可
		fmt.Println("json CONVERT pb CONVERT redis-cli string")
		jsonData, _ := readFile(jsonfile)
		rediscliData, _ := redisconvert.Json2pb2rediscli(pbmessage, jsonData)
		fmt.Println(rediscliData)

	} else {
		flag.Usage()
		fmt.Printf("CONVERT json to pb(redis-cli string):  go run redis2pb.go -pbm 'TestInfo' -json 'test.json'\n")
		fmt.Printf("CONVERT pb(redis-cli string) to debug string:  go run redis2pb.go -pbm 'TestInfo' -rcd '\\b\\x01\\x12\\x0c\\n\\n1651939200'\n")
	}

	return
}

func readFile(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	filesize := fileInfo.Size()
	buffer := make([]byte, filesize)
	if _, err = file.Read(buffer); err != nil {
		fmt.Println(err)
		return "", err
	}

	return string(buffer), nil
}

