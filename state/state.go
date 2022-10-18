package state

import (
	"context"
	"sync"
	"time"

	"google.golang.org/grpc/resolver"
)

/*
有状态的结构体，搞清楚状态变化，就搞清楚一大半逻辑了
*/

//如果结构体是有状态的，需要注意结构体状态的变化
type hostResolver struct {
	m      sync.Mutex
	cc     resolver.ClientConn
	wait   <-chan struct{}
	cancel context.CancelFunc
	closed bool
}

//为了避免同时更新，更新的时候先判断状态，如果发现正在更新就取消，然后等取消完了再执行更新
func (r *hostResolver) UpdateHosts(app string) error {
	r.m.Lock()
	defer r.m.Unlock()

	if r.closed {
		return nil
	}
	r.cancelHostsUpdate()

	wait := make(chan struct{})
	go func() {
		//更新完以后关闭
		time.Sleep(1 * time.Second)
		close(wait)
	}()
	r.wait = wait
	r.closed = false
	return nil
}

func (r *hostResolver) cancelHostsUpdate() {
	if r.cancel != nil {
		r.cancel()
		r.cancel = nil
	}
	if r.wait != nil {
		<-r.wait
		r.wait = nil
	}
}
