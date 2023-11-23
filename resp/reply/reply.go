package reply

import (
	"bytes"
	"go-redis/interface/resp"
	"strconv"
)

var (
	nullBulkReplyBytes = []byte("$-1")
	CRLF               = "\r\n" // 结尾
)

type BulkReply struct {
	Arg []byte // "moody" "$5\r\nmoody\r\n"
}

func (b *BulkReply) ToBytes() []byte {
	if len(b.Arg) == 0 {
		return nullBulkReplyBytes
	}
	return []byte("$" + strconv.Itoa(len(b.Arg)) + CRLF + string(b.Arg) + CRLF)
}

func MakeBulkReply(arg []byte) *BulkReply {
	return &BulkReply{Arg: arg}
}

// MultiBulkReply 多个字符串回复
type MultiBulkReply struct {
	Args [][]byte
}

func MakeMultiBulkReply(arg [][]byte) *MultiBulkReply {
	return &MultiBulkReply{Args: arg}
}

func (m *MultiBulkReply) ToBytes() []byte {
	argLen := len(m.Args) // 行度
	var buf bytes.Buffer
	buf.WriteString("*" + strconv.Itoa(argLen) + CRLF)
	for _, arg := range m.Args {
		if arg == nil {
			buf.WriteString(string(nullBulkReplyBytes) + CRLF)
		} else {
			buf.WriteString("$" + strconv.Itoa(len(arg)) + CRLF + string(arg) + CRLF)
		}
	}
	return buf.Bytes()
}

// StatusReply 回复状态
type StatusReply struct {
	Status string
}

var theStatusReply = new(StatusReply)

func MakeStatusReply(status string) *StatusReply {
	return theStatusReply
}
func (s *StatusReply) ToBytes() []byte {
	return []byte("+" + s.Status + CRLF)
}

// IntReply 通用的数字回复
type IntReply struct {
	Code int64
}

func (i *IntReply) ToBytes() []byte {
	return []byte(":" + strconv.FormatInt(i.Code, 10) + CRLF)
}

func MakeIntReply(code int64) *IntReply {
	return &IntReply{Code: code}
}

// ErrorReply 错误回复接口
type ErrorReply interface {
	Error() string
	ToBytes() []byte
}

// StandardErrReply 一般化的错误回复
type StandardErrReply struct {
	Status string
}

func MakeStandardErrReply(status string) *StandardErrReply {
	return &StandardErrReply{Status: status}
}
func (s *StandardErrReply) ToBytes() []byte {
	return []byte("-" + s.Status + CRLF)
}

func (s *StandardErrReply) Error() string {
	return s.Status
}

// IsErrReply 判断回复是正常还是异常
func IsErrReply(reply resp.Reply) bool {
	return reply.ToBytes()[0] == '-'
}
