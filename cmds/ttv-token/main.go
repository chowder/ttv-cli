package main

import (
	"fmt"
	"log"
	"os"
	"ttv-cli/internals/pkg/twitch/login"
)

func main() {
	authToken, err := login.GetAccessToken(os.Getenv("TWITCH_USERNAME"), os.Getenv("TWITCH_PASSWORD"))
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(authToken)
}
