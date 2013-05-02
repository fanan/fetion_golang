package fetion

import (
    "encoding/json"
    "strconv"
    "strings"
)

type groupListInfo struct {
    Total int `json:"total"`
    FriendGroupIds string `jsoin:"total"`
}

func parseGroupListInfo (contens *[]byte) (gli *groupListInfo) {
    gli = new(groupListInfo)
    err := json.Unmarshal(*contens, gli)
    if err != nil {
        println(err.Error())
        return nil
    }
    return gli
}

func (gli *groupListInfo) parseGroups() []int {
    ignored_group_ids := []int{9998, 9999}
    _groupids := strings.Split(gli.FriendGroupIds, ",")
    groupids := make([]int, 0)
    for _, s := range _groupids {
        i, _ := strconv.Atoi(s)
        flag := true
        for _, b := range ignored_group_ids {
            if b == i {
                flag = false
                break
            }
        }
        if flag {
            groupids = append(groupids, i)
        }
    }
    return groupids
}
