package cluster

import "go-redis/interface/resp"

func execSelect(cluster *ClusterDatabase, conn resp.Connection, cmdArgs [][]byte) resp.Reply {
	return cluster.db.Exec(conn,cmdArgs)
}