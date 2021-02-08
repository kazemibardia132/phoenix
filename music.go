package main

import (
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
)

func (app *App) HandleMusicUpload() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		file, handler, err := r.FormFile("file")
		if err != nil {
			panic(err)
		}
		defer file.Close()
		rn := rand.New(rand.NewSource(99))
		code := rn.Uint64()
		stringified_code := strconv.FormatUint(code, 10)
		address := "/home/vader/phoenix-music/" + stringified_code + ".mp3"
		log.Println(address)
		f, err := os.Create(address)
		if err != nil {
			log.Println(err)
			res, _ := json.Marshal(map[string]interface{}{
				"error": "cannot upload the file.",
			})
			w.Write(res)
			return
		}
		defer f.Close()
		_, _ = io.Copy(f, file)
		app.db.Create(&Music{
			Name:     r.PostForm.Get("name"),
			code:     code,
			FileName: handler.Filename,
		})
		res, _ := json.Marshal(map[string]interface{}{
			"status": "ok",
			"code":   code,
		})
		w.Write(res)
	}
}