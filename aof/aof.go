package aof

import (
	"go-redis/config"
	databaseface "go-redis/interface/database"
	"go-redis/lib/logger"
	"go-redis/lib/utils"
	"go-redis/resp/connection"
	"go-redis/resp/parser"
	"go-redis/resp/reply"
	"io"
	"os"
	"strconv"
)

type CmdLine = [][]byte

const aofBufferSize = 1 << 16

type payload struct {
	cmdLine CmdLine
	dbIndex int
}

// AofHandler 全局只有一个
type AofHandler struct {
	database    databaseface.DataBase // 业务核心信息
	aofChan     chan *payload         // channel充当缓冲区
	aofFile     *os.File              // aofFile文件
	aofFileName string
	currentDB   int // 记录上一条指令存放在哪个 DB上判断需不需要切换
}

// NewAofHandler :
func NewAofHandler(database databaseface.DataBase) (*AofHandler, error) {
	handler := &AofHandler{}

	handler.aofFileName = config.Properties.AppendFilename
	handler.database = database
	// LoadAof
	handler.LoadAof()
	aofFile, err := os.OpenFile(handler.aofFileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	handler.aofFile = aofFile
	// channel
	handler.aofChan = make(chan *payload, aofBufferSize)
	go func() {
		handler.handleAof()
	}()
	return handler, nil
}

// AddAof Add payload(set k v) -> aofChan
func (handler *AofHandler) AddAof(dbIndex int, cmd CmdLine) {
	if config.Properties.AppendOnly && handler.aofChan != nil {
		handler.aofChan <- &payload{
			cmdLine: cmd,
			dbIndex: dbIndex,
		}
	}

}

// handleAof  payload(set k v) <- aofChan(落盘)
func (handler *AofHandler) handleAof() {
	// TODO: payload(set k v) <- aofChan(落盘)
	handler.currentDB = 0
	for p := range handler.aofChan {
		if p.dbIndex != handler.currentDB {
			data := reply.MakeMultiBulkReply(utils.ToCmdLine("select", strconv.Itoa(p.dbIndex))).ToBytes()
			_, err := handler.aofFile.Write(data)
			if err != nil {
				logger.Error(err)
				continue
			}
			handler.currentDB = p.dbIndex
		}
		data := reply.MakeMultiBulkReply(p.cmdLine).ToBytes()
		_, err := handler.aofFile.Write(data)
		if err != nil {
			logger.Error(err)
		}
	}
}

// LoadAof (读盘)
func (handler *AofHandler) LoadAof() {
	file, err := os.Open(handler.aofFileName)
	if err != nil {
		logger.Error(err)
		return
	}

	defer file.Close()
	ch := parser.ParseStream(file)
	fackConn := &connection.Connection{}
	for p := range ch {
		if p.Err != nil {
			if p.Err == io.EOF {
				break
			}
			logger.Error("parse error: " + p.Err.Error())
			continue
		}
		if p.Data == nil {
			logger.Error("empty payload")
			continue
		}
		r, ok := p.Data.(*reply.MultiBulkReply)
		if !ok {
			logger.Error("need multi bulk")
			continue
		}
		rep := handler.database.Exec(fackConn, r.Args)
		if reply.IsErrReply(rep) {
			logger.Error("exec err", err)
		}
	}
}
