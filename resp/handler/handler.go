package handler

import (
	"context"
	"go-redis/database"
	databaseface "go-redis/interface/database"
	"go-redis/lib/logger"
	"go-redis/lib/sync/atomic"
	"go-redis/resp/connection"
	"go-redis/resp/parser"
	"go-redis/resp/reply"
	"io"
	"net"
	"strings"
	"sync"
)

var (
	unknownErrReplyBytes = []byte("-ERR unknown\r\n")
)

// RespHandler 处理 RESP 协议的结构体
type RespHandler struct {
	activeConn sync.Map
	db         databaseface.DataBase
	closing    atomic.Boolean
}

func MakeHandler() *RespHandler {
	var db databaseface.DataBase
	db = database.NewDatabase()
	return &RespHandler{
		db: db,
	}
}

// closeClient 关闭一个连接
func (r *RespHandler) closeClient(client *connection.Connection) {
	_ = client.Close()
	r.db.AfterClientClose(client)
	r.activeConn.Delete(client)
}

func (r *RespHandler) Handler(ctx context.Context, conn net.Conn) {
	if r.closing.Get() {
		_ = conn.Close()
	}
	client := connection.NewConn(conn)
	r.activeConn.Store(client, struct{}{})
	// parser 已经独立工作了 启协程 读取conn 里面的报文 解析放入channel
	ch := parser.ParseStream(conn)
	for payload := range ch {
		// error
		if payload.Err != nil {
			if payload.Err == io.EOF ||
				payload.Err == io.ErrUnexpectedEOF ||
				strings.Contains(payload.Err.Error(), "use of closed network connection") { // 挥手关闭连接
				r.closeClient(client)
				logger.Info("connection closed: " + client.RemoteAddr().String())
			}
			// protocol error
			errReply := reply.MakeStatusReply(payload.Err.Error())
			err := client.Write(errReply.ToBytes())
			if err != nil {
				r.closeClient(client)
				logger.Info("connection closed: " + client.RemoteAddr().String())
				return
			}
			continue
		}
		//exec
		if payload.Data == nil {
			continue
		}
		res, ok := payload.Data.(*reply.MultiBulkReply)
		if !ok {
			logger.Error("require multi bulk reply")
			continue
		}
		result := r.db.Exec(client, res.Args)
		if result != nil {
			client.Write(result.ToBytes())
		} else {
			client.Write(unknownErrReplyBytes)
		}
	}
}

// Close 关闭所有连接
func (r *RespHandler) Close() error {
	logger.Info("handler shutting down")
	r.closing.Set(true)
	r.activeConn.Range(func(key interface{}, value interface{}) bool {
		client := key.(*connection.Connection)
		_ = client.Close()
		return true
	})
	r.db.Close()
	return nil
}
