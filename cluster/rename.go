package cluster

import (
	"go-redis/interface/resp"
	"go-redis/resp/reply"
)
// rename k1 k2
func rename(cluster *ClusterDatabase, conn resp.Connection, cmdArgs [][]byte) resp.Reply {
	if len(cmdArgs) != 3 {
		return reply.MakeErrReply("ERR Wrong number args")
	}
	src := string(cmdArgs[1])
	dest := string(cmdArgs[2])

	srcPeer := cluster.peerPicker.PickNode(src) // 192.168.....:6379
	destPeer := cluster.peerPicker.PickNode(dest)
	if srcPeer != destPeer {
		return reply.MakeErrReply("ERR rename must within on peer")
	}
	return cluster.relay(srcPeer,conn,cmdArgs)
}