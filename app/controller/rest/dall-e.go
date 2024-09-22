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

func (r *restHandler) generateImage(c *fiber.Ctx) (err error) {

	var (
		req = domainCt.RequestImageGenerate{}
		res = domainCt.ResponseImageGenerate{}
		content = make(chan string)
		image = make(chan string)
		finish = make(chan bool)
	)

	c.BodyParser(&req)

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
			case img := <-image:
				var buf bytes.Buffer
                enc := json.NewEncoder(&buf)

				res.ImageURL = img
				res.Content = ""

				err := enc.Encode(res)
				if err != nil {
					return
				}

				sb := strings.Builder{}
				sb.WriteString(fmt.Sprintf("data: %v\n", buf.String()))
	
				fmt.Println(img)

				_, err = fmt.Fprint(w, sb.String())
				if err != nil {
					log.Println(err)
				}
				err = w.Flush()
				if err != nil {
					log.Println(err)
				}
			case ev := <-content:
				var buf bytes.Buffer
				enc := json.NewEncoder(&buf)
			
				res.ImageURL = ""
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

	disableDallE := viper.GetBool("DALL_E_DISABLE")

	if engine == "gpt" {
		if disableDallE {
			go func ()  {
				time.Sleep(1000 * time.Millisecond)
				text_split := strings.Split("Mohon maaf, untuk sementara kami sedang membatasi untuk generate gambar, kamu bisa mencoba lagi nanti.", " ")
				for _, v := range text_split {
					content <- v + " "
					time.Sleep(25 * time.Millisecond)
				}
				finish <- true
			}()
		} else {
			go r.uc.gpt.GenerateImage(c.UserContext(), content, image, finish, req.Prompt)
		}
	} else if engine == "gemini" {
		// go r.uc.gemini.AskGPT(c.UserContext)	
	}

	return nil
}