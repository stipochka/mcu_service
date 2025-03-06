package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	if len(os.Args) < 2 {
		fmt.Println("Provide slavePath to run")
		return
	}

	slavePath := os.Args[1]
	file, err := os.OpenFile(slavePath, os.O_RDWR, 0600)
	if err != nil {
		log.Fatalf("failed to open slavePath: %v", err)
	}

	defer file.Close()

	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("failed to read line: %v", err)
		}
		fmt.Println(line)

		resp := fmt.Sprintf("temperature:%d, humidity:%d\n", rand.Uint64()%200, rand.Uint64()%100)
		_, err = file.Write([]byte(resp))
		if err != nil {
			log.Fatalf("failed to write to master %w", err)
		}

	}
}
