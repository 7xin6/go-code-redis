package tcp

import (
	"context"
	"go-redis/interface/tcp"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Config struct {
	Address string
}

// ListenAndServeWithSignal 监听新连接 并 携带信号量
func ListenAndServeWithSignal(cfg *Config, handler tcp.Handler) error {
	closeChan := make(chan struct{})
	//获取 操作系统传递的信号量 给 closeChan 发布信号
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		sig := <-sigChan
		switch sig {
		case syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			closeChan <- struct{}{}
		}
	}()
	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return err
	}
	// todo logg 打印
	ListenAndServe(listener, handler, closeChan)
	return nil
}

// ListenAndServe 接受新连接
func ListenAndServe(listener net.Listener, handler tcp.Handler, closeChan <-chan struct{}) {

	go func() {
		<-closeChan // 没有内容 会阻塞在这
		// todo 打印日志
		_ = listener.Close()
		_ = handler.Close()
	}()

	// 关闭连接
	defer func() {
		_ = listener.Close()
		_ = handler.Close()
	}()
	// 创建一个新的上下文环境变量
	ctx := context.Background()
	var waitDone sync.WaitGroup
	for true {
		conn, err := listener.Accept()
		if err != nil {
			break
		}
		// todo logg 打印
		waitDone.Add(1)
		// 开启一个协程 用来处理连接
		go func() {
			defer waitDone.Done()
			handler.Handler(ctx, conn)
		}()
	}
	waitDone.Wait()
}
