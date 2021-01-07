package lib

import (
	"bytes"
	"encoding/json"
	jsoniter "github.com/json-iterator/go"
	"log"
	"strings"
)

// 定义JSON操作
//var (
//	json              = jsoniter.ConfigCompatibleWithStandardLibrary
//	JSONMarshal       = json.Marshal
//	JSONUnmarshal     = json.Unmarshal
//	JSONMarshalIndent = json.MarshalIndent
//	JSONNewDecoder    = json.NewDecoder
//	JSONNewEncoder    = json.NewEncoder
//)

// JSONMarshalToString JSON编码为字符串
func JSONMarshalToString(v interface{}) string {
	//return ParseParams(v.(map[string]interface{}))
	log.Printf("vvv====%v" ,v)
	s, err := jsoniter.MarshalToString(v)
	log.Println("JSONMarshalToString====" +s)
	if err != nil {
		return ""
	}
	ss := strings.Replace(s,"\\\"","'",-1)
	sss := strings.Replace(ss,"\"","'",-1)
	//log.Println("JSONMarshalToString_old====" +s)
	log.Println("JSONMarshalToString_new====" +sss)

	return sss
}

func DisableEscapeHtml(data interface{}) (string, error) {
	bf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(bf)
	jsonEncoder.SetEscapeHTML(false)
	if err := jsonEncoder.Encode(data); err != nil {
		return "", err
	}
	return bf.String(), nil
}