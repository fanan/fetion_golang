package fetion

import (
    "encoding/json"
    "log"
)

type loginStatus struct {
    UserId string `json:"idUser"`
    State string `json:"loginstate"`
    Tip string `json:"tip"`
}


//func init() {
    //log.SetFlags(log.LstdFlags| log.Lshortfile)
//}

func parseLoginStatus (contents *[]byte) (ls *loginStatus) {
    ls = new(loginStatus)
    err := json.Unmarshal(*contents, ls)
    if err != nil {
        log.Println(err)
        return nil
    }
    return ls
}


type logoutStatus struct {
    Tip string `json:"tip"`
}

func ParseLogoutStatus (contents *[]byte) (ls *logoutStatus) {
    ls = new(logoutStatus)
    err := json.Unmarshal(*contents, ls)
    if err != nil {
        log.Println(err)
        return nil
    }
    return ls
}

type sendSMSStatus struct {
    Info string `json:"info"`
}

func ParseSendSMSStatus (contents *[]byte) (sss *sendSMSStatus) {
    sss = new(sendSMSStatus)
    err := json.Unmarshal(*contents, sss)
    if err != nil {
        log.Println(err)
        return nil
    }
    return sss
}
