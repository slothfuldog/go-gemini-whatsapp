package infrastructure

import (
	"context"
	"fmt"
	"gemini-gen-ai/function"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/google/generative-ai-go/genai"
	"github.com/mdp/qrterminal"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
)

func WhatsappGo(client genai.Client, ctxGem context.Context) {
	dbLog := waLog.Stdout("Database", "DEBUG", true)

	container, err := sqlstore.New("sqlite3", "file:examplestore.db?_foreign_keys=on", dbLog)
	if err != nil {
		panic(err)
	}
	whatsmeow.KeepAliveIntervalMax = 90 * time.Second

	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		panic(err)
	}

	clientLog := waLog.Stdout("Client", "INFO", true)
	clients := whatsmeow.NewClient(deviceStore, clientLog)

	clients.AddEventHandler(function.WhatsappHandler(client, clients, ctxGem))

	if clients.Store.ID == nil {
		// No ID stored, new login
		qrChan, _ := clients.GetQRChannel(context.Background())
		err = clients.Connect()
		if err != nil {
			panic(err)
		}
		for evt := range qrChan {
			if evt.Event == "code" {
				// Render the QR code here
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
			} else {
				fmt.Println("Login event:", evt.Event)
			}
		}
	} else {
		// Already logged in, just connect
		err = clients.Connect()
		if err != nil {
			panic(err)
		}
	}

	// Listen to Ctrl+C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	clients.Disconnect()
}
