package ziface

type IMessage interface {
	GetDataLength() uint32
	GetMsgID() uint32
	GetData() []byte

	SetMsgID(uint32)
	SetData([]byte)
	SetDataLength(uint32)
}
