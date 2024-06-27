package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Config holds the server configuration
type Config struct {
	Port        string `json:"port"`
	Password    string `json:"password"`
	LogFile     string `json:"logfile"`
	AuthMethods string `json:"auth_methods"`
	Version     string `json:"version"`
}

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "", "Path to the configuration file")
	flag.Parse()

	var config Config
	if configPath != "" {
		configData, err := ioutil.ReadFile(configPath)
		if err != nil {
			log.Fatalf("Error reading config file: %v", err)
		}
		if err := json.Unmarshal(configData, &config); err != nil {
			log.Fatalf("Error parsing config file: %v", err)
		}
	} else {
		config.Port = "9051"
		config.Password = "your_password"
		config.LogFile = "tor_control.log"
		config.AuthMethods = "HASHEDPASSWORD"
		config.Version = "0.4.6.5"
	}

	// 设置日志输出到文件
	file, err := os.OpenFile(config.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	log.SetOutput(file)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	listener, err := net.Listen("tcp", ":"+config.Port)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer listener.Close()
	log.Infof("[Tor-Control] Server started on port %s", config.Port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Errorf("Error accepting connection: %v", err)
			continue
		}
		go handleConnection(conn, config)
	}
}

func handleConnection(conn net.Conn, config Config) {
	defer conn.Close()
	log.Infof("[Tor-Control] Client connected: %s", conn.RemoteAddr().String())
	reader := bufio.NewReader(conn)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Errorf("[Tor-Control] Error reading message: %v", err)
			return
		}
		message = strings.TrimSpace(message)
		log.Infof("[Tor-Control] Received: %s", message)

		// 将命令转换为大写以确保不区分大小写
		command := strings.ToUpper(message)

		if strings.HasPrefix(command, "AUTHENTICATE") {
			handleAuthenticate(conn, message, config.Password)
		} else if strings.HasPrefix(command, "PROTOCOLINFO") {
			handleProtocolInfo(conn, config)
		} else if strings.HasPrefix(command, "GETINFO") {
			handleGetInfo(conn, message, config)
		} else if strings.HasPrefix(command, "SIGNAL") {
			handleSignal(conn, message)
		} else if strings.HasPrefix(command, "MAPADDRESS") {
			handleMapAddress(conn, message)
		} else if strings.HasPrefix(command, "GETCONF") {
			handleGetConf(conn, message)
		} else if strings.HasPrefix(command, "SETCONF") {
			handleSetConf(conn, message)
		} else if strings.HasPrefix(command, "QUIT") {
			handleQuit(conn)
			return
		} else {
			conn.Write([]byte("552 Unrecognized command\r\n"))
		}
	}
}

func handleAuthenticate(conn net.Conn, message, password string) {
	parts := strings.Split(message, " ")
	if len(parts) == 2 && parts[1] == fmt.Sprintf("\"%s\"", password) {
		conn.Write([]byte("250 OK\r\n"))
		log.Infof("[Tor-Control] Authentication successful")
	} else {
		conn.Write([]byte("515 Authentication failed\r\n"))
		log.Infof("[Tor-Control] Authentication failed")
	}
}

func handleProtocolInfo(conn net.Conn, config Config) {
	response := fmt.Sprintf("250-PROTOCOLINFO 1\r\n250-AUTH METHODS=%s\r\n250-VERSION Tor=\"%s\"\r\n250 OK\r\n", config.AuthMethods, config.Version)
	conn.Write([]byte(response))
	log.Infof("[Tor-Control] Sent protocol info")
}

func handleGetInfo(conn net.Conn, message string, config Config) {
	if strings.HasSuffix(strings.ToUpper(message), "VERSION") {
		response := fmt.Sprintf("250-version=%s\r\n250 OK\r\n", config.Version)
		conn.Write([]byte(response))
		log.Infof("[Tor-Control] Sent version info")
	} else {
		conn.Write([]byte("552 Unrecognized command\r\n"))
		log.Infof("[Tor-Control] Unrecognized GETINFO command")
	}
}

func handleSignal(conn net.Conn, message string) {
	parts := strings.Split(message, " ")
	if len(parts) == 2 && strings.ToUpper(parts[1]) == "RELOAD" {
		conn.Write([]byte("250 OK\r\n"))
		log.Infof("[Tor-Control] Signal RELOAD received")
	} else {
		conn.Write([]byte("552 Unrecognized command\r\n"))
		log.Infof("[Tor-Control] Unrecognized SIGNAL command")
	}
}

func handleMapAddress(conn net.Conn, message string) {
	conn.Write([]byte("250 OK\r\n"))
	log.Infof("[Tor-Control] MAPADDRESS command received")
}

func handleGetConf(conn net.Conn, message string) {
	conn.Write([]byte("250 OK\r\n"))
	log.Infof("[Tor-Control] GETCONF command received")
}

func handleSetConf(conn net.Conn, message string) {
	conn.Write([]byte("250 OK\r\n"))
	log.Infof("[Tor-Control] SETCONF command received")
}

func handleQuit(conn net.Conn) {
	conn.Write([]byte("250 OK\r\n"))
	log.Infof("[Tor-Control] Client disconnected")
}
