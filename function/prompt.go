package function

import (
	"context"
	"fmt"

	"github.com/google/generative-ai-go/genai"
)

var text string

func WritePrompt(prompt string, ctx context.Context, client genai.Client, models string) string {

	fmt.Println("Write Prompt Start.....")
	defer fmt.Println("Write Prompt End....")
	model := client.GenerativeModel(models) // Adjust based
	prompts := prompt
	IsWorking = true

	text = fmt.Sprintf("###### AI RESPONSE ######## \n\nResponding to \"%s\"\n\n", prompts)

	resp, err := model.GenerateContent(ctx, genai.Text(prompts))
	if err != nil {
		fmt.Println("Error generating content:", err.Error())
		return fmt.Sprintf("ERROR: %s", err.Error())
	}
	printResponse(resp)

	IsWorking = false
	return text
}

func printResponse(resp *genai.GenerateContentResponse) {
	fmt.Println("----")
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				text = fmt.Sprintf("%s%s", text, part)
			}
		}
	}
	fmt.Println("---")
}
