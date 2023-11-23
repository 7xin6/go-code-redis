package database

import (
	"go-redis/aof"
	"go-redis/config"
	"go-redis/interface/resp"
	"go-redis/lib/logger"
	"go-redis/resp/reply"
	"strconv"
	"strings"
)

type Database struct {
	dbSet      []*DB
	aofHandler *aof.AofHandler
}

func NewDatabase() *Database {
	database := &Database{}
	if config.Properties.Databases == 0 {
		config.Properties.Databases = 16
	}
	database.dbSet = make([]*DB, config.Properties.Databases)
	for i := range database.dbSet {
		db := MakeDB()
		db.index = i
		database.dbSet[i] = db
	}
	if config.Properties.AppendOnly {
		aofHandler, err := aof.NewAofHandler(database)
		if err != nil {
			panic(err)
		}
		database.aofHandler = aofHandler
		for _, db := range database.dbSet { // db 在for循环上 db的值会变，但是地址不会变
			// db = dbSet[0]
			// db = dbSet[15] 发生逃逸，逃逸到了堆上
			singleDb := db                         // 还是会逃逸到堆上，但是 在for 循环内部 值会变，地址也会变
			singleDb.addAof = func(line CmdLine) { // 闭包问题
				database.aofHandler.AddAof(singleDb.index, line)
			}
		}
	}
	return database
}

// set k v
// get k v
// select 2

func (data *Database) Exec(client resp.Connection, args [][]byte) resp.Reply {
	defer func() {
		err := recover()
		if err != nil {
			logger.Error(err)
		}
	}()
	cmdName := strings.ToLower(string(args[0]))
	if cmdName == "select" {
		if len(args) != 2 {
			return reply.MakeArgNumErrReply("select")
		}
		return execSelect(client, data, args[1:])
	}
	dbIndex := client.GetDBIndex()
	db := data.dbSet[dbIndex]
	return db.Exec(client, args)
}

func (data *Database) Close() {

}

func (data *Database) AfterClientClose(c resp.Connection) {

}

// 切换 db           select 2
func execSelect(c resp.Connection, database *Database, args [][]byte) resp.Reply {
	dbIndex, err := strconv.Atoi(string(args[0]))
	if err != nil {
		return reply.MakeErrReply("ERR invalid DB index")
	}
	if dbIndex >= len(database.dbSet) {
		return reply.MakeErrReply("ERR DB index is out of range")
	}
	c.SelectDB(dbIndex)
	return reply.MakeOkReply()
}
