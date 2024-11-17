package function

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

var AITurnedON bool
var models = "gemini-1.5-flash"
var IsWorking = false

func WhatsappHandler(client genai.Client, clientW *whatsmeow.Client, ctx context.Context) func(interface{}) {
	defer fmt.Println("Whatsapp Handler End....")
	fmt.Println("Whatsapp Handler Start.......")

	dir, _ := os.Getwd()
	currDir := fmt.Sprintf("%s/.env", dir)
	err := godotenv.Load(currDir)

	if err != nil {
		fmt.Println("ERROR: ", err)
	}
	done := make(chan bool)

	return func(evt interface{}) {
		switch v := evt.(type) {
		case *events.Message:
			var messageBody = v.Message.GetConversation()
			fmt.Println(v.Info.MessageSource.Chat.User)
			user := v.Info.MessageSource.Chat.User
			list := os.Getenv("LIST")
			lists := strings.Split(list, ",")
			if isTrue := findID(lists, user); isTrue {
				fmt.Println("Message Body: ", messageBody)
				fmt.Println("Extended Message: ", v.Message.GetExtendedTextMessage().GetText())
				if messageBody == "" {
					messageBody = v.Message.GetExtendedTextMessage().GetText()
				}
				if messageBody == "!turnon" && !AITurnedON {
					AITurnedON = true
					strings := fmt.Sprintf("AITurnedON Status: %t", AITurnedON)
					clientW.SendMessage(context.Background(), v.Info.Chat, &waProto.Message{
						Conversation: proto.String(strings),
					})
				} else if messageBody == "!checkSts" {
					strings := fmt.Sprintf("AITurnedON status: %t", AITurnedON)
					clientW.SendMessage(context.Background(), v.Info.Chat, &waProto.Message{
						Conversation: proto.String(strings),
					})
				} else if messageBody == "!useFlash" {
					models = "gemini-1.5-flash"
					clientW.SendMessage(context.Background(), v.Info.Chat, &waProto.Message{
						Conversation: proto.String("Model changed into Flash"),
					})
				} else if messageBody == "!usePro" {
					models = "gemini-1.5-pro"
					clientW.SendMessage(context.Background(), v.Info.Chat, &waProto.Message{
						Conversation: proto.String("Model changed into Pro"),
					})
				} else if messageBody == "!checkModel" {
					strings := fmt.Sprintf("Current model: %s", models)
					clientW.SendMessage(context.Background(), v.Info.Chat, &waProto.Message{
						Conversation: proto.String(strings),
					})
				} else if messageBody == "!checkList" {
					strings := fmt.Sprintf("####Command List####\n!turnon: turn on AI command\n!turnoff: turnoff AI command\n!checkSts: Check AI Status\n!useFlash: Use Flash Model\n!usePro: use Pro Model\n!getResponse: Get response from AI (i.e: !getResponse what is love?)\nNote: User should turnon AI command before using !getResponse")
					clientW.SendMessage(context.Background(), v.Info.Chat, &waProto.Message{
						Conversation: proto.String(strings),
					})
				} else if AITurnedON && messageBody == "!turnon" {
					strings := "AITurnedON already TRUE!"
					clientW.SendMessage(context.Background(), v.Info.Chat, &waProto.Message{
						Conversation: proto.String(strings),
					})
				} else if len(messageBody) >= 12 {
					if AITurnedON && len(messageBody) == 12 && messageBody[:12] == "!getResponse" {
						clientW.SendMessage(context.Background(), v.Info.Chat, &waProto.Message{
							Conversation: proto.String("tulis sesuatu!"),
						})
					} else if AITurnedON && messageBody[:12] == "!getResponse" {
						prompt := messageBody[12:]
						tries := 0
						for IsWorking {
							clientW.SendMessage(context.Background(), v.Info.Chat, &waProto.Message{
								Conversation: proto.String(fmt.Sprintf("Waiting job to be done ...(%d secs)", tries*5)),
							})
							time.Sleep(5 * time.Second)
							tries++
						}
						go func() {
							clientW.SendMessage(context.Background(), v.Info.Chat, &waProto.Message{
								Conversation: proto.String("Generating data...."),
							})
							resp := WritePrompt(prompt, ctx, client, models)
							clientW.SendMessage(context.Background(), v.Info.Chat, &waProto.Message{
								Conversation: proto.String(resp),
							})
							done <- true
							tries = 0
						}()

						go func() {
							// Monitor elapsed time
							interval := 10 * time.Second
							start := time.Now()
							ticker := time.NewTicker(interval)
							for {
								select {
								case <-ticker.C:
									elapsed := time.Since(start)
									fmt.Printf("Elapsed time: %v\n", elapsed)
									stringss := fmt.Sprintf("Elapsed time: %v\n", elapsed)
									clientW.SendMessage(context.Background(), v.Info.Chat, &waProto.Message{
										Conversation: proto.String(stringss),
									})
								case <-done:
									ticker.Stop()
									return
								}
							}
						}()

					} else if !AITurnedON && messageBody[:12] == "!getResponse" {
						clientW.SendMessage(context.Background(), v.Info.Chat, &waProto.Message{
							Conversation: proto.String("Turnon AI first!"),
						})
					}
				} else if messageBody == "!turnoff" {
					if !AITurnedON {
						clientW.SendMessage(context.Background(), v.Info.Chat, &waProto.Message{
							Conversation: proto.String("AITurnedON already false!!!"),
						})
					} else {
						AITurnedON = false
						strings := fmt.Sprintf("AITurnedON status: %t", AITurnedON)
						clientW.SendMessage(context.Background(), v.Info.Chat, &waProto.Message{
							Conversation: proto.String(strings),
						})
					}
				}
			}

		}
	}
}

func findID(str []string, match string) bool {
	for _, val := range str {
		if val == match {
			return true
		}
	}
	return false
}
