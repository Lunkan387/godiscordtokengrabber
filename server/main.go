package main

import (
	"crypto/md5"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

var routePath, routePath2, routePath3, routePath4, routePath5 string
var mu sync.Mutex

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./server <salt>")
		return
	}

	salt, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("Invalid salt value")
		return
	}

	handle(salt)
}

func handle(salty int) {
	app := fiber.New()

	go func() {
		for {
			mu.Lock()
			routePath, routePath2, routePath3, routePath4, routePath5 = GenerateDest(salty)
			mu.Unlock()
			time.Sleep(1 * time.Second)
		}
	}()

	app.All("/*", func(c *fiber.Ctx) error {
		mu.Lock()
		currentRoute := routePath
		currentRoute2 := routePath2
		currentRoute3 := routePath3
		currentRoute4 := routePath4
		currentRoute5 := routePath5
		mu.Unlock()

		if c.Path() == currentRoute || c.Path() == currentRoute2 || c.Path() == currentRoute3 || c.Path() == currentRoute4 || c.Path() == currentRoute5 {
			body := c.Body()
			fmt.Println(string(body))
			return c.SendString("Data received")
		}
		return c.SendStatus(fiber.StatusNotFound)
	})

	log.Fatal(app.Listen(":4444"))
}

func GenerateDest(salty int) (string, string, string, string, string) {
	location, _ := time.LoadLocation("Europe/Paris")
	nowtime := time.Now().In(location)

	nowsecond := nowtime.Second() + salty
	nowsecond2 := nowtime.Second() + -1 + salty
	nowsecond3 := nowtime.Second() + -2 + salty
	nowsecond4 := nowtime.Second() + -3 + salty
	nowsecond5 := nowtime.Second() + -4 + salty

	hash := md5.Sum([]byte(fmt.Sprintf("%d", nowsecond)))
	hashString := fmt.Sprintf("/%x", hash)

	hash2 := md5.Sum([]byte(fmt.Sprintf("%d", nowsecond2)))
	hashString2 := fmt.Sprintf("/%x", hash2)

	hash3 := md5.Sum([]byte(fmt.Sprintf("%d", nowsecond3)))
	hashString3 := fmt.Sprintf("/%x", hash3)

	hash4 := md5.Sum([]byte(fmt.Sprintf("%d", nowsecond4)))
	hashString4 := fmt.Sprintf("/%x", hash4)

	hash5 := md5.Sum([]byte(fmt.Sprintf("%d", nowsecond5)))
	hashString5 := fmt.Sprintf("/%x", hash5)

	return hashString, hashString2, hashString3, hashString4, hashString5
}
