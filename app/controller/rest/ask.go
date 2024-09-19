package rest

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	domainCt "github.com/AIGPTku/api-aigptku.id/app/controller/domain"
	domainUc "github.com/AIGPTku/api-aigptku.id/app/usecase/domain"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
)

func (r *restHandler) ask(c *fiber.Ctx) (err error) {

	var (
		req = domainCt.RequestAsk{}
		res = domainCt.ResponseAsk{}
		funcCall = make(chan domainCt.FuncCall)
		content = make(chan string)
		finish = make(chan bool)
	)

	c.BodyParser(&req)
	fmt.Println(req)

	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")

	ctx := c.UserContext()

	c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		keepAliveTickler := time.NewTicker(15 * time.Second)

		for loop := true; loop; {
			select {
			case stop := <-finish:
				if stop {
					keepAliveTickler.Stop()
					loop = false
				}
			case fc := <- funcCall:
				var buf bytes.Buffer
				enc := json.NewEncoder(&buf)

				fmt.Println(fc)

				if fc.Name != "web_search" && fc.Name != "image_generate" {
					go r.uc.gpt.HandleFunctionText(ctx, content, finish, fc)
					continue
				}
		
				res.FuncCall = &fc
				res.Content = ""
			
				err := enc.Encode(res)
				if err != nil {
					return
				}

				sb := strings.Builder{}
				sb.WriteString(fmt.Sprintf("data: %v\n", buf.String()))

				_, err = fmt.Fprint(w, sb.String())
				if err != nil {
					log.Println(err)
				}
				time.AfterFunc(100 * time.Millisecond, func() {
					finish <- true
				})
				err = w.Flush()
				if err != nil {
					log.Println(err)
				}
			case ev := <-content:
				var buf bytes.Buffer
				enc := json.NewEncoder(&buf)
			
				res.Content = ev
			
				err := enc.Encode(res)
				if err != nil {
					return
				}

				sb := strings.Builder{}
				sb.WriteString(fmt.Sprintf("data: %v\n", buf.String()))
	
				fmt.Print(ev)

				_, err = fmt.Fprint(w, sb.String())
				if err != nil {
					log.Println(err)
				}
				err = w.Flush()
				if err != nil {
					log.Println(err)
				}
			case <- keepAliveTickler.C:
				var buf bytes.Buffer
				enc := json.NewEncoder(&buf)
			
				m := domainCt.ResponseAsk{
					Content: "",
				}
			
				err := enc.Encode(m)
				if err != nil {
					return
				}

				sb := strings.Builder{}
				sb.WriteString(fmt.Sprintf("data: %v\n", buf.String()))

				_, err = fmt.Fprint(w, sb.String())
				if err != nil {
					log.Println(err)
				}
				err = w.Flush()
				if err != nil {
					log.Printf("Error while flushing: %v.\n", err)
					keepAliveTickler.Stop()
					loop = false
				}
			}
		}
	}))
	
	engine := viper.GetString("NLP_ENGINE")
	if engine == "" {
		engine = "gpt"
	}

	if engine == "gpt" {
		go r.uc.gpt.AskGPT(ctx, domainUc.RequestAsk{
			FuncCall: funcCall,
			Result: content,
            Finish: finish,
            AskContent: req.Contents,
            UseDefaultSystem: true,
		})
	} else if engine == "gemini" {
		// go r.uc.gemini.AskGPT(c.UserContext)	
	}

	return nil
}