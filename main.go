// Copyright 2012 JessonChan and his Girl-friend. All rights reserved.
// Use of this source code is governed by a BSD-style 'license',
// that means "Take it down to the copy center and make as many copies as you want"
// you can find the BSD license on the website http://www.fsf.org/licensing/licenses/index_html#OriginalBSD

/*
   可以直接使用
   gorun main.go t Hello,终端发微博 
   来发微博，也可以打包成二进制的。
   注意：程序中的client_id和client_secret可以替换成自己的应用。保留只为了大家方便运行
   TODO 错误检查、Bug修复、帮助信息、功能完善、代码注释 
*/
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type AccessToken struct {
	Access_token string
	Expires_in   int
	Uid          string
}

func (a AccessToken) String() string {
	return a.Access_token
}

var (
	access_token string = ""
	config_path  string = os.Getenv("HOME") + "/.clifaweibo_oauth_2"
)

const (
	faweibo_status      string = "/tmp/.clifaweibo_status"
	client_id           string = "3319129991"
	client_secret       string = "9762d236559d87f4e6e9a0e5c966e4fd"
	grant_type          string = "authorization_code"
	redirect_uri        string = "http://clifaweibo.sinaapp.com/ilovecliofubuntu.php"
	access_token_url    string = "https://api.weibo.com/oauth2/access_token"
	statuses_update_url string = "https://api.weibo.com/2/statuses/update.json"
	statuses_upload_url string = "https://upload.api.weibo.com/2/statuses/upload.json"
)

func get_access_token_from_file() bool {
	config_file, err := os.OpenFile(config_path, os.O_RDONLY, 0644)
	if err != nil {
		return false
	}
	defer config_file.Close()
	a := new(AccessToken)
	err = json.NewDecoder(config_file).Decode(a)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	access_token = a.Access_token
	return true
}

func get_access_token_from_http() bool {
	fmt.Printf("请复制下面的地址在浏览器中打开完成授权\nhttp://t.cn/zlUnqCC\n请输入(右键粘贴)授权码:")
	var grant_code string
	fmt.Scanf("%s", &grant_code)
	grant_url_values := make(url.Values)
	grant_url_values.Set("client_id", client_id)
	grant_url_values.Set("client_secret", client_secret)
	grant_url_values.Set("grant_type", grant_type)
	grant_url_values.Set("code", grant_code)
	grant_url_values.Set("redirect_uri", redirect_uri)
	r, _ := http.Post(access_token_url, "application/x-www-form-urlencoded", strings.NewReader(grant_url_values.Encode()))
	defer r.Body.Close()
	a := new(AccessToken)
	err := json.NewDecoder(r.Body).Decode(a)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	access_token = a.Access_token
	body, _ := ioutil.ReadAll(r.Body)
	config_file, err := os.OpenFile(config_path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return false
	}
	defer config_file.Close()
	config_file.Write(body)
	return true
}

func send_text_weibo(text string) bool {
	weibo_values := make(url.Values)
	weibo_values.Set("access_token", access_token)
	weibo_values.Set("status", text)
	r, err := http.Post(statuses_update_url, "application/x-www-form-urlencoded", strings.NewReader(weibo_values.Encode()))
	defer r.Body.Close()
	if err != nil {
		return false
	}
	return true
}
func send_pic_weibo(text, file_path string) bool {
	var b bytes.Buffer
	body_writer := multipart.NewWriter(&b)
	body_writer.WriteField("access_token", access_token)
	body_writer.WriteField("status", text)

	pic_writer, _ := body_writer.CreateFormFile("pic", file_path)
	fh, _ := os.Open(file_path)
	io.Copy(pic_writer, fh)
	form_type := body_writer.FormDataContentType()
	body_writer.Close()
	resp, err := http.Post(statuses_upload_url, form_type, &b)
	defer resp.Body.Close()
	if err != nil {
		return false
	}
	return true
}

func main() {
	if false == get_access_token_from_file() {
		if false == get_access_token_from_http() {
			return
		}
	}
	argc := len(os.Args)
	if argc < 3 {
		return
	}
	switch os.Args[1] {
	case "-t", "t":
		text := ""
		for i := 2; i < argc; i++ {
			text += os.Args[i]
			if i < argc-1 {
				text += " "
			}
		}
		if true == send_text_weibo(text) {
			fmt.Println("发送成功")
		}
	case "-p", "p":
		if true == send_pic_weibo("分享图片", os.Args[argc-1]) {
			fmt.Println("发送成功")
		}
	case "-tp", "-pt", "pt", "tp":
		text := ""
		for i := 2; i < argc-1; i++ {
			text += os.Args[i]
			if i < argc-1 {
				text += " "
			}
		}
		if true == send_pic_weibo(text, os.Args[argc-1]) {
			fmt.Println("发送成功")
		}
	default:
		return
	}
}
