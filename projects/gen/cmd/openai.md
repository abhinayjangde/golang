package main

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
)

func main() {
	godotenv.Load()
	client := openai.NewClient(
		option.WithAPIKey(os.Getenv("GROQ_API_KEY")),
		option.WithBaseURL("https://api.groq.com/openai/v1"),
	)

	response, err := client.Responses.New(context.Background(), responses.ResponseNewParams{
		Model: "groq/compound",
		Input: responses.ResponseNewParamsInputUnion{
			OfString: openai.String("hello world"),
		},
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(response.OutputText())
}
