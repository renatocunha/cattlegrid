package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type Forms struct {
	FileName  string
	Label     string
	Published bool
}

func upLoader(w http.ResponseWriter, req *http.Request) {

	if req.Method == http.MethodPost {

		fmt.Println("Upload!")

		f, h, err := req.FormFile("uploadfile")

		req.ParseForm()
		fmt.Println(req.Form)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			fmt.Println("error Read")
			return
		}

		form := &Forms{}

		if len(h.Filename) > 0 {
			form.FileName = h.Filename
		} else {
			fmt.Println("invalid file")
			return
		}

		if label, ok := req.Form["label"]; ok {
			if len(label[0]) > 0 {
				form.Label = label[0]
			}
		}

		if pub, ok := req.Form["published"]; ok {
			if len(pub) > 0 {
				form.Published = true
			}

		} else {
			form.Published = false
		}

		defer f.Close()

		//print values
		bs, err := ioutil.ReadAll(f)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			fmt.Println("error ReadAll")
			return
		}

		// store on server
		dst, err := os.Create(filepath.Join("./storage/", h.Filename))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			fmt.Println("error Create")
			return
		}
		defer dst.Close()

		_, err = dst.Write(bs)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			fmt.Println("error Write")
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}
}

func main() {
	http.HandleFunc("/upload", upLoader)
	http.Handle("/", http.FileServer(http.Dir("./web")))
	log.Fatal(http.ListenAndServe(":8090", nil))
}
