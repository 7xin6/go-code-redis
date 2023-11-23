package tcp

import (
	"context"
	"net"
)

/**
处理 TCP连接
*/

type Handler interface {
	Handler(ctx context.Context, conn net.Conn)
	Close() error
}
