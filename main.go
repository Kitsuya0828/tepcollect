package main

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/sclevine/agouti"
)

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
	time.Sleep(3 * time.Second)

	log.Printf("info: close popup")
	for closeAttempts := 0; ; closeAttempts++ {
		closeClasses := []string{"close_icon", "close_about", "close_button"}
		ok := true
		for _, closeClass := range closeClasses {
			if err := page.FindByClass(closeClass).Click(); err != nil {
				if !strings.Contains(err.Error(), "element not found") {
					log.Fatal(err)
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
			log.Fatal("error: too many close popup")
		}
	}

	log.Printf("info: monthly usage")
	nowPrice, err := page.FindByClass("price").Text()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("info: nowPrice=%v", nowPrice)

	forecastPrice, err := page.FindByClass("price_forecast").FindByClass("txt_red").Text()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("info: forecastPrice=%v", forecastPrice)

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
			log.Fatal(err)
		}
		break
	}
	time.Sleep(3 * time.Second)

	yesterdayPrice, err := page.FindByClass("price").Text()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("info: yesterdayPrice=%v", yesterdayPrice)

	yesterdayUsage, err := page.FindByClass("kwh").Text()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("info: yesterdayUsage=%v", yesterdayUsage)
	time.Sleep(3 * time.Second)
}