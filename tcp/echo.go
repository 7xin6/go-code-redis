package tcp

import (
	"bufio"
	"context"
	"go-redis/lib/logger"
	"go-redis/lib/sync/atomic"
	"go-redis/lib/sync/wait"
	"io"
	"net"
	"sync"
	"time"
)

/**
  回复信息
*/

// EchoClient 客户端信息
type EchoClient struct {
	Conn    net.Conn
	Waiting wait.Wait
}

func (e *EchoClient) Close() error {
	e.Waiting.WaitWithTimeout(10 * time.Second)
	_ = e.Conn.Close()
	return nil
}

type EchoHandler struct {
	activeConn sync.Map // 记录当前连接数
	closing    atomic.Boolean
}

func MakeHandler() *EchoHandler {
	return &EchoHandler{}
}

// Handler 接受 Tcp 服务器传递过来的连接
func (handler *EchoHandler) Handler(ctx context.Context, conn net.Conn) {
	if handler.closing.Get() {
		_ = conn.Close()
	}
	client := &EchoClient{
		Conn: conn,
	}
	// 记录当前 客户端连接
	handler.activeConn.Store(client, struct{}{})
	reader := bufio.NewReader(conn)
	for {
		// 因为用户传输消息可能是断断续续的传输，所以需要一个 带有缓存区的bufio
		msg, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF { // 操作系统的结束符号，客户端退出了
				logger.Info("Connecting close")
				handler.activeConn.Delete(client)
			} else {
				logger.Warn(err)
			}
			return
		}
		// 代表一个 业务正在进行
		client.Waiting.Add(1)
		b := []byte(msg)
		_, _ = conn.Write(b)
		client.Waiting.Done()
	}
}

func (handler *EchoHandler) Close() error {
	logger.Info("handler shutting down")
	handler.closing.Set(true)

	// 关闭所有客户端
	handler.activeConn.Range(func(key, value any) bool {
		client := key.(*EchoClient)
		_ = client.Conn.Close()
		// 返回 true 就处理下一个
		return true
	})

	return nil
}
