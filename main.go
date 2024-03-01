package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"sync"
	"time"
)

func main() {

	start := time.Now()
	var userHash string
	fmt.Print("Please enter hash: ")
	fmt.Scanln(&userHash)

	file, err := os.Open("hash-file.txt")
	if err != nil {
		fmt.Println("An error occurred opening the file.:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var wg sync.WaitGroup

	lines := make(chan string)
	stop := make(chan struct{})

	wg.Add(1)
	go func() {
		defer close(lines)
		defer wg.Done()
		for scanner.Scan() {
			select {
			case lines <- scanner.Text():
			case <-stop:
				return
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for line := range lines {

			hasher := md5.New()
			hasher.Write([]byte(line))
			hashed := hasher.Sum(nil)
			hashedStr := hex.EncodeToString(hashed)

			if hashedStr == userHash {
				fmt.Println("Found matching hash!")
				fmt.Println("Data:", line)
				close(stop)
				return
			}
		}
		fmt.Println("No matching hash found.")
	}()

	wg.Wait()

	if err := scanner.Err(); err != nil {
		fmt.Println("An error occurred reading the file.:", err)
		return
	}
	elapsed := time.Since(start)
	fmt.Println("Total time spent:", elapsed)
}
