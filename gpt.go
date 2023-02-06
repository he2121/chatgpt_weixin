package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	GptToken = "xxx"
)

func init() {
	http.DefaultClient.Timeout = 4900 * time.Millisecond
}

func DoGPTRequest(d string) string {
	req := newGptRequest(d)
	resp, err := doRequest(req)
	var content string
	if err != nil {
		content = err.Error()
	} else {
		content = "服务器爆满，请稍后重试"
		if len(resp.Choices) > 0 {
			content = strings.TrimSpace(resp.Choices[0].Text)
		}
	}
	return content
}

func newGptRequest(d string) *http.Request {
	p := payload{
		Model:       "text-davinci-003",
		Prompt:      d,
		Temperature: 0,
		MaxTokens:   1000,
	}
	bs, err := json.Marshal(p)
	if err != nil {
		log.Println(err)
	}
	reader := bytes.NewReader(bs)
	request, err := http.NewRequest(http.MethodPost, "https://api.openai.com/v1/completions", reader)
	if err != nil {
		log.Println(err)
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+GptToken)
	return request
}

func doRequest(req *http.Request) (*openAICompletionsResp, error) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	bs, err := io.ReadAll(resp.Body)
	log.Println(string(bs))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var s openAICompletionsResp
	err = json.Unmarshal(bs, &s)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &s, nil
}

type openAICompletionsResp struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Text         string      `json:"text"`
		Index        int         `json:"index"`
		Logprobs     interface{} `json:"logprobs"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

type payload struct {
	Model            string  `json:"model"`
	Prompt           string  `json:"prompt"`
	Temperature      float64 `json:"temperature"`
	MaxTokens        int     `json:"max_tokens"`
	FrequencyPenalty float64 `json:"frequency_penalty"`
	PresencePenalty  float64 `json:"presence_penalty"`
}
