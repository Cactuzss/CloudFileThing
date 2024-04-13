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

var workingDir = ""
var cfgFile = "config.json"

var conf = Configurator{}

func errcheck(err error) {
	if err != nil {
		panic(err)
	}
}

func help() {
	fmt.Println("\nUsage:")
	fmt.Println("set-server <host address>")
	fmt.Println("ping")
	fmt.Println("send <filename>")
	fmt.Println("list")

	fmt.Println()
}

func main() {
	//fmt.Println(len(os.Args), os.Args)

	configure()

	if len(os.Args) < 2 {
		help()
		return
	}

	switch os.Args[1] {
	case "set-server":
		if len(os.Args) != 3 {
			help()
			os.Exit(1)
		}
		setServer(os.Args[2])

	case "ping":
		pingServer()

	case "send":
		if len(os.Args) != 3 {
			help()
			os.Exit(1)
		}
		sendFile(os.Args[2])

	case "list":
		listFiles()

	case "get":
		if len(os.Args) != 3 {
			help()
			os.Exit(1)
		}
		getFile(os.Args[2])

	case "delete":
		if len(os.Args) != 3 {
			help()
			os.Exit(1)
		}
		deleteFile(os.Args[2])

	case "help":
		help()

	default:
		fmt.Println("Unknown command. Type 'help' for help.")
	}
}

func configure() {
	user, err := user.Current()
	errcheck(err)

	workingDir = "/home/" + user.Username + "/GFS/"

	conf = LoadFromJson(workingDir + cfgFile)
}

func setServer(arg string) {
	conf.HostAddress = arg
	SaveAsJson(workingDir+cfgFile, conf)
}

func pingServer() {
	conn, err := net.Dial("tcp", conf.HostAddress)
	errcheck(err)

	fmt.Println("Sending ping to " + conf.HostAddress)

	writer := bufio.NewWriter(conn)

	_, err = writer.WriteString(("ping\n"))
	errcheck(err)
	writer.Flush()

	time.Sleep(1 * time.Second)

	data, err := bufio.NewReader(conn).ReadString('\n')
	errcheck(err)

	if len(data) != 0 {
		fmt.Println("Server responded with: " + data)
		return
	}

	panic("No response from server.")
}

func sendFile(filename string) {
	conn, err := net.Dial("tcp", conf.HostAddress)
	errcheck(err)
	defer conn.Close()

	fmt.Println("Sending file " + filename + " to " + conf.HostAddress)

	writer := bufio.NewWriter(conn)
	writer.WriteString("send\n")
	writer.Flush()

	writer.WriteString(filename + "\n")
	writer.Flush()

	writer.WriteString(conf.Username + "\n")
	writer.Flush()

	src, err := os.Open(filename)
	errcheck(err)

	defer src.Close()
	io.Copy(conn, src)
}

func listFiles() {
	conn, err := net.Dial("tcp", conf.HostAddress)
	errcheck(err)
	defer conn.Close()

	writer := bufio.NewWriter(conn)
	writer.WriteString("list\n")
	writer.Flush()

	writer.WriteString(conf.Username + "\n")
	writer.Flush()

	reader := bufio.NewReader(conn)
	data, err := reader.ReadString('\n')
	errcheck(err)

	res := strings.Split(data, "\t")

	for i := 0; i < len(res); i++ {
		fmt.Println(res[i])
	}
}

func getFile(filename string) {
	conn, err := net.Dial("tcp", conf.HostAddress)
	errcheck(err)
	defer conn.Close()

	writer := bufio.NewWriter(conn)
	writer.WriteString("get\n")
	writer.Flush()

	writer.WriteString(filename + "\n")
	writer.Flush()

	writer.WriteString(conf.Username + "\n")
	writer.Flush()

	dst, err := os.Create(workingDir + filename)
	errcheck(err)

	defer dst.Close()
	io.Copy(dst, conn)
}

func deleteFile(filename string) {
	conn, err := net.Dial("tcp", conf.HostAddress)
	errcheck(err)
	defer conn.Close()

	writer := bufio.NewWriter(conn)
	writer.WriteString("delete\n")
	writer.Flush()

	writer.WriteString(filename + "\n")
	writer.Flush()

	writer.WriteString(conf.Username + "\n")
	writer.Flush()

}
