package infrastructure

func Env() {
	clients, ctx := GeminiGo()
	WhatsappGo(*clients, ctx)
}
