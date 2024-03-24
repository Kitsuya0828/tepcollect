package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/sclevine/agouti"
	"github.com/slack-go/slack"
)

func logout(page *agouti.Page) {
	log.Printf("info: click logout button")
	if err := page.FindByClass("logout").Click(); err != nil {
		log.Fatal(err)
	}
	time.Sleep(3 * time.Second)
}

func main() {
	driver := agouti.ChromeDriver(
		agouti.ChromeOptions(
			"args", []string{
				"--no-sandbox",
				"--disable-dev-shm-usage",
				"--headless",
				"--disable-gpu",
				"lang=ja",
				"--disable-desktop-notifications",
				"--disable-blink-features=AutomationControlled",
				"--ignore-certificate-errors",
				"--disable-extensions",
				"--window-size=1920,1080",
				"--user-agent=Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Mobile Safari/537.36",
			}),
	)
	driver.Start()
	defer driver.Stop()

	slackToken := os.Getenv("SLACK_TOKEN")
	slackChannelID := os.Getenv("SLACK_CHANNEL_ID")
	fmt.Println("slackToken: ", slackToken, "slackChannelID: ", slackChannelID)
	api := slack.New(slackToken)

	page, _ := driver.NewPage()

	url := "https://epauth.tepco.co.jp/u/login"
	log.Printf("info: url=%v", url)
	if err := page.Navigate(url); err != nil {
		log.Fatal(err)
	}
	time.Sleep(3 * time.Second)

	username := os.Getenv("TEPCO_USERNAME")
	password := os.Getenv("TEPCO_PASSWORD")

	log.Printf("info: username=%v", username)
	if err := page.FindByID("username").Fill(username); err != nil {
		log.Fatal(err)
	}
	time.Sleep(1 * time.Second)
	log.Printf("info: password=%v", strings.Repeat("*", len(password)))
	if err := page.FindByID("password").Fill(password); err != nil {
		log.Fatal(err)
	}
	time.Sleep(5 * time.Second)

	log.Printf("info: click login button")
	if err := page.FindByButton("ログイン").Click(); err != nil {
		log.Fatal(err)
	}
	defer logout(page)
	time.Sleep(3 * time.Second)

	log.Printf("info: close popup")
	for closeAttempts := 0; ; closeAttempts++ {
		closeClasses := []string{"close_icon", "close_about", "close_button"}
		ok := true
		for _, closeClass := range closeClasses {
			if err := page.FindByClass(closeClass).Click(); err != nil {
				if !strings.Contains(err.Error(), "element not found") {
					log.Printf("error: %v", err)
					return
				}
			} else {
				log.Printf("info: click %v", closeClass)
				ok = false
			}
			time.Sleep(1 * time.Second)
		}
		if ok {
			break
		}
		if closeAttempts > 20 {
			page.Screenshot("screenshot.png")
			log.Print("error: too many close popup")
		}
	}

	msg := ""

	log.Printf("info: monthly usage")
	nowPrice, err := page.FindByClass("price").Text()
	if err != nil {
		log.Printf("error: %v", err)
	}
	log.Printf("info: nowPrice=%v", nowPrice)
	msg += "nowPrice=" + nowPrice + "\n"

	forecastPrice, err := page.FindByClass("price_forecast").FindByClass("txt_red").Text()
	if err != nil {
		log.Printf("error: %v", err)
	}
	log.Printf("info: forecastPrice=%v", forecastPrice)
	msg += "forecastPrice=" + forecastPrice + "\n"

	log.Printf("info: daily usage")
	cnt, err := page.AllByClass("gaclick").Count()
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < cnt; i++ {
		ele := page.AllByClass("gaclick").At(i)
		if txt, err := ele.Text(); err != nil || txt != "日" {
			continue
		}
		if err := ele.Click(); err != nil {
			log.Printf("error: %v", err)
		}
		break
	}
	time.Sleep(3 * time.Second)

	yesterdayPrice, err := page.FindByClass("price").Text()
	if err != nil {
		log.Printf("error: %v", err)
	}
	log.Printf("info: yesterdayPrice=%s", yesterdayPrice)
	msg += "yesterdayPrice=" + yesterdayPrice + "\n"

	yesterdayUsage, err := page.FindByClass("kwh").Text()
	if err != nil {
		log.Printf("error: %v", err)
	}
	log.Printf("info: yesterdayUsage=%s", yesterdayUsage)
	msg += "yesterdayUsage=" + yesterdayUsage + "\n"
	time.Sleep(3 * time.Second)

	if _, _, err = api.PostMessage(slackChannelID, slack.MsgOptionText(msg, true)); err != nil {
		log.Printf("error: %v", err)
	}
}
