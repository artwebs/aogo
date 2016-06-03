package fsnotify

import (
	"log"
	"os/exec"
	"sync"
	"time"

	"github.com/go-fsnotify/fsnotify"
)

type FSNotify struct {
	sleeptime int
	path      string
	cmd       string
	args      []string
}

func NewFSNotify() *FSNotify {
	return &FSNotify{sleeptime: 30, path: "./", cmd: "", args: []string{}}
}

func (this *FSNotify) Run() {
	Watch, err := fsnotify.NewWatcher()
	if err != nil {
		log.Println("Init monitor error: ", err.Error())
		return
	}
	if err := Watch.Add(this.path); err != nil {
		log.Println("Add monitor path error: ", this.path)
		return
	}
	var (
		cron bool = false
		lock      = new(sync.Mutex)
	)
	for {
		select {
		case event := <-Watch.Events:
			log.Printf("Monitor event %s", event.String())
			if !cron {
				cron = true
				go func() {
					T := time.After(time.Second * time.Duration(this.sleeptime))
					<-T
					if err := call(this.cmd, this.args...); err != nil {
						log.Println(err)
					}
					lock.Lock()
					cron = false
					lock.Unlock()
				}()
			}
		case err := <-Watch.Errors:
			log.Println(err)
			return
		}
	}
}

func call(programe string, args ...string) error {
	cmd := exec.Command(programe, args...)
	buf, err := cmd.Output()
	if err != nil {
		return err
	}
	log.Printf("\n%s\n", string(buf))
	return nil
}
