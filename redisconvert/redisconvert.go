package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/liudong1994/goutil/rediscli"
	_ "github.com/liudong1994/goutil/redisconvert/pb"	// 提前读取pb目录下的所有pb文件, 方便后面全局找protobuf message类名
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
	"os"
)

// TODO 指定目录加载pb文件

func mustParseFile(s string) *descriptorpb.FileDescriptorProto {
	pb := new(descriptorpb.FileDescriptorProto)
	if err := prototext.Unmarshal([]byte(s), pb); err != nil {
		panic(err)
	}
	return pb
}

func main() {
	var pbfile string
	flag.StringVar(&pbfile , "pbf", "", "protobuf file name, package name must be 'pb', *pb.go file directory must be './pb'")

	var pbmessage string
	flag.StringVar(&pbmessage , "pbm", "", "protobuf message name")

	var rediscliData string
	flag.StringVar(&rediscliData , "rcd", "", "redis-cli string data")

	var jsonfile string
	flag.StringVar(&jsonfile , "json", "", "json file name")

	flag.Parse()

	if len(pbmessage) != 0 && len(rediscliData) != 0 {
		// 输入pb文件, redis-cli字符串数据, 自动解析pb 打印debugstring
		fmt.Println("redis-cli string CONVERT pb debug string")
		rediscli2pb(pbmessage, rediscliData)

	} else if len(pbmessage) != 0 && len(jsonfile) != 0 {
		// 输入pb文件, json文件, 自动转换为redis-cli字符串数据, 直接set即可
		fmt.Println("json CONVERT pb CONVERT redis-cli string")
		json2pb2rediscli(pbmessage, jsonfile)

	} else {
		flag.Usage()
		fmt.Printf("CONVERT json to pb(redis-cli string):  go run redisconvert.go -pbf \"test.proto\" -pbm \"TestInfo\" -json test.json\n")
		fmt.Printf("CONVERT pb(redis-cli string) to debug string:  go run redisconvert.go -pbf \"test.proto\" -pbm \"TestInfo\" -rcd \"\\b\\x01\\x12\\x0c\\n\\n1651939200\"\n")
	}

	return
}

func rediscli2pb(pbmessage string, redisData string) {
	// redisdata转换pb数据 打印
	redisBinaryData, _ := rediscli.String2binary(redisData)

	msg, err := genPBMessageByName(pbmessage)
	if err != nil {
		fmt.Printf("rediscli2pb gen pb message by name:%s err:%s\n", pbmessage, err)
		return
	}

	if err := proto.Unmarshal(redisBinaryData, msg); err != nil {
		fmt.Printf("rediscli2pb proto unmarshal err:%s\n", err)
		return
	}
	// fmt.Printf("rediscli2pb data:\n%s", proto.MarshalTextString(msg))

	pb2json := jsonpb.Marshaler{}
	jsonStr, _ := pb2json.MarshalToString(msg)

	var jsonDebug bytes.Buffer
	if err = json.Indent(&jsonDebug, []byte(jsonStr), "", "    "); err != nil {
		fmt.Printf("rediscli2pb json indent err:%s\n", err)
		return
	}

	fmt.Printf("rediscli2pb data:\n%s\n", jsonDebug.String())
}

func json2pb2rediscli(pbmessage string, jsonfile string) {
	// 找到pb message
	msg, err := genPBMessageByName(pbmessage)
	if err != nil {
		fmt.Printf("json2pb2rediscli gen protobuf message err:%s\n", err)
		return
	}

	// 读取json
	jsonData, err := readfile(jsonfile)
	if err != nil {
		fmt.Printf("json2pb2rediscli read json err:%s\n", err)
		return
	}

	// json转pb
	if err := jsonpb.UnmarshalString(jsonData, msg); err != nil {
		fmt.Printf("json2pb2rediscli json 2 protobuf err: %s\n", err)
		return
	}

	// pb转rediscli字符串
	pbData, err := proto.Marshal(msg)
	if err != nil {
		fmt.Printf("json2pb2rediscli proto marshal err:%s\n", err)
		return
	}

	output := rediscli.Binary2string(pbData)
	fmt.Printf("json2pb2rediscli data: \n%s\n", output)
	return
}

func genPBMessageByName(fullName string) (proto.Message, error) {
	msgName := protoreflect.FullName(fullName)
	msgType, err := protoregistry.GlobalTypes.FindMessageByName(msgName)
	if err != nil {
		return nil, err
	}

	return proto.MessageV1(msgType.New()), nil
}

func readfile(filename string) (string, error) {
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

