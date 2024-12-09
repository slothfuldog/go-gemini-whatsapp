package function

import (
	"context"
	"fmt"
	"os"
	_ "strings"

	"github.com/google/generative-ai-go/genai"
	_ "google.golang.org/api/iterator"
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

/* func StartStreams(prompt string, ctx context.Context, client genai.Client, models string, answer *string) {
	fmt.Println("StartStreams Prompt Start.....")
	defer fmt.Println("StartStreams Prompt End....")
	model := client.GenerativeModel(models) // Adjust based
	prompts := prompt

	text = fmt.Sprintf("###### AI RESPONSE ######## \n\nResponding to \"%s\"\n\n", prompts)

	iter := model.GenerateContentStream(ctx, genai.Text(prompts))

	for {
		resp, err := iter.Next()
		splitted := strings.Split(prompt, " ")
		if splitted[0] == "!endChat" {
			*answer = fmt.Sprintf("###### AI RESPONSE ######## \n\nResponding to \"%s\"\n\n CHAT ENDED", prompts)
			break
		}
		if err == iterator.Done {
			fmt.Println("Iterator Done")
			break
		}
		if err != nil {
			fmt.Println("error,", err)
		}
		printResponse(resp)
		*answer = text
	}

} */

func WriteImgPrompt(prompt string, ctx context.Context, client genai.Client, models string, filepath string) string {

	fmt.Println("WriteImgPrompt Start.....")
	defer fmt.Println("WriteImgPrompt End....")
	model := client.GenerativeModel(models) // Adjust based
	prompts := prompt
	IsWorking = true

	text = fmt.Sprintf("###### AI RESPONSE ######## \n\nResponding to \"%s\"\n\n", prompts)

	imgData, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Println("Error generating content:", err.Error())
		errs := os.Remove(filepath)
		fmt.Println(errs)
		return fmt.Sprintf("ERROR: %s", err.Error())
	}

	resp, err := model.GenerateContent(ctx,
		genai.Text(prompt),
		genai.ImageData("jpeg", imgData))

	if err != nil {
		fmt.Println("Error generating content:", err.Error())
		errs := os.Remove(filepath)
		fmt.Println(errs)
		return fmt.Sprintf("ERROR: %s", err.Error())
	}
	printResponse(resp)

	err = os.Remove(filepath)
	if err != nil {
		fmt.Println("Error removing data:", err.Error())
		return fmt.Sprintf("ERROR: %s", err.Error())
	}

	IsWorking = false
	return text
}

func printResponse(resp *genai.GenerateContentResponse) {
	fmt.Println("----")
	fmt.Println("Generating using", models)
	fmt.Println("========================")
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				text = fmt.Sprintf("%s%s", text, part)
			}
		}
	}
	fmt.Println(text)
	fmt.Println("---")
}
