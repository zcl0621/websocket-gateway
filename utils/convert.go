package utils

import (
	"bytes"
	"compress/zlib"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"strconv"
)

// StructToInterface convert struct to interface
func StructToInterface(obj interface{}) interface{} {
	vp := reflect.New(reflect.TypeOf(obj))
	vp.Elem().Set(reflect.ValueOf(obj))
	return vp.Interface()
}

func StructToBytes(obj interface{}) []byte {
	data, err := json.Marshal(obj)
	if err != nil {
		return nil
	}
	return data
}

// InterfaceToStruct convert interface to struct
func InterfaceToStruct(obj interface{}) interface{} {
	vp := reflect.ValueOf(obj)
	return vp.Elem().Interface()
}

// ByteToStruct convert byte to struct
func ByteToStruct(obj interface{}, data []byte) error {
	err := json.Unmarshal(data, obj)
	return err
}

func InterfaceToString(obj interface{}) string {
	variable, ok := obj.(string)
	if !ok {
		return ""
	}
	return variable
}

func InterfaceToInt(obj interface{}) int {
	variable, ok := obj.(float64)
	if !ok {
		return 0
	}
	return int(variable)
}

func InterfaceToInt8(obj interface{}) int8 {
	variable, ok := obj.(float64)
	if !ok {
		return 0
	}
	return int8(variable)
}

func StringToInt(str string) int {
	variable, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return variable
}

// zipString zip string
func ZipString(str string) string {
	return Base64Encode(zipStr(str))
}

func zipStr(origin string) (content string) {
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	w.Write([]byte(origin))
	w.Close()
	return b.String()
}

func UnZipString(zipContent string) string {
	var originInfo []byte
	newZipContent := Base64Decode(zipContent)
	var b bytes.Buffer
	b.WriteString(newZipContent)
	r, err := zlib.NewReader(&b)
	if err != nil {
		fmt.Println(" err : ", err)
	}
	defer r.Close()

	//r.Close()
	originInfo, err = ioutil.ReadAll(r)
	if err != nil {
		fmt.Println(" err : ", err)
	}

	return fmt.Sprintf("%s", originInfo)
}

// string to base64 string
func Base64Encode(str string) string {
	return b64.StdEncoding.EncodeToString([]byte(str))
}

func Base64Decode(str string) string {
	newStr, _ := b64.StdEncoding.DecodeString(str)
	return fmt.Sprintf("%s", newStr)
}
