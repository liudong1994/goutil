package redisconvert

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/liudong1994/goutil/rediscli"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

func Rediscli2pb2json(pbMsgName string, redisData string) (string, error) {
	// redisdata转换pb数据 打印
	redisBinaryData, _ := rediscli.String2binary(redisData)

	pbMsg, err := genPBMessageByName(pbMsgName)
	if err != nil {
		fmt.Printf("rediscli2pb gen pb message by name:%s err:%s\n", pbMsgName, err)
		return "", errors.New("gen pb message by name err")
	}

	if err := proto.Unmarshal(redisBinaryData, pbMsg); err != nil {
		fmt.Printf("rediscli2pb proto unmarshal err:%s\n", err)
		return "", errors.New("proto unmarshal err")
	}
	// fmt.Printf("rediscli2pb data:\n%s", proto.MarshalTextString(msg))

	pb2json := jsonpb.Marshaler{}
	jsonStr, _ := pb2json.MarshalToString(pbMsg)

	var jsonDebug bytes.Buffer
	if err = json.Indent(&jsonDebug, []byte(jsonStr), "", "    "); err != nil {
		fmt.Printf("rediscli2pb json indent err:%s\n", err)
		return "", errors.New("json indent err")
	}

	return jsonDebug.String(), nil
}

func Json2pb2rediscli(pbMsgName string, jsonData string) (string, error) {
	// 找到pb message
	pbMsg, err := genPBMessageByName(pbMsgName)
	if err != nil {
		fmt.Printf("json2pb2rediscli gen protobuf message err:%s\n", err)
		return "", errors.New("gen protobuf message err")
	}

	// json转pb
	if err := jsonpb.UnmarshalString(jsonData, pbMsg); err != nil {
		fmt.Printf("json2pb2rediscli json 2 protobuf err: %s\n", err)
		return "", errors.New("json 2 protobuf err")
	}

	// pb转rediscli字符串
	pbData, err := proto.Marshal(pbMsg)
	if err != nil {
		fmt.Printf("json2pb2rediscli proto marshal err:%s\n", err)
		return "", errors.New("proto marshal err")
	}

	output := rediscli.Binary2string(pbData)
	return output, nil
}

func genPBMessageByName(fullName string) (proto.Message, error) {
	msgName := protoreflect.FullName(fullName)
	msgType, err := protoregistry.GlobalTypes.FindMessageByName(msgName)
	if err != nil {
		return nil, err
	}

	return proto.MessageV1(msgType.New()), nil
}

