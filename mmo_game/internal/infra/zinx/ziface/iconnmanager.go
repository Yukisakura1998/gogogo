package ziface

type IConnManager interface {
	Add(conn IConnection)
	Remove(conn IConnection)
	Get(ID string) (IConnection, error)
	Len() int
	ClearConn()
}
