package database

import "go-redis/interface/resp"

type CmdLine = [][]byte

// DataBase redis 业务层接口
type DataBase interface {
	Exec(client resp.Connection, args [][]byte) resp.Reply //执行业务
	Close()
	AfterClientClose(c resp.Connection)
}

type DataEntity struct {
	Data interface{}
}
