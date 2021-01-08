package lib

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var TimeLocation *time.Location
var TimeFormat = "2006-01-02 15:04:05"
var DateFormat = "2006-01-02"
var LocalIP = net.ParseIP("127.0.0.1")



func HttpGET(trace *TraceContext, urlString string, urlParams url.Values, msTimeout int, header http.Header) (*http.Response, []byte, error) {
	//startTime := time.Now().UnixNano()
	client := http.Client{
		Timeout: time.Duration(msTimeout) * time.Millisecond,
	}
	urlString = AddGetDataToUrl(urlString, urlParams)
	req, err := http.NewRequest("GET", urlString, nil)
	if err != nil {
		//Log.TagWarn(trace, DLTagHTTPFailed, map[string]interface{}{
		//	"url":       urlString,
		//	"proc_time": float32(time.Now().UnixNano()-startTime) / 1.0e9,
		//	"method":    "GET",
		//	"args":      urlParams,
		//	"err":       err.Error(),
		//})
		return nil, nil, err
	}
	if len(header) > 0 {
		req.Header = header
	}
	req = addTrace2Header(req, trace)
	resp, err := client.Do(req)
	if err != nil {
		//Log.TagWarn(trace, DLTagHTTPFailed, map[string]interface{}{
		//	"url":       urlString,
		//	"proc_time": float32(time.Now().UnixNano()-startTime) / 1.0e9,
		//	"method":    "GET",
		//	"args":      urlParams,
		//	"err":       err.Error(),
		//})
		return nil, nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		//Log.TagWarn(trace, DLTagHTTPFailed, map[string]interface{}{
		//	"url":       urlString,
		//	"proc_time": float32(time.Now().UnixNano()-startTime) / 1.0e9,
		//	"method":    "GET",
		//	"args":      urlParams,
		//	"result":    Substr(string(body), 0, 1024),
		//	"err":       err.Error(),
		//})
		return nil, nil, err
	}
	//Log.TagInfo(trace, DLTagHTTPSuccess, map[string]interface{}{
	//	"url":       urlString,
	//	"proc_time": float32(time.Now().UnixNano()-startTime) / 1.0e9,
	//	"method":    "GET",
	//	"args":      urlParams,
	//	"result":    Substr(string(body), 0, 1024),
	//})
	return resp, body, nil
}

func HttpPOST(trace *TraceContext, urlString string, urlParams url.Values, msTimeout int, header http.Header, contextType string) (*http.Response, []byte, error) {
	//startTime := time.Now().UnixNano()
	client := http.Client{
		Timeout: time.Duration(msTimeout) * time.Millisecond,
	}
	if contextType == "" {
		contextType = "application/x-www-form-urlencoded"
	}
	urlParamEncode := urlParams.Encode()
	req, err := http.NewRequest("POST", urlString, strings.NewReader(urlParamEncode))
	if len(header) > 0 {
		req.Header = header
	}
	req = addTrace2Header(req, trace)
	req.Header.Set("Content-Type", contextType)
	resp, err := client.Do(req)
	if err != nil {
		//Log.TagWarn(trace, DLTagHTTPFailed, map[string]interface{}{
		//	"url":       urlString,
		//	"proc_time": float32(time.Now().UnixNano()-startTime) / 1.0e9,
		//	"method":    "POST",
		//	"args":      Substr(urlParamEncode, 0, 1024),
		//	"err":       err.Error(),
		//})
		return nil, nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		//Log.TagWarn(trace, DLTagHTTPFailed, map[string]interface{}{
		//	"url":       urlString,
		//	"proc_time": float32(time.Now().UnixNano()-startTime) / 1.0e9,
		//	"method":    "POST",
		//	"args":      Substr(urlParamEncode, 0, 1024),
		//	"result":    Substr(string(body), 0, 1024),
		//	"err":       err.Error(),
		//})
		return nil, nil, err
	}
	//Log.TagInfo(trace, DLTagHTTPSuccess, map[string]interface{}{
	//	"url":       urlString,
	//	"proc_time": float32(time.Now().UnixNano()-startTime) / 1.0e9,
	//	"method":    "POST",
	//	"args":      Substr(urlParamEncode, 0, 1024),
	//	"result":    Substr(string(body), 0, 1024),
	//})
	return resp, body, nil
}

func HttpJSON(trace *TraceContext, urlString string, jsonContent string, msTimeout int, header http.Header) (*http.Response, []byte, error) {
	//startTime := time.Now().UnixNano()
	client := http.Client{
		Timeout: time.Duration(msTimeout) * time.Millisecond,
	}
	req, err := http.NewRequest("POST", urlString, strings.NewReader(jsonContent))
	if len(header) > 0 {
		req.Header = header
	}
	req = addTrace2Header(req, trace)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		//Log.TagWarn(trace, DLTagHTTPFailed, map[string]interface{}{
		//	"url":       urlString,
		//	"proc_time": float32(time.Now().UnixNano()-startTime) / 1.0e9,
		//	"method":    "POST",
		//	"args":      Substr(jsonContent, 0, 1024),
		//	"err":       err.Error(),
		//})
		return nil, nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		//Log.TagWarn(trace, DLTagHTTPFailed, map[string]interface{}{
		//	"url":       urlString,
		//	"proc_time": float32(time.Now().UnixNano()-startTime) / 1.0e9,
		//	"method":    "POST",
		//	"args":      Substr(jsonContent, 0, 1024),
		//	"result":    Substr(string(body), 0, 1024),
		//	"err":       err.Error(),
		//})
		return nil, nil, err
	}
	//Log.TagInfo(trace, DLTagHTTPSuccess, map[string]interface{}{
	//	"url":       urlString,
	//	"proc_time": float32(time.Now().UnixNano()-startTime) / 1.0e9,
	//	"method":    "POST",
	//	"args":      Substr(jsonContent, 0, 1024),
	//	"result":    Substr(string(body), 0, 1024),
	//})
	return resp, body, nil
}

func AddGetDataToUrl(urlString string, data url.Values) string {
	if strings.Contains(urlString, "?") {
		urlString = urlString + "&"
	} else {
		urlString = urlString + "?"
	}
	return fmt.Sprintf("%s%s", urlString, data.Encode())
}

func addTrace2Header(request *http.Request, trace *TraceContext) *http.Request {
	traceId := trace.TraceId
	cSpanId := NewSpanId()
	if traceId != "" {
		request.Header.Set("didi-header-rid", traceId)
	}
	if cSpanId != "" {
		request.Header.Set("didi-header-spanid", cSpanId)
	}
	trace.CSpanId = cSpanId
	return request
}

func GetMd5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func Encode(data string) (string, error) {
	h := md5.New()
	_, err := h.Write([]byte(data))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}


func NewTrace() *TraceContext {
	trace := &TraceContext{}
	trace.TraceId = GetTraceId()
	trace.SpanId = NewSpanId()
	return trace
}
func NewTraceByTag(logTag string) *TraceContext {
	trace := NewTrace()
	trace.LogTag = logTag
	return trace
}

func NewSpanId() string {
	timestamp := uint32(time.Now().Unix())
	ipToLong := binary.BigEndian.Uint32(LocalIP.To4())
	b := bytes.Buffer{}
	b.WriteString(fmt.Sprintf("%08x", ipToLong^timestamp))
	b.WriteString(fmt.Sprintf("%08x", rand.Int31()))
	return b.String()
}

//GetTraceId .生成TraceId
func GetTraceId() (traceId string) {
	return calcTraceId(LocalIP.String())
}

func calcTraceId(ip string) (traceId string) {
	now := time.Now()
	timestamp := uint32(now.Unix())
	timeNano := now.UnixNano()
	pid := os.Getpid()

	b := bytes.Buffer{}
	netIP := net.ParseIP(ip)
	if netIP == nil {
		b.WriteString("00000000")
	} else {
		b.WriteString(hex.EncodeToString(netIP.To4()))
	}
	b.WriteString(fmt.Sprintf("%08x", timestamp&0xffffffff))
	b.WriteString(fmt.Sprintf("%04x", timeNano&0xffff))
	b.WriteString(fmt.Sprintf("%04x", pid&0xffff))
	b.WriteString(fmt.Sprintf("%06x", rand.Int31n(1<<24)))
	b.WriteString("b0") // 末两位标记来源,b0为go

	return b.String()
}

//GetLocalIPs 获取本地Ip
func GetLocalIPs() (ips []net.IP) {
	interfaceAddr, err := net.InterfaceAddrs()
	if err != nil {
		return nil
	}
	for _, address := range interfaceAddr {
		ipNet, isValidIpNet := address.(*net.IPNet)
		if isValidIpNet && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ips = append(ips, ipNet.IP)
			}
		}
	}
	return ips
}

//InArrayString 字符串是否在数组中
func InArrayString(s string, arr []string) bool {
	for _, i := range arr {
		if i == s {
			return true
		}
	}
	return false
}

//Substr 字符串的截取
func Substr(str string, start int64, end int64) string {
	length := int64(len(str))
	if start < 0 || start > length {
		return ""
	}
	if end < 0 {
		return ""
	}
	if end > length {
		end = length
	}
	return string(str[start:end])
}

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


var ConfEnvPath string //配置文件夹
var ConfEnv string     //配置环境名 比如：dev prod test

// 解析配置文件目录
//
// 配置文件必须放到一个文件夹中
// 如：config=conf/dev/base.json 	ConfEnvPath=conf/dev	ConfEnv=dev
// 如：config=conf/base.json		ConfEnvPath=conf		ConfEnv=conf
func ParseConfPath(config string) error {
	path := strings.Split(config, "/")
	prefix := strings.Join(path[:len(path)-1], "/")
	ConfEnvPath = prefix
	ConfEnv = path[len(path)-2]
	return nil
}

//获取配置环境名
func GetConfEnv() string {
	return ConfEnv
}

func GetConfPath(fileName string) string {
	return ConfEnvPath + "/" + fileName + ".toml"
}

func GetConfFilePath(fileName string) string {
	return ConfEnvPath + "/" + fileName
}

//本地解析文件
func ParseLocalConfig(fileName string, st interface{}) error {
	path := GetConfFilePath(fileName)
	err := ParseConfig(path, st)
	if err != nil {
		return err
	}
	return nil
}

func ParseConfig(path string, conf interface{}) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("Open config %v fail, %v", path, err)
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("Read config fail, %v", err)
	}

	v := viper.New()
	v.SetConfigType("toml")
	v.ReadConfig(bytes.NewBuffer(data))
	if err := v.Unmarshal(conf); err != nil {
		return fmt.Errorf("Parse config fail, config:%v, err:%v", string(data), err)
	}
	return nil
}



