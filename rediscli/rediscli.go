package rediscli

import (
	"encoding/hex"
	"errors"
	"fmt"
)

func Binary2string(inData []byte) (outData string) {
	outData += "\""

	for _, b := range inData {
		switch b {
		case '\\', '"':
			// \ -> \\
			// " -> \"
			outData += "\\"
			outData += string(b)
		case '\n':
			outData += "\\n"
		case '\r':
			outData += "\\r"
		case '\t':
			outData += "\\t"
		case '\a':
			outData += "\\a"
		case '\b':
			outData += "\\b"
		default:
			if 0x20 <= b && b <= 0x7E {
				// 可打印
				outData += string(b)
			} else {
				// 不可打印字符, 打印它的16进制
				outData += "\\x"
				outData += fmt.Sprintf("%02x", b)
			}
		}
	}

	outData += "\""
	return outData
}

func String2binary(inData string) (outData []byte, err error) {
	// 去掉开头结尾的"
	if len(inData) >= 2 && inData[0] == '"' && inData[len(inData)-1] == '"' {
		inData = inData[1:len(inData)-1]
	}

	for i:=0; i<len(inData); i++ {
		switch inData[i] {
		case '\\':
			i++
			if i >= len(inData) {
				fmt.Println("ERROR len data: ", inData)
				return outData, errors.New("error len")
			}

			switch inData[i] {
			case '\\':
				outData = append(outData, '\\')
			case '"':
				outData = append(outData, '"')
			case 'n':
				outData = append(outData, '\n')
			case 'r':
				outData = append(outData, '\r')
			case 't':
				outData = append(outData, '\t')
			case 'a':
				outData = append(outData, '\a')
			case 'b':
				outData = append(outData, '\b')
			case 'x':
				hexString := inData[i+1:i+3]
				i += 2
				hexByte, _ := hex.DecodeString(hexString)
				// fmt.Println("DEBUG hex string: ", hexString, ", hex byte: ", hexByte)
				outData = append(outData, hexByte...)
			}
		default:
			outData = append(outData, inData[i])
		}
	}

	return outData, nil
}
