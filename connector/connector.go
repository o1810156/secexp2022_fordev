package connector

import (
	"fmt"
	"net"
	"time"
)

const PORT = ":10000"
const PROTOCOL = "tcp"
const DEADLINE = time.Duration(10) * time.Second

type Communicator interface {
	GetConn() net.Conn
}

func raw_send(com Communicator, data []byte) error {
	conn := com.GetConn()
	conn.Write(data)

	return nil
}

func send(com Communicator, data []byte) error {
	err := raw_send(com, data)
	if err != nil {
		return err
	}

	_, err = raw_receive(com, []byte{})

	return err
}

func raw_receive(com Communicator, buffer []byte) (int, error) {
	conn := com.GetConn()
	len, err := conn.Read(buffer)

	return len, err
}

func receive(com Communicator, buffer []byte) (int, error) {
	len, err := raw_receive(com, buffer)
	if err != nil {
		return 0, err
	}

	err = send(com, []byte{})
	if err != nil {
		return 0, err
	}

	return len, nil
}

// 予備校に該当するよい英単語がないためローマ字表記
type Yobikou struct {
	conn net.Conn
}

func NewYobikouServer() (Yobikou, error) {
	tcpAddr, err := net.ResolveTCPAddr(PROTOCOL, PORT)
	if err != nil {
		return Yobikou{}, err
	}

	listener, err := net.ListenTCP(PROTOCOL, tcpAddr)
	if err != nil {
		return Yobikou{}, err
	}

	/*
		err = listener.SetDeadline(time.Now().Add(DEADLINE))
		if err != nil {
			return Yobikou{}, err
		}
	*/

	fmt.Println("中学側に接続します……")

	conn, err := listener.Accept()
	if err != nil {
		return Yobikou{}, err
	}

	return Yobikou{conn}, nil
}

func (yobikou Yobikou) GetConn() net.Conn {
	return yobikou.conn
}

func (yobikou Yobikou) Close() error {
	return yobikou.conn.Close()
}

func (yobikou Yobikou) Send(data []byte) error {
	return send(yobikou, data)
}

func (yobikou Yobikou) Receive(buffer []byte) (int, error) {
	return receive(yobikou, buffer)
}

// 予備校に合わせ中学校もローマ字で
type Tyuugaku struct {
	conn net.Conn
}

func NewTyuugakuClient(ServerAddr string) (Tyuugaku, error) {
	fmt.Println("予備校に接続します……")

	conn, err := net.DialTimeout(PROTOCOL, ServerAddr+PORT, DEADLINE)
	if err != nil {
		return Tyuugaku{}, err
	}

	return Tyuugaku{conn}, nil
}

func (tyuugaku Tyuugaku) GetConn() net.Conn {
	return tyuugaku.conn
}

func (tyuugaku Tyuugaku) Close() error {
	return tyuugaku.conn.Close()
}

func (tyuugaku Tyuugaku) Send(data []byte) error {
	return send(tyuugaku, data)
}

func (tyuugaku Tyuugaku) Receive(buffer []byte) (int, error) {
	return receive(tyuugaku, buffer)
}
