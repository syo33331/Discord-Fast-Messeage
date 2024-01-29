package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter Channel ID: ")
	channelID, _ := reader.ReadString('\n')
	channelID = strings.TrimSpace(channelID)

	fmt.Print("Enter the message: ")
	message, _ := reader.ReadString('\n')
	message = strings.TrimSpace(message)

	fmt.Print("Enter the number of times to send the message: ")
	timesStr, _ := reader.ReadString('\n')
	timesStr = strings.TrimSpace(timesStr)
	times, err := strconv.Atoi(timesStr)
	if err != nil {
		fmt.Println("Invalid number")
		return
	}

	token := "Token"

	var wg sync.WaitGroup
	for i := 0; i < times; i++ {
		wg.Add(1)
		go func(count int) {
			defer wg.Done()
			sendMessage(channelID, message, token, count)
		}(i + 1)
	}

	wg.Wait()
	fmt.Println("All messages sent")
}

func sendMessage(channelID, message, token string, count int) {
	url := fmt.Sprintf("https://discord.com/api/v9/channels/%s/messages", channelID)
	content := fmt.Sprintf(`{"content": "%s"}`, message)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(content)))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("Message %d sent: StatusCode = %d\n", count, resp.StatusCode)
}
