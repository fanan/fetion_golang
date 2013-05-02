package fetion

import (
    "sync"
    "net/http"
    "net/url"
)


type Jar struct {
    lk sync.Mutex
    cookies map[string][]*http.Cookie
}

func NewJar() *Jar {
    j := Jar{sync.Mutex{}, make(map[string][]*http.Cookie)}
    return &j
}

func (jar *Jar) SetCookies(u *url.URL, cookies []*http.Cookie)  {
    jar.lk.Lock()
    jar.cookies[u.Host] = cookies
    jar.lk.Unlock()
}

func (jar *Jar) Cookies(u *url.URL) []*http.Cookie {
    return jar.cookies[u.Host]
}
