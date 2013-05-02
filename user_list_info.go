package fetion

import (
    "encoding/json"
)

type UserListInfo struct {
    Total int `json:"total"`
    Users []UserInfo `json:"contacts"`
}

func parseUserListInfo (contents *[]byte)(*UserListInfo) {
    uli := new(UserListInfo)
    err := json.Unmarshal(*contents, uli)
    if err != nil {
        println(err.Error())
        return nil
    }
    return uli
}
