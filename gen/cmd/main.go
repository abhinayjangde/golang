package main

import (
	"context"
	"fmt"
	"os"

	"github.com/abhinayjangde/sarvam"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	client := sarvam.NewClient(os.Getenv("SARVAM_API_KEY"))
	var messages []sarvam.Message

	for true {
		var query string
		fmt.Print("> ")
		fmt.Scanln(&query)

		if query == "exit" {
			break
		}
		messages = append(messages, sarvam.Message{
			Role:    "user",
			Content: query,
		})
		resp, err := client.Chat.Completions.Create(
			context.Background(),
			sarvam.ChatCompletionRequest{
				Model:    sarvam.ModelSarvam105BConversations,
				Messages: messages,
			},
		)

		if err != nil {
			panic(err)
		}

		messages = append(messages, sarvam.Message{
			Role:    "assistant",
			Content: resp.Choices[0].Message.Content,
		})

		fmt.Println("> " + resp.Choices[0].Message.Content)
	}

}
