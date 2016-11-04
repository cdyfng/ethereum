package main

import (
	"fmt"
	"log"

	"flag"
	"github.com/PuerkitoBio/goquery"
	"github.com/robfig/config"
	"net/smtp"
	"strconv"
	"time"
)

func getEtherPrice() float32 {
	doc, err := goquery.NewDocument("http://coinmarketcap.com/currencies/ethereum/")
	if err != nil {
		log.Fatal(err)
	}

	var global_price float32

	// Find the review items
	//doc.Find(".sidebar-reviews article .content-block").Each(func(i int, s *goquery.Selection) {
	doc.Find("span#quote_price.text-large").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		ether_price, _ := s.Html()
		price, error := strconv.ParseFloat(ether_price[1:], 32)
		if error != nil {
			fmt.Println("字符串转换成整数失败")
		}
		//fmt.Printf("Review %d: %.3f \n", i, price)
		global_price = float32(price)

	})
	return global_price

}

func sendEmail() error {
	return nil
}

func loadConfig() *Email_setting {
	var emailSetting = make(map[string]string)
	var configFile = flag.String("configfile", "/Users/acer/config.ini", "General configuration file")
	flag.Parse()
	cfg, err := config.ReadDefault(*configFile)
	if err != nil {
		log.Fatalf("Fail to find", *configFile, err)
	}
	//set config file std End

	//Initialized topic from the configuration
	if cfg.HasSection("emailSetting") {
		section, err := cfg.SectionOptions("emailSetting")
		if err == nil {
			for _, v := range section {
				options, err := cfg.String("emailSetting", v)
				if err == nil {
					emailSetting[v] = options
				}
			}
		}
	}
	//Initialized topic from the configuration END

	fmt.Println(emailSetting)
	fmt.Println(emailSetting["username"])

	var email_setting *Email_setting = &Email_setting{}
	email_setting.identity = ""
	email_setting.username = emailSetting["username"]
	email_setting.password = emailSetting["password"]
	email_setting.host = emailSetting["host"]
	email_setting.addr = emailSetting["addr"]
	email_setting.from = emailSetting["from"]
	email_setting.to = []string{emailSetting["to"]}
	//email_setting.msg_str = emailSetting["msg_str"]
	return email_setting
}

func mail(emailSettingPtr *Email_setting) {
	// Set up authentication information.
	auth := smtp.PlainAuth(
		emailSettingPtr.identity,
		emailSettingPtr.username,
		emailSettingPtr.password,
		emailSettingPtr.host,
	)
	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	err := smtp.SendMail(
		emailSettingPtr.addr,
		auth,
		emailSettingPtr.from,
		emailSettingPtr.to,
		emailSettingPtr.msg,
	)
	if err != nil {
		log.Fatal(err)
	}
}

type Price_Setting struct {
	lastSavePrice float32
	notifyPercent float32
	nowPrice      float32
}

type Email_setting struct {
	identity, username, password, host string
	addr, from                         string
	to                                 []string
	msg                                []byte
}

func price_compare(ps *Price_Setting, ps_out chan *Price_Setting) {
	for {
		fmt.Print(".")
		ps.nowPrice = getEtherPrice()
		if ps.nowPrice > ps.lastSavePrice*(1+0.05) {
			//if ps.nowPrice != ps.lastSavePrice {
			fmt.Printf("diff , %f > %f \n", ps.nowPrice, ps.lastSavePrice)
			ps_out <- ps
			time.Sleep(time.Millisecond * 5)
			ps.lastSavePrice = ps.nowPrice
		} else if ps.nowPrice < ps.lastSavePrice*(1-0.05) {
			fmt.Printf("warning low, %f < %f \n", ps.nowPrice, ps.lastSavePrice)
			ps_out <- ps
			time.Sleep(time.Millisecond * 5)
			ps.lastSavePrice = ps.nowPrice

		}
		//fmt.Println(ps)
		time.Sleep(time.Second * 60)
	}
}

func main() {
	//init price, return last_save_price, notify_percent, email_setting
	// for{}
	//get new price
	//check if exceed the bounds, if so, send notify-email & save new last_save_price
	//delay 1 miniter
	//p := getEtherPrice()
	//fmt.Println(p)
	//mail()
	es := loadConfig()
	ps := &Price_Setting{10.0, 0.05, 0}
	ps_out := make(chan *Price_Setting)
	go price_compare(ps, ps_out)

	for {
		select {
		case out := <-ps_out:
			fmt.Println("diff in main: ", out)
			msg := fmt.Sprintf("To: "+es.to[0]+"\r\nFrom: f<f>\r\nSubject: %f,%f\r\nContent-Type: text/plain\r\n\r\nno content", out.nowPrice, out.lastSavePrice)
			fmt.Println(msg)
			es.msg = []byte(msg)
			mail(es)

		}
	}

}
