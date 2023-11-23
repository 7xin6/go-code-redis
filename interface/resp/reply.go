package resp

// Reply 代表一切数据 客户端的回复
type Reply interface {
	ToBytes() []byte
}
