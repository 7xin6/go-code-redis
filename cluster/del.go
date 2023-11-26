package cluster

import (
	"go-redis/interface/resp"
	"go-redis/resp/reply"
)

// del k1 k2 k3 k4 k5
func del(cluster *ClusterDatabase, conn resp.Connection, cmdArgs [][]byte) resp.Reply  {
	replies := cluster.broadcast(conn, cmdArgs) // 广播出去了,若有些删除成功了，有些报错，不应该回滚么？
	var errReply reply.ErrorReply
	var deleted int64 = 0
	for _, r := range replies {
		if reply.IsErrReply(r) {
			errReply = r.(reply.ErrorReply)
			break
		}
		intReply,ok := r.(*reply.IntReply)
		if !ok {
			errReply = reply.MakeErrReply("error")
		}
		deleted += intReply.Code
	}
	if errReply == nil {
		return reply.MakeIntReply(deleted)
	}
	return reply.MakeErrReply("error: "+errReply.Error())
}
