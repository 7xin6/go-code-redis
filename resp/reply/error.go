package reply

type ErrReply struct {
	Status string
}



var errBytes = []byte("-err\r\n")

func (e *ErrReply) ToBytes() []byte {
	return errBytes
}
func (e *ErrReply) Error() string {
	return "Err"
}

func MakeErrReply(status string) *ErrReply {
	return &ErrReply{Status: status}
}

// UnKnowErrReply 未知错误
type UnKnowErrReply struct {
}

var unknownErrBytes = []byte("-Err unknown\r\n")

func (u *UnKnowErrReply) Error() string {
	return "Err unknown"
}

func (u *UnKnowErrReply) ToBytes() []byte {
	return unknownErrBytes
}

// ArgNumErrReply 参数异常
type ArgNumErrReply struct {
	Cmd string // 指令本身

}

func MakeArgNumErrReply(cmd string) *ArgNumErrReply {

	return &ArgNumErrReply{Cmd: cmd}
}

func (r *ArgNumErrReply) Error() string {
	return "ERR wrong number of arguments for '" + r.Cmd + "' command\\r\\n"
}

func (r *ArgNumErrReply) ToBytes() []byte {
	return []byte("-ERR wrong number of arguments for " + r.Cmd + "' command\\r\\n")
}

// SyntaxErrReply 语法错误
type SyntaxErrReply struct {
}

var theSyntaxErrReply = &SyntaxErrReply{}

func MakeSyntaxErrReply() *SyntaxErrReply {
	return theSyntaxErrReply
}

var syntaxErrBytes = []byte("-Err syntax error\r\n")

func (r *SyntaxErrReply) ToBytes() []byte {
	return syntaxErrBytes
}
func (r *SyntaxErrReply) Error() string {
	return "Err syntax error"
}

// WrongTypeErrReply 数据类型错误represents operation against a key holding the wrong kind of value
type WrongTypeErrReply struct{}

var wrongTypeErrBytes = []byte("-WRONGTYPE Operation against a key holding the wrong kind of value\r\n")

// ToBytes marshals redis.Reply
func (r *WrongTypeErrReply) ToBytes() []byte {
	return wrongTypeErrBytes
}

func (r *WrongTypeErrReply) Error() string {
	return "WRONGTYPE Operation against a key holding the wrong kind of value"
}

// ProtocolErr

// ProtocolErrReply 协议类型错误represents meeting unexpected byte during parse requests
type ProtocolErrReply struct {
	Msg string
}

// ToBytes marshals redis.Reply
func (r *ProtocolErrReply) ToBytes() []byte {
	return []byte("-ERR Protocol error: '" + r.Msg + "'\r\n")
}

func (r *ProtocolErrReply) Error() string {
	return "ERR Protocol error: '" + r.Msg
}
