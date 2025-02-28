package network

import (
	"net"
	"strconv"
)

// 测试网络是否连通
func TestNetwork(ip string, port int) bool {
	// 先测试是否能连接到设备
	// 创建一个UDP连接
	conn, err := net.Dial("udp", ip+":"+strconv.Itoa(port))
	if err != nil {
		return false
	}
	defer conn.Close()

	// 发送一个测试包
	message := []byte("test")
	_, err = conn.Write(message)
	if err != nil {
		return false
	}
	conn, err = net.Dial("tcp", ip+":"+strconv.Itoa(port))
	if err != nil {
		return false
	}
	conn.Close()
	return true
}
