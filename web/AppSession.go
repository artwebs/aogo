package web

import (
	"github.com/astaxie/beego/session"
	"net/http"
)

var (
	appSession *Session
)

type Store interface {
	Set(key, value interface{}) error     //set session value
	Get(key interface{}) interface{}      //get session value
	Delete(key interface{}) error         //delete session value
	SessionID() string                    //back current sessionID
	SessionRelease(w http.ResponseWriter) // release the resource & save data to provider & return the data
	Flush() error                         //delete all data
}

type Session struct {
	manager *session.Manager
	store   Store
}

func InitSession() *Session {
	if appSession == nil {
		appSession = &Session{}
		appSession.manager, _ = session.NewManager("memory", `{"cookieName":"gosessionid", "enableSetCookie,omitempty": true, "gclifetime":3600, "maxLifetime": 3600, "secure": false, "sessionIDHashFunc": "sha1", "sessionIDHashKey": "", "cookieLifeTime": 3600, "providerConfig": ""}`)
		go appSession.manager.GC()
	}
	return appSession
}

func (this *Session) Start(w http.ResponseWriter, r *http.Request) {
	this.store, _ = this.manager.SessionStart(w, r)
}

func (this *Session) Release(w http.ResponseWriter) {
	this.store.SessionRelease(w)
}

func (this *Session) Flush() error {
	return this.store.Flush()
}

func (this *Session) Set(key, value interface{}) {
	this.store.Set(key, value)
}

func (this *Session) Get(key interface{}) interface{} {
	return this.store.Get(key)
}
