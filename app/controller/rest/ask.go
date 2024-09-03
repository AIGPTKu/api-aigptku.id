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
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
)

func (r *restHandler) ask(c *fiber.Ctx) (err error) {

	var (
		req = domainCt.RequestAsk{}
		res = domainCt.ResponseAsk{}
		content = make(chan string)
		finish = make(chan bool)
	)

	c.BodyParser(&req)
	fmt.Println(req)

	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")

	c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		keepAliveTickler := time.NewTicker(15 * time.Second)

		for loop := true; loop; {
			select {
			case stop := <-finish:
				if stop {
					keepAliveTickler.Stop()
					loop = false
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
		go r.uc.gpt.AskGPT(c.UserContext(), content, finish, req.Contents)
	} else if engine == "gemini" {
		// go r.uc.gemini.AskGPT(c.UserContext)	
	}

	return nil
}