fetion golang 接口
======

## 原理
模拟wap飞信登陆，无需输入验证码，适合用作定期发送短信、命令行操作等。

如果有某次登陆密码输入错误，直接登录会需要输入验证码。可以在网页上手动登录一遍，之后再使用本库登陆就不需要输入验证码了。


## 安装

```bash
go get -u github.com/fanan/fetion_golang
```


## 使用
```go

package main

import (
	"fmt"
    "github.com/fanan/fetion_golang"
)

func main() {
	mobileNumber := "13888888888"
	password := "88888888"
	f := NewFetion(mobileNumber, password)
	err := f.Login()
	if err != nil {
		fmt.Println(err)
	} else {
		f.getGroupList()
	}
	defer f.Logout()
	f.getGroupList()
	f.BuildUserDb()
    fmt.Println("total", len(f.groupids), "groups", len(f.friends), "users.")
    users := []string{"12345678901", "98765432109"}
    msg := "Hello 世界"
    f.SendSms(msg, users)
    f.SendOneself("发送成功")
}

```
