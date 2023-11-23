package resp

type Connection interface {
	Write([]byte) error // 给客户端回复消息
	GetDBIndex() int    // 查询现在使用的是哪个 DB库
	SelectDB(int)       // 切换 DB库
	Close() error
}
