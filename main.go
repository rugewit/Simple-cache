package main

import (
	"bufio"
	"cache_project/cache"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	for {
		reader := bufio.NewReader(os.Stdin)
		inputString, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			break
		}
		inputString = strings.Replace(inputString, "\n", "", -1)
		inputString = strings.Replace(inputString, string(rune(13)), "", -1)

		if inputString == "end" {
			fmt.Println("bye")
			break
		} else if inputString == "start" {
			c := cache.New(5 * time.Minute)

			c.Set("a", "b", 1*time.Minute)

			obj, _ := c.Get("a")

			fmt.Println(obj.(string))
		} else {
			fmt.Println("unknown command")
		}
	}
}
