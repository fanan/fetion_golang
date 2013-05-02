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
)

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
	fmt.Println(string(body))
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
	return nil
}

func (f *Fetion) getGroupList() {
	resp, _ := f.client.Get(FetionGroupListURL)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(body))
	gil := parseGroupListInfo(&body)
	f.groupids = gil.parseGroups()
}

func (f *Fetion) getFriends(groupid int) error {
	url := fmt.Sprintf("%s%d", FetionFriendsListURL, groupid)
	resp, err := f.client.Get(url)
	defer resp.Body.Close()
	if err != nil {
		println(err.Error())
		return err
	}
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		println(err.Error())
		return err
	}
	uli := parseUserListInfo(&contents)
	if uli == nil {
		println("parse error")
		return errors.New(ERROR_JSON_PARSE)
	}
	for _, user := range uli.Users {
		if user.BasicServiceStatus == 1 && user.MobileNumer != "" {
			f.friends[user.MobileNumer] = user.IdContact
		}
	}
	return nil
}

func (f *Fetion) Logout() error {
	resp, err := f.client.PostForm(FetionLogoutURL, url.Values{})
	defer resp.Body.Close()
	if err != nil {
		println(err.Error())
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		println(err.Error())
		return err
	}
	ls := ParseLogoutStatus(&body)
	if ls == nil || ls.Tip != "退出成功" {
		return errors.New(ERROR_LOGOUT)
	}
	return nil
}

func (f *Fetion) BuildUserDb() {
	for _, groupid := range f.groupids {
		f.getFriends(groupid)
	}
	//for i, _ := range f.groupids {
		//println(i, "groupid", <-finished, "finished")
	//}
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
    fmt.Println("*******************")
	fmt.Println(string(contents))
    fmt.Println("*******************")
	return 0, nil
}

func (f *Fetion) SendSms(msg string, users []string) (err error) {
	if msg == "" {
		return errors.New(ERROR_EMPTY_MSG)
	}
	fetionContacts := make([]string, 0)
	for _, user := range users {
		id, err := f.QueryFriendId(user)
		if err == nil {
			fetionContacts = append(fetionContacts, strconv.Itoa(id))
		}
	}
	if len(fetionContacts) == 0 {
		return errors.New(ERROR_EMPTY_USERS)
	}
	touserids := fmt.Sprintf(",%s", strings.Join(fetionContacts, ","))
	data := url.Values{"touserid": {touserids}, "msg": {msg}}
	resp, err := f.client.PostForm(FetionSendGroupSMSURL, data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	sss := ParseSendSMSStatus(&contents)
	if sss == nil {
		return errors.New(ERROR_JSON_PARSE)
	}
	if sss.Info != "发送成功" {
		return errors.New(ERROR_SENDSMS)
	}
	return nil
}

func (f *Fetion) SendOneself(msg string) (err error) {
	l := []string{f.mobileNumber}
	return f.SendSms(msg, l)
}
