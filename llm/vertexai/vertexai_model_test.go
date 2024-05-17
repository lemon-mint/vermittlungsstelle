package vertexai_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/lemon-mint/vermittlungsstelle/llm"
	"github.com/lemon-mint/vermittlungsstelle/llm/vertexai"

	"cloud.google.com/go/vertexai/genai"
	"github.com/lemon-mint/godotenv"
)

var client *genai.Client = func() *genai.Client {
	envfile, err := os.ReadFile("../../.env")
	if err != nil {
		panic(err)
	}
	for k, v := range godotenv.Parse(string(envfile)) {
		os.Setenv(k, v)
	}

	client, err := vertexai.NewClient(
		context.Background(),
		os.Getenv("PROJECT_ID"),
		os.Getenv("REGION"),
	)
	if err != nil {
		panic(err)
	}
	return client
}()

func TestVertexAIGenerate(t *testing.T) {
	var model llm.LLM = vertexai.NewModel(client, "gemini-pro", nil)
	defer model.Close()

	output := model.GenerateStream(
		context.Background(),
		&llm.ChatContext{},
		&llm.Content{
			Role:  llm.RoleUser,
			Parts: []llm.Segment{llm.Text("Hello!")},
		},
	)

	for segment := range output.Stream {
		fmt.Print(segment)
	}
	fmt.Println()

	if output.Err != nil {
		t.Error(output.Err)
		return
	}

	if output.UsageData.OutputTokens <= 0 {
		t.Errorf("expected output tokens > 0, got %d\n", output.UsageData.OutputTokens)
		return
	}
}
