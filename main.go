package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Response struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Text     string `json:"text"`
		Index    int    `json:"index"`
		Logprobs struct {
			Tokens       []string    `json:"tokens"`
			TopLogprobs  [][]float64 `json:"top_logprobs"`
			TextOffset   []int       `json:"text_offset"`
			FinishReason string      `json:"finish_reason"`
		} `json:"logprobs"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error to load .env")
	}

	inputReader := bufio.NewReader(os.Stdin)
	input, _ := inputReader.ReadString('\n')

	requestBody := fmt.Sprintf(`{
		"model": "text-davinci-003",
		"prompt": "%s",
		"max_tokens": 300
	}`, strings.TrimSpace(input))

	url := "https://api.openai.com/v1/completions"
	req, err := http.NewRequest("POST", url, strings.NewReader(requestBody))

	if err != nil {
		fmt.Println("create request error:", err)
		return
	}
	apikey := os.Getenv("CHATGPT_API_KEY")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+apikey)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("send request error:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("read and parse response error:", err)
		return
	}

	var response Response
	json.Unmarshal(body, &response)
	fmt.Println(response.Choices[0].Text)
}
