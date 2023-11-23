package reply

type PongReply struct {
}

var pongBytes = []byte("+PONG\r\n")

func (p *PongReply) ToBytes() []byte {
	return pongBytes
}

var thePongReply = new(PongReply)

func MakePongReply() *PongReply {
	return thePongReply
}

type OKReply struct {
}

var okBytes = []byte("+OK\r\n")

func (o *OKReply) ToBytes() []byte {
	return okBytes
}

// 在本地内存创建一个OKReply,就不需要反复的创建对象，可以重复利用
var theOKReply = new(OKReply)

func MakeOkReply() *OKReply {
	return theOKReply
}

// NullBulkReply 空的块回复
type NullBulkReply struct {
}

var nullBulkBytes = []byte("$-1\r\n") // null

func (p *NullBulkReply) ToBytes() []byte {
	return nullBulkBytes
}

var theNullBulkReply = new(NullBulkReply)

func MakeNullBulkReply() *NullBulkReply {
	return theNullBulkReply
}

// EmptyMultiBulkReply 空的块回复
type EmptyMultiBulkReply struct {
}

var emptyMultiBulkBytes = []byte("*0\r\n") // null

func (p *EmptyMultiBulkReply) ToBytes() []byte {
	return emptyMultiBulkBytes
}

var emptyMultiBulkReply = new(EmptyMultiBulkReply)

func MakeEmptyMultiBulkReply() *EmptyMultiBulkReply {
	return emptyMultiBulkReply
}

type NoReply struct {
}

var noBytes = []byte("")

func (p *NoReply) ToBytes() []byte {
	return noBytes
}

var noReply = new(NoReply)

func makeNoReply() *NoReply {
	return noReply
}
