package r2m

import "fmt"

var JobChannel = make(chan Job, 100)

type Job struct {
	Key      string
	Data     string
	QuitChan chan bool
}

func (j Job) Start() {
	go func() {
		for {
			select {
			case job := <-JobChannel:
				SaveCache(job.Key, job.Data)
			case <-j.QuitChan:
				fmt.Println("job exit")
				return
			}
		}
	}()
}

func (j Job) Stop() {
	go func() {
		j.QuitChan <- true
	}()
}
