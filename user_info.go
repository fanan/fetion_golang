package fetion

import (
    "encoding/json"
)

type UserInfo struct {
    IdContact int `json:"idContact"`
    MobileNumer string `json:"mobileNo"`
    BasicServiceStatus int `json:"basicServiceStatus"`
}

func parseUserInfo (contents *[]byte)(ui *UserInfo) {
    ui = new(UserInfo)
    err := json.Unmarshal(*contents, ui)
    if err != nil {
        println(err.Error())
        return nil
    }
    return ui
}
