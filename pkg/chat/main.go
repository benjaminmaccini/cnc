package chat

import (
	"context"

	openai "github.com/sashabaranov/go-openai"
	log "github.com/sirupsen/logrus"
)

type ChatClient struct {
	Client *openai.Client
}

func InitClient(key string) *ChatClient {
	return &ChatClient{
		Client: openai.NewClient(key),
	}
}

func (client *ChatClient) Query(msg string) string {
	resp, err := client.Client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: msg,
				},
			},
		},
	)

	if err != nil {
		log.Fatal(err)
	}

	return resp.Choices[0].Message.Content
}
