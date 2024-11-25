package function

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
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
var models = "gemini-1.5-flash-8b"
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
			sender := v.Info.MessageSource.Sender.User
			ruler := os.Getenv("RULER")
			list := os.Getenv("LIST")
			lists := strings.Split(list, ",")
			isUserOk := checkUser(sender)
			if sender == ruler {
				isUserOk = true
			}
			if isTrue := findID(lists, user); isTrue && isUserOk {
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
					if sender != ruler {
						clientW.SendMessage(context.Background(), v.Info.Chat, &waProto.Message{
							Conversation: proto.String("You are not authorized to change the model!"),
						})
						return
					}
					models = "gemini-1.5-flash-8b"
					clientW.SendMessage(context.Background(), v.Info.Chat, &waProto.Message{
						Conversation: proto.String("Model changed into Flash"),
					})
				} else if messageBody == "!usePro" {
					if sender != ruler {
						clientW.SendMessage(context.Background(), v.Info.Chat, &waProto.Message{
							Conversation: proto.String("You are not authorized to change the model!"),
						})
						return
					}
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
					strings := "####Command List####\n!turnon: turn on AI command\n!turnoff: turnoff AI command\n!checkSts: Check AI Status\n!useFlash: Use Flash Model\n!usePro: use Pro Model\n!getResponse: Get response from AI (i.e: !getResponse what is love?)\nNote: User should turnon AI command before using !getResponse"
					clientW.SendMessage(context.Background(), v.Info.Chat, &waProto.Message{
						Conversation: proto.String(strings),
					})
				} else if AITurnedON && messageBody == "!turnon" {
					strings := "AITurnedON already TRUE!"
					clientW.SendMessage(context.Background(), v.Info.Chat, &waProto.Message{
						Conversation: proto.String(strings),
					})
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
				} else if messageBody == "!checkUsers" {

					var data []map[string]interface{}
					var message string

					if sender != ruler {
						clientW.SendMessage(context.Background(), v.Info.Chat, &waProto.Message{
							Conversation: proto.String("You are not authorized to use this command!"),
						})
						return
					}

					content, err := GetData()
					if err != nil {
						clientW.SendMessage(context.Background(), v.Info.Chat, &waProto.Message{
							Conversation: proto.String(fmt.Sprintf("%v", err)),
						})
						return
					}

					if content == "" {
						clientW.SendMessage(context.Background(), v.Info.Chat, &waProto.Message{
							Conversation: proto.String("There are no user in the list!"),
						})
						return
					}

					if err = json.Unmarshal([]byte(content), &data); err != nil {
						clientW.SendMessage(context.Background(), v.Info.Chat, &waProto.Message{
							Conversation: proto.String(fmt.Sprintf("%v", err)),
						})
						return
					} else {
						for i, val := range data {
							message += fmt.Sprintf("%d)\tname: %s\n\tphone: +%s\n\n", i+1, val["name"], val["phone"])
						}
						clientW.SendMessage(context.Background(), v.Info.Chat, &waProto.Message{
							Conversation: proto.String(fmt.Sprintf("###### USER LIST ########\n\n%s", message)),
						})
					}
				} else if len(messageBody) >= 12 {
					if messageBody[:4] == "!add" {
						numbers := strings.Split(messageBody, " ")

						if sender != ruler {
							clientW.SendMessage(context.Background(), v.Info.Chat, &waProto.Message{
								Conversation: proto.String("You are not authorized to use this command!"),
							})
							return
						}

						if len(numbers) != 3 {
							clientW.SendMessage(context.Background(), v.Info.Chat, &waProto.Message{
								Conversation: proto.String("Invalid Format!\n use !add phoneNumber name"),
							})
							return
						}
						if res := checkNumber(numbers[1]); res {
							content, err := GetData()

							if err != nil {
								clientW.SendMessage(context.Background(), v.Info.Chat, &waProto.Message{
									Conversation: proto.String(fmt.Sprintf("%v", err)),
								})
							}

							if isDuplicate, whichFound := checkDuplicate(numbers[1], numbers[2]); isDuplicate {
								clientW.SendMessage(context.Background(), v.Info.Chat, &waProto.Message{
									Conversation: proto.String(fmt.Sprintf("%s ALREADY ADDED", whichFound)),
								})
								return
							}

							var data []map[string]interface{}

							if content == "" {
								data = []map[string]interface{}{
									{"name": numbers[2],
										"phone": numbers[1],
									},
								}
							} else {
								err := json.Unmarshal([]byte(content), &data)
								if err != nil {
									clientW.SendMessage(context.Background(), v.Info.Chat, &waProto.Message{
										Conversation: proto.String(fmt.Sprintf("%v", err)),
									})
								}
								currentData := map[string]interface{}{
									"name":  numbers[2],
									"phone": numbers[1],
								}
								data = append(data, currentData)
							}

							json, err := json.Marshal(data)

							if err != nil {
								clientW.SendMessage(context.Background(), v.Info.Chat, &waProto.Message{
									Conversation: proto.String(fmt.Sprintf("%s - %s cannot be added: [%v]", numbers[2], numbers[1], err)),
								})
								return
							}

							CreateFile(string(json))

							clientW.SendMessage(context.Background(), v.Info.Chat, &waProto.Message{
								Conversation: proto.String(fmt.Sprintf("%s - %s has been added!", numbers[2], numbers[1])),
							})
						} else {
							fmt.Printf("is null %t\n", res)
							clientW.SendMessage(context.Background(), v.Info.Chat, &waProto.Message{
								Conversation: proto.String(fmt.Sprintf("%s is not a number", numbers[1])),
							})
						}
					} else if messageBody[:7] == "!remove" {
						body := strings.Split(messageBody, " ")
						if sender != ruler {
							clientW.SendMessage(context.Background(), v.Info.Chat, &waProto.Message{
								Conversation: proto.String("You are not authorized to use this command!"),
							})
							return
						}
						err := removeUser(body[1])
						if err != nil {
							clientW.SendMessage(context.Background(), v.Info.Chat, &waProto.Message{
								Conversation: proto.String(fmt.Sprintf("%v", err)),
							})
							return
						}
						clientW.SendMessage(context.Background(), v.Info.Chat, &waProto.Message{
							Conversation: proto.String(fmt.Sprintf("%s is successfully removed!", body[1])),
						})

					} else if AITurnedON && len(messageBody) == 12 && messageBody[:12] == "!getResponse" {
						clientW.SendMessage(context.Background(), v.Info.Chat, &waProto.Message{
							Conversation: proto.String("Write something!"),
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

func checkNumber(str string) bool {

	re := regexp.MustCompile(`^[0-9]+$`)

	return re.MatchString(str)
}

func checkUser(str string) bool {
	var data []map[string]interface{}
	content, err := GetData()
	if err != nil {
		fmt.Println("checkUser: ", err)
	}

	err = json.Unmarshal([]byte(content), &data)

	if err != nil {
		fmt.Println("checkUser: ", err)
	}

	for _, val := range data {
		if val["phone"] == str {
			return true
		}
	}

	return false
}

func checkDuplicate(phone string, name string) (isDuplicate bool, whichFound string) {
	var data []map[string]interface{}
	content, err := GetData()
	if err != nil {
		fmt.Println("checkDuplicate: ", err)
	}

	err = json.Unmarshal([]byte(content), &data)

	if err != nil {
		fmt.Println("checkDuplicate: ", err)
	}

	for _, val := range data {
		if val["phone"] == phone {
			return true, phone
		} else if val["name"] == name {
			return true, name
		}
	}

	return false, ""
}

func removeUser(str string) error {
	var data []map[string]interface{}
	dir, _ := os.Getwd()
	currDir := fmt.Sprintf("%s/.env", dir)
	err := godotenv.Load(currDir)
	if err != nil {
		fmt.Println("checkUser: ", err)
	}
	idx := -1
	ruler := os.Getenv("RULER")
	content, err := GetData()
	if err != nil {
		fmt.Println("checkUser: ", err)
	}

	err = json.Unmarshal([]byte(content), &data)

	if err != nil {
		fmt.Println("removeUser", err)
		return err
	}

	for i, val := range data {
		if val["name"] == str {
			if val["phone"] == ruler {
				return fmt.Errorf("what are you doing?")
			}
			idx = i
		}
	}

	if idx == -1 {
		return fmt.Errorf("user not found")
	}

	data = append(data[:idx], data[idx+1:]...)

	result, err := json.Marshal(data)
	if err != nil {
		fmt.Println("removeUser", err)
		return err
	}

	CreateFile(string(result))

	return nil
}
