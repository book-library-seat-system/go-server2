/***************************************************************************
Copyright(C) 2018
Author: huziang
Description: user方面服务端处理，包含user方面路由函数
	PS：本层不涉及逻辑判断，逻辑判断在user包部分
Date: 2018年5月4日 星期五 下午1:13
****************************************************************************/

package server

import (
	"errors"
	"net/http"

	"github.com/book-library-seat-system/go-server/entity/user"
	. "github.com/book-library-seat-system/go-server/util"
	"github.com/unrolled/render"
)

// UserReturnjson 用于返回student的模板Json
type StudentRtnJson struct {
	// 包含userItem属性
	user.Item
	// 包含错误信息
	ErrorRtnJson
}

// 返回错误表单
func errResponse(w http.ResponseWriter, formatter *render.Render) {
	if err := recover(); err != nil {
		var rtn ErrorRtnJson
		rtn.Errorcode, rtn.Errorinformation = HandleError(err)
		formatter.JSON(w, 500, rtn)
	}
}

// test
func test(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}

// 创建一个新的用户
func createStudentHandle(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer errResponse(w, formatter)

		// 解析json数据
		js := parseJSON(r)
		user.RegisterStudent(
			js.Get("ID").MustString(),
			js.Get("name").MustString(),
			js.Get("password").MustString(),
			js.Get("email").MustString(),
			js.Get("school").MustString())

		// 发回json
		formatter.JSON(w, http.StatusOK, ErrorRtnJson{})
	}
}

// 登录用户
func loginStudentHandle(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer errResponse(w, formatter)

		// 解析json
		js := parseJSON(r)
		// 解析cookie
		_, _, err := parseCookie(r)
		if err == nil {
			CheckErr(errors.New("6|学生当前处于登陆状态"))
		}

		pitem := user.LoginStudent(
			js.Get("ID").MustString(),
			js.Get("password").MustString())

		// 如果成功登录，设置cookie
		cookie := getCookie("ID", pitem.ID)
		http.SetCookie(w, cookie)
		cookie = getCookie("school", pitem.School)
		http.SetCookie(w, cookie)

		// 返回json信息
		formatter.JSON(w, http.StatusOK, ErrorRtnJson{})
	}
}

// 登出用户
func logoutStudentHandle(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer errResponse(w, formatter)

		// 解析cookie
		_, _, err := parseCookie(r)
		CheckNewErr(err, "7|用户当前未登陆")

		// 返回json信息
		formatter.JSON(w, http.StatusOK, ErrorRtnJson{})
	}
}

// 显示用户信息
func listStudentInfoHandle(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer errResponse(w, formatter)

		// 解析cookie
		ID, _, err := parseCookie(r)
		CheckNewErr(err, "7|用户当前未登陆")

		pitem := user.GetStudent(ID)

		// 解析json
		formatter.JSON(w, http.StatusOK, StudentRtnJson{
			Item: *pitem,
		})
	}
}
