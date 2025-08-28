package znet

type Message struct {
	Id     uint32
	Length uint32
	Data   []byte
}

func (m *Message) GetDataLength() uint32 {
	return m.Length
}

func (m *Message) GetMsgID() uint32 {
	return m.Id
}

func (m *Message) GetData() []byte {
	return m.Data
}

func (m *Message) SetMsgID(u uint32) {
	m.Id = u
}

func (m *Message) SetData(bytes []byte) {
	m.Data = bytes
}

func (m *Message) SetDataLength(u uint32) {
	m.Length = u
}

func NewMessage(id uint32, data []byte) *Message {
	return &Message{
		Id:     id,
		Length: uint32(len(data)),
		Data:   data,
	}
}
