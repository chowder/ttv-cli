package main

import (
	"fmt"
	"os"
	"ttv-cli/internals/pkg/twitch/login"
)

func main() {
	resp := login.GetAccessToken(os.Getenv("TWITCH_USERNAME"), os.Getenv("TWITCH_PASSWORD"))
	fmt.Println(resp)
}
