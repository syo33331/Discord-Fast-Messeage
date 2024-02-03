package main

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"fmt"
	"io/ioutil"
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

	fmt.Print("Add random 3-character alphanumeric to each message? (yes/no): ")
	addRandomStr, _ := reader.ReadString('\n')
	addRandomStr = strings.TrimSpace(addRandomStr)

	fmt.Print("Enter the number of times to send the message: ")
	timesStr, _ := reader.ReadString('\n')
	timesStr = strings.TrimSpace(timesStr)
	times, err := strconv.Atoi(timesStr)
	if err != nil {
		fmt.Println("Invalid number")
		return
	}

	tokens, err := readTokens("token.txt")
	if err != nil {
		fmt.Println("Error reading tokens:", err)
		return
	}

	var wg sync.WaitGroup
	tokenCount := len(tokens)
	for i := 0; i < times; i++ {
		modifiedMessage := message
		if strings.ToLower(addRandomStr) == "yes" {
			randomStr := generateRandomString(3)
			modifiedMessage += randomStr // ランダムな文字列をメッセージの最後に追加
		}
		wg.Add(1)
		go sendMessage(channelID, modifiedMessage, tokens[i%tokenCount], i+1, &wg)
	}

	wg.Wait()
	fmt.Println("All messages sent")
}

func readTokens(filePath string) ([]string, error) {
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	tokens := strings.Split(strings.TrimSpace(string(fileContent)), "\n")
	return tokens, nil
}

func sendMessage(channelID, message, token string, count int, wg *sync.WaitGroup) {
	defer wg.Done()

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

	fmt.Printf("Message %d sent with token: StatusCode = %d\n", count, resp.StatusCode)
}

// generateRandomString は指定された長さのランダムな英数字の文字列を生成します。
func generateRandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bytes := make([]byte, n)
	_, err := rand.Read(bytes)
	if err != nil {
		fmt.Println("Error generating random string:", err)
		return ""
	}
	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}
	return string(bytes)
}
