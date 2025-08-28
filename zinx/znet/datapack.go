package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"zinx/utils"
	"zinx/ziface"
)

type DataPack struct {
}

func (d *DataPack) GetHeadLen() uint32 {
	return 8
}

func (d *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	dataBuff := bytes.NewBuffer([]byte{})
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetDataLength()); err != nil {
		return nil, err
	}
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgID()); err != nil {
		return nil, err
	}
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}
	return dataBuff.Bytes(), nil
}

func (d *DataPack) Unpack(b []byte) (ziface.IMessage, error) {
	dataBuff := bytes.NewReader(b)
	msg := &Message{}
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Length); err != nil {
		return nil, err
	}
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}
	/*if err := binary.Read(dataBuff, binary.LittleEndian, msg.Data); err != nil {
		return nil, err
	}*/
	if utils.GlobalObject.MaxPacketSize > 0 && msg.Length > utils.GlobalObject.MaxPacketSize {
		return nil, errors.New("too large msg data received")
	}
	return msg, nil
}

func NewDataPack() *DataPack {
	return &DataPack{}
}
