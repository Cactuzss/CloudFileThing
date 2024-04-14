package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"os/user"
	"strings"
	"time"
)

var server_workingDir = ""

//var cfgFile = "server_config.json"

var listen_ip = "localhost"
var listen_port = "1109"

func Server() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	server_workingDir = "/home/" + user.Username + "/GFS/Server/"

	if _, err := os.Stat(server_workingDir); os.IsNotExist(err) {
		os.Mkdir(server_workingDir, 0777)
	}

	listener, err := net.Listen("tcp", listen_ip+":"+listen_port)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}

		go handle_connection(conn)
	}
}

func handle_connection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	data, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	switch data {
	case "ping\n":
		pinged(conn)
		return
	case "send\n":
		gettingFile(conn, reader)
		return
	case "list\n":
		list(conn, reader)
		return
	case "get\n":
		sendingFile(conn, reader)
		return
	case "delete\n":
		deletingFile(conn, reader)
		return
	}

}

func pinged(conn net.Conn) {
	writer := bufio.NewWriter(conn)

	_, err := writer.WriteString("meow\n")
	if err != nil {
		panic(err)
	}

	writer.Flush()
	conn.Close()
}

func gettingFile(conn net.Conn, reader *bufio.Reader) {
	defer conn.Close()

	time.Sleep(1 * time.Second)

	data, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	username, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	username = username[:len(username)-1]

	i := strings.LastIndex(data, "/")
	data = data[i+1 : len(data)-1]

	fmt.Println(server_workingDir + username + "/" + data)

	dst, err := os.Create(server_workingDir + username + "/" + data)
	if err != nil {
		err = os.Mkdir(server_workingDir+username+"/", os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	defer dst.Close()
	io.Copy(dst, conn)
}

func list(conn net.Conn, reader *bufio.Reader) {
	defer conn.Close()

	username, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	username = username[:len(username)-1]

	data, err := os.ReadDir(server_workingDir + username + "/")
	if err != nil {
		panic(err)
	}

	res := ""
	for _, v := range data {
		res += v.Name() + "\t"
	}

	writer := bufio.NewWriter(conn)

	_, err = writer.WriteString(res + "\n")
	if err != nil {
		panic(err)
	}

	writer.Flush()
	conn.Close()
}

func sendingFile(conn net.Conn, reader *bufio.Reader) {
	defer conn.Close()

	filename, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	filename = filename[:len(filename)-1]

	username, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	username = username[:len(username)-1]

	src, err := os.Open(server_workingDir + username + "/" + filename)
	if err != nil {
		panic(err)
	}

	defer src.Close()
	io.Copy(conn, src)
}

func deletingFile(conn net.Conn, reader *bufio.Reader) {
	defer conn.Close()

	filename, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	filename = filename[:len(filename)-1]

	username, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	username = username[:len(username)-1]

	err = os.Remove(server_workingDir + username + "/" + filename)
	if err != nil {
		panic(err)
	}

}
