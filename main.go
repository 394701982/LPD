package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	ServerName    string `json:"server_name"`
	BannerMessage string `json:"banner_message"`
	Version       string `json:"version"`
	Port          int    `json:"port"`
	LogFilePath   string `json:"log_file_path"`
}

func LoadConfig(filename string) (*Config, error) {
	absPath, err := filepath.Abs(filename)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(absPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := Config{}
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func setupLogging(logFilePath string) (*os.File, error) {
	absPath, err := filepath.Abs(logFilePath)
	if err != nil {
		return nil, err
	}

	logFile, err := os.OpenFile(absPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime)
	return logFile, nil
}

func handleConnection(conn net.Conn, config *Config) {
	defer conn.Close()

	clientAddr := conn.RemoteAddr().String()
	log.Printf("[LPD] Connection from %s\n", clientAddr)

	reader := bufio.NewReader(conn)

	for {
		command, err := reader.ReadByte()
		if err != nil {
			if err != io.EOF {
				log.Printf("[LPD] Error reading from %s: %v\n", clientAddr, err)
			}
			return
		}

		switch command {
		case 0x02:
			// 打印机状态请求或打印作业提交
			line, _ := reader.ReadString('\n')
			line = strings.TrimSpace(line)
			log.Printf("[LPD] Received command: 0x02 %s from %s\n", line, clientAddr)

			if line == "print_queue" {
				// 响应特定的打印机状态请求
				response := fmt.Sprintf("%s [@Boyk]: Print-services are not available to your host (aF2qXkQ2m).\n", config.ServerName)
				conn.Write([]byte(response))
				log.Printf("[LPD] Sent custom status to %s: %s\n", clientAddr, response)
			} else {
				// 模拟打印作业提交
				queueName := line
				fileContent, _ := reader.ReadString('\n')
				log.Printf("[LPD] Received print job for queue: %s from %s\n", queueName, clientAddr)
				log.Printf("[LPD] Received file content for queue %s: %s\n", queueName, fileContent)
				jobConfirmation := fmt.Sprintf("%s [%s]: Print job for queue %s received\n", config.ServerName, config.Version, queueName)
				conn.Write([]byte(jobConfirmation))
				log.Printf("[LPD] Sent job confirmation to %s: %s\n", clientAddr, jobConfirmation)
			}

		case 0x03:
			// 打印队列列表请求
			log.Printf("[LPD] Received queue list request from %s\n", clientAddr)
			queueListResponse := fmt.Sprintf("%s [%s]: Available queues: queue1, queue2, queue3\n", config.ServerName, config.Version)
			conn.Write([]byte(queueListResponse))
			log.Printf("[LPD] Sent queue list to %s: %s\n", clientAddr, queueListResponse)

		default:
			log.Printf("[LPD] Received unknown command from %s: 0x%x\n", clientAddr, command)
		}
	}
}

func startServer(config *Config) {
	addr := fmt.Sprintf(":%d", config.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Error starting server: %v\n", err)
	}
	defer listener.Close()

	log.Printf("[LPD] Server started on port %d\n", config.Port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v\n", err)
			continue
		}

		go handleConnection(conn, config)
	}
}

func main() {
	configPath := flag.String("config", "config.json", "Path to configuration file")
	port := flag.Int("port", 0, "Port to run the server on (overrides config file)")

	flag.Parse()

	config, err := LoadConfig(*configPath)
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	if *port != 0 {
		config.Port = *port
	}

	logFile, err := setupLogging(config.LogFilePath)
	if err != nil {
		fmt.Printf("Error setting up logging: %v\n", err)
		return
	}
	defer logFile.Close()

	startServer(config)
}
