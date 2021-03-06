/***************************************************************************
Copyright(C) 2018
Author: huziang
Description: 常用的工具函数包
Date: 2018年5月4日 星期五 下午1:08
****************************************************************************/

package util

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	ini "github.com/vaughan0/go-ini"
)

var logDivPath = "src/github.com/book-library-seat-system/go-server/log"
var logFilePath = "/" + time.Now().Format("2006-01-02") + ".txt"

// GetGOPATH 获得用户环境的gopath
func GetGOPATH() string {
	var sp string
	if runtime.GOOS == "windows" {
		sp = ";"
	} else {
		sp = ":"
	}
	goPath := strings.Split(os.Getenv("GOPATH"), sp)
	for _, v := range goPath {
		if _, err := os.Stat(filepath.Join(v, "/src/github.com/book-library-seat-system/go-server/util/util.go")); !os.IsNotExist(err) {
			return v
		}
	}
	return ""
}

// getFileHandle 获得文件handle
func getFileHandle() *os.File {
	abspath := GetGOPATH() + logDivPath + logFilePath
	if _, err := os.Open(abspath); err != nil {
		os.Create(abspath)
	}

	// 以追加模式打开文件,并向文件写入
	fi, _ := os.OpenFile(abspath, os.O_RDWR|os.O_APPEND, 0)
	return fi
}

// CheckErr panic错误
func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

// CheckDBErr 本地输出原错误，然后抛出自定义错误
func CheckNewErr(err error, str string) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		panic(errors.New(str))
	}
}

// HandleError 处理错误，并且返回成int，string形式
func HandleError(err interface{}) (int, string) {
	strs := strings.Split(err.(error).Error(), "|")
	if len(strs) != 2 {
		return 200, "未定义错误"
	}
	errcode, err := strconv.Atoi(strs[0])
	if err != nil {
		return 200, "未定义错误"
	}
	return errcode, strs[1]
}

// MD5Hash MD5哈希函数
func MD5Hash(input string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(input))
	return hex.EncodeToString(md5Ctx.Sum(nil))
}

// ReadFromIniFile 从配置文件中读取信息并返回
func ReadFromIniFile(blockname string, rowname string) string {
	file, err := ini.LoadFile(filepath.Join(GetGOPATH(), "/src/github.com/book-library-seat-system/go-server/config.ini"))
	CheckErr(err)
	str, ok := file.Get(blockname, rowname)
	if !ok {
		panic(errors.New("202|读取配置文件发生错误"))
	}
	return str
}

// String2Int string转换成int
func String2Int(str string) int {
	i, _ := strconv.Atoi(str)
	return i
}

// Bool2Int bool转换成int
func Bool2Int(b bool) int {
	if b {
		return 1
	} else {
		return 0
	}
}

// Int2Bool int转换成bool
func Int2Bool(i int) bool {
	return i > 0
}
