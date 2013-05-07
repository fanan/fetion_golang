package fetion

import (
	//"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
    "log"
)

const (
	FetionLoginURL        = "http://f.10086.cn/im5/login/loginHtml5.action"
	FetionLogoutURL       = "http://f.10086.cn/im5/index/logoutsubmit.action"
	FetionGroupListURL    = "http://f.10086.cn/im5/index/loadGroupContactsAjax.action"
	FetionFriendsListURL  = "http://f.10086.cn/im5/index/contactlistView.action?fromUrl=&idContactList="
	FetionSendGroupSMSURL = "http://f.10086.cn/im5/chat/sendNewGroupShortMsg.action"
	FetionFriendQureyURL  = "http://f.10086.cn/im5/index/searchFriendsByQueryKey.action"
)

const (
	ERROR_EMPTY_MSG   = "empty message"
	ERROR_JSON_PARSE  = "json parse error"
	ERROR_LOGOUT      = "log out error"
	ERROR_EMPTY_USERS = "empty users"
	ERROR_SENDSMS     = "sendsms error"
    max_http_connections = 5
)

var connectionChannels chan int

func init () {
    connectionChannels = make(chan int, max_http_connections)
    log.SetFlags(log.Lmicroseconds | log.Lshortfile | log.Ldate)
}

type Fetion struct {
	mobileNumber string
	password     string
	userid       string
	client       *http.Client
	groupids     []int
	friends      map[string]int
}




func NewFetion(mobileNumber, password string) (f *Fetion) {
	f = new(Fetion)
	f.mobileNumber = mobileNumber
	f.password = password
	f.client = &http.Client{nil, nil, NewJar()}
	f.groupids = make([]int, 0)
	f.friends = make(map[string]int)
	return f
}

func (f *Fetion) Login() error {
	data := url.Values{"m": {f.mobileNumber}, "pass": {f.password}, "captchaCode": {""}, "checkCodeKey": {"null"}}
    log.Println("start login")
	resp, err := f.client.PostForm(FetionLoginURL, data)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		println(err.Error())
		return err
	}
	//fmt.Println(string(body))
	//u, _ := url.Parse("http://f.10086.cn")
	//fmt.Printf("%v", f.client.Jar.Cookies(u))
	ls := parseLoginStatus(&body)
	if ls == nil {
		return errors.New(ERROR_JSON_PARSE)
	}
	if ls.Tip != "" || ls.State != "200" {
		return errors.New(ls.Tip)
	}
	f.userid = ls.UserId
	f.friends[f.mobileNumber], _ = strconv.Atoi(f.userid)
    log.Println("login succeed")
	return nil
}

func (f *Fetion) getGroupList() {
    log.Println("start get group list")
	resp, _ := f.client.Get(FetionGroupListURL)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(body))
	gil := parseGroupListInfo(&body)
	f.groupids = gil.parseGroups()
    log.Println("finish get group list")
}

func (f *Fetion) getFriends(groupid int, ret chan bool){
	url := fmt.Sprintf("%s%d", FetionFriendsListURL, groupid)
    connectionChannels <- 1
    log.Println("start get friends in group", groupid)
	resp, err := f.client.Get(url)
	defer resp.Body.Close()
	if err != nil {
		log.Println(err)
        <-connectionChannels
        ret <- false
		return
	}
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
        log.Println(err)
        <-connectionChannels
        ret <- false
		return
	}
	uli := parseUserListInfo(&contents)
	if uli == nil {
        log.Println(ERROR_JSON_PARSE)
        <-connectionChannels
        ret <- false
		return
	}
    number := 0
	for _, user := range uli.Users {
		//if user.BasicServiceStatus == 1 && user.MobileNumer != "" {
        //whar does BasicServiceStatus mean????
		if user.MobileNumer != "" {
			f.friends[user.MobileNumer] = user.IdContact
            number++
		}
	}
    log.Println("finish get", number, "friends in group", groupid)
    <-connectionChannels
    ret <- true
	return
}

func (f *Fetion) Logout() error {
    log.Println("start logout")
	resp, err := f.client.PostForm(FetionLogoutURL, url.Values{})
	defer resp.Body.Close()
	if err != nil {
		log.Println(err)
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
        log.Println(err)
		return err
	}
	ls := ParseLogoutStatus(&body)
	if ls == nil || ls.Tip != "退出成功" {
        log.Println("logout failed!", ls.Tip)
		return errors.New(ERROR_LOGOUT)
	}
    log.Println("logout succeed")
	return nil
}

func (f *Fetion) BuildUserDb() {
    log.Println("start build user db")
    f.getGroupList()
    numberOfErrors := len(f.groupids)
    returnChannel := make(chan bool)
    flag := true
	for _, groupid := range f.groupids {
        go f.getFriends(groupid, returnChannel)
	}
    for i:=0; i < numberOfErrors; i++ {
        b := <-returnChannel
        flag = flag && b
    }
    log.Println("finished build user db")
}

func (f *Fetion) QueryFriendId(mobileNumber string) (int, error) {
	id, ok := f.friends[mobileNumber]
	if ok {
		return id, nil
	}
	resp, err := f.client.PostForm(FetionFriendQureyURL, url.Values{"queryKey": []string{mobileNumber}})
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
    //fmt.Println("*******************")
	//fmt.Println(len(contents), string(contents))
    //fmt.Println("*******************")
    uli := parseUserListInfo(&contents)
    if uli == nil || len(uli.Users) != 1 {
        return 0, nil
    }
    user := uli.Users[0]
    if user.BasicServiceStatus != 1 || user.MobileNumer == "" {
        return 0, nil
    }
    f.friends[mobileNumber] = user.IdContact
    return user.IdContact, nil
}

func (f *Fetion) SendSms(msg string, users []string) (err error) {
    log.Println(msg, users)
	if msg == "" {
		return errors.New(ERROR_EMPTY_MSG)
	}
	fetionContacts := make([]string, 0)
    //log.Println(fetionContacts)
	for _, user := range users {
		id, err := f.QueryFriendId(user)
		if err == nil {
            //log.Println(id)
			fetionContacts = append(fetionContacts, strconv.Itoa(id))
            log.Println(fetionContacts)
		}
	}
	if len(fetionContacts) == 0 {
		return errors.New(ERROR_EMPTY_USERS)
	}
	touserids := fmt.Sprintf(",%s", strings.Join(fetionContacts, ","))
	data := url.Values{"touserid": {touserids}, "msg": {msg}}
    log.Println(data)
	resp, err := f.client.PostForm(FetionSendGroupSMSURL, data)
	if err != nil {
        log.Println(err)
		return err
	}
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
        log.Println(err)
		return err
	}
    //log.Println(string(contents))
	sss := ParseSendSMSStatus(&contents)
	if sss == nil {
		return errors.New(ERROR_JSON_PARSE)
	}
	if sss.Info != "发送成功" {
		return errors.New(ERROR_SENDSMS)
	}
    log.Println("message sent!")
	return nil
}

func (f *Fetion) SendOneself(msg string) (err error) {
	l := []string{f.mobileNumber}
	return f.SendSms(msg, l)
}

func (f *Fetion) ListFriends()  {
    for k, v := range f.friends{
        fmt.Println(k, ":", v)
    }
}
