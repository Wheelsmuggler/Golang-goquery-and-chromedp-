package main

import (
	"./mongo"
	"context"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"log"
	"strconv"
	"strings"
	"time"
)


func getHostComment(url string,ch chan bool) {
	options := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless",true),
		chromedp.Flag("blink-settings","imageEnable=false"),
		chromedp.Flag("no-sandbox",true),
		chromedp.UserAgent(`Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko)`),
	}

	c,_ := chromedp.NewExecAllocator(context.Background(),options...)

	chromeCtx, cancel := chromedp.NewContext(c,chromedp.WithLogf(log.Printf))
	_ = chromedp.Run(chromeCtx, make([]chromedp.Action, 0, 1)...)

	timeOutCtx, cancel := context.WithTimeout(chromeCtx,80*time.Second)
	defer cancel()

	log.Printf("chrome visit page%s\n",url)
	var htmlContent string
	err := chromedp.Run(timeOutCtx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(`div[class="List-item"]`),
		chromedp.Sleep(2*time.Second),
		chromedp.MultiClick(`button[class="Button ContentItem-action Button--plain Button--withIcon Button--withLabel"]`),
		chromedp.WaitVisible(`div[class="CommentListV2"]`),
		chromedp.OuterHTML(`document.querySelector("body")`,&htmlContent,chromedp.ByJSPath),

	)
	if err!=nil {
		log.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		log.Fatal(err)
	}
	doc.Find("div[class=\"List-item\"]").
		Each(func(i int, selection *goquery.Selection) {

		card, exist := selection.Find("div[class=\"ContentItem PinItem\"]").Attr("data-za-extra-module")
		if !exist {
			log.Fatal("card not exist")
		}
		cardMap := make(map[string]map[string]map[string]string)
		err:= json.Unmarshal([]byte(card),&cardMap)
		if err != nil {
			log.Fatal(err)
		}
		id := cardMap["card"]["content"]["id"]
		//fmt.Println(id)
		txt := selection.Find("div[class=\"RichContent-inner\"]").Find("span").Text()
		//log.Println("txt:"+txt)
		keyWords := ExtractKeyWords(txt)
		var comments [4]string
		_ = selection.Find("div[class=\"Comments-container\"]").
			Find("div[class=\"CommentsV2 CommentsV2--withEditor CommentsV2-withPagination\"]").Eq(0).
			Find("div[class=\"CommentListV2\"]").
			Find("ul[class=\"NestComment\"]").Each(func(i int, selection *goquery.Selection) {
			if i < 4 {
				hotComment := selection.Find("div[class=\"RichText ztext\"]").Text()
				comments[i] = hotComment
				//log.Println("comments"+hotComment)
			}
		})
		record := &Ping{
			Id:		  id,
			OnePing:  txt,
			Comments: comments,
			Tags:     keyWords,
		}
		err =mongo.Insert(mongo.Collection,record)
		if err != nil {
			log.Fatal(err)
		}
	})
	ch <- true
}


func main() {
	defer mongo.CloseConn()
	start := time.Now()
	queue := make(chan bool)
	for j := 1; j <= 9; j++ {
		for i := 1; i <= 5; i++ {
			url := UrlZhihu + strconv.Itoa(j*5+i)
			go getHostComment(url,queue)
		}
		for i := 0; i < 5; i++ {
			<-queue
		}
	}

	elapsed := time.Since(start)
	fmt.Printf("It takes%s\n",elapsed)
}