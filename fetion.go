package fetion

/*
本包主要用来模拟wap飞信进行登陆,发送短信和登出.无需验证码,不存在任何破解行为.
使用时首先Login,用SeachUser找出所需发送对象的touserid,然后SendSms,最后Logout即可.
*/

import (
    "net/http"
    "net/url"
    "bufio"
    "fmt"
    "strings"
    "regexp"
)

const (
    loginURL = "http://f.10086.cn/im/login/inputpasssubmit1.action"
    logoutURL = "http://f.10086.cn/im/index/logout.action"
    smsSelfURL = "http://f.10086.cn/im/user/sendMsgToMyselfs.action"
    smsURL = "http://f.10086.cn/im/chat/sendShortMsg.action?touserid="
    searchURL = "http://f.10086.cn/im/index/searchOtherInfoList.action"
    /*4表示隐身*/
    loginstatus = "4"
)

var (
    touseridReg = regexp.MustCompile(`/im/chat/toinputMsg\.action\?touserid=(\d+)`)
)

func Echo(a interface{}) {
    fmt.Printf("%+v\n", a)
}

func Login(mobileNumber, password string) (cookies []*http.Cookie) {
    data := make(url.Values)
    data.Set("onlinestatus", loginstatus)
    data.Set("m", mobileNumber)
    data.Set("pass", password)
    resp, err := http.DefaultClient.PostForm(loginURL, data)
    if err != nil {
        Echo(err)
    }
    for _,cookie := range(resp.Cookies()) {
        cookies = append(cookies, cookie)
    }
    return
}

func SendSms(cookies []*http.Cookie, touserid string, msg string) {
    var realurl string
    /*touserid为空表示给自己发短信*/
    if touserid == "" {
        realurl = smsSelfURL
    } else {
        realurl = smsURL + touserid
    }
    data := make(url.Values)
    data.Set("msg", msg)
    s := data.Encode()
    req,_ := http.NewRequest("POST", realurl, strings.NewReader(s))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
    for _,cookie := range cookies {
        req.AddCookie(cookie)
    }
    http.DefaultClient.Do(req)
    /*如果有需要的话,可以在这里分析一下response.Body,判断是否发送成功*/
    return
}

func SearchUser(cookies []*http.Cookie, mobileNumber string) (touserid string) {
    data := make(url.Values)
    data.Set("searchText", mobileNumber)
    s := data.Encode()
    req,_ := http.NewRequest("POST", searchURL, strings.NewReader(s))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
    for _,cookie := range cookies {
        req.AddCookie(cookie)
    }
    resp,err := http.DefaultClient.Do(req)
    if err != nil {
        Echo(err)
        return
    }
    buf := bufio.NewReader(resp.Body)
    defer resp.Body.Close()
    A: for {
        b,_,err := buf.ReadLine()
        if err != nil {
            break A
        }
        if touseridReg.Match(b) {
            touserid = touseridReg.FindStringSubmatch(string(b))[1]
            return
        }
    }
    return
}

func Logout(cookies []*http.Cookie) {
    req,_ := http.NewRequest("GET", logoutURL, nil)
    for _,cookie := range cookies {
        req.AddCookie(cookie)
    }
    http.DefaultClient.Do(req)
}

