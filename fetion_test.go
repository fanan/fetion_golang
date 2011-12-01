package fetion_test

import (
    "fetion"
    "testing"
    "time"
)

const (
    mobileNumber = "your mobilephone number"
    password = "your fetion password"
    testMobile = "your friend mobilephone number"
)

func TestFetion(t *testing.T) {
    cookies := fetion.Login(mobileNumber, password)
    /*user := fetion.SearchUser(cookies, testMobile)*/
    /*fetion.SendSms(cookies, user, "your message here")*/
    time.Sleep(2*1e9)
    fetion.SendSms(cookies, "", "a message to yourself")
    fetion.Logout(cookies)
}
