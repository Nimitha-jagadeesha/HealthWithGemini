package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/googleai"
	"github.com/tmc/langchaingo/schema"
)

var apiKey = "AIzaSyCjolnfa-0MENwsHohybMgE5Hr5X0wGQo0"

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), "usage: web [options]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

var (
	addr = flag.String("addr", "localhost:8080", "address to serve")
)

func generateHandler(w http.ResponseWriter, r *http.Request, llm *googleai.GoogleAI) {

	prompt := r.FormValue("prompt")
	file, header, err := r.FormFile("chosen-image")
	if err != nil {
		log.Printf("Unable to read image %s %v\n", header.Filename, err)
		http.Error(w, "Error: unable to generate content", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	imgData := bytes.NewBuffer(nil)
	if _, err := io.Copy(imgData, file); err != nil {
		log.Printf("Unable to read image %v\n", err)
		http.Error(w, "Error: unable to generate content", http.StatusInternalServerError)
		return
	}

	// imgData, err := os.ReadFile(filepath.Join("static", "images", "ingridients_2.jpeg"))
	if err != nil {
		log.Printf("Unable to read image %v\n", err)
		http.Error(w, "Error: unable to generate content", http.StatusInternalServerError)
		return
	}

	content := []llms.MessageContent{
		{
			Role: schema.ChatMessageTypeHuman,
			Parts: []llms.ContentPart{
				llms.BinaryPart("image/jpeg", imgData.Bytes()),
				llms.TextPart(prompt),
			},
		},
	}

	_, err = llm.GenerateContent(r.Context(), content,
		llms.WithModel("gemini-pro-vision"),
		llms.WithMaxTokens(500), // default of 256 is not enough.
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			fmt.Fprint(w, string(chunk))
			return nil
		}),
	)
	if err != nil {
		log.Printf("Error generating content: %v\n", err)
		http.Error(w, "Error: unable to generate content", http.StatusInternalServerError)
		return
	}
}

type Page struct {
	Images []string
}

var tmpl = template.Must(template.ParseFiles("static/index.html"))

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// Load all baked goods images from the static/images directory.
	matches, err := filepath.Glob(filepath.Join("static", "images"))
	if err != nil {
		log.Printf("Error loading baked goods images: %v", err)
	}
	var page = &Page{Images: make([]string, len(matches))}
	for i, match := range matches {
		page.Images[i] = filepath.Base(match)
	}
	switch r.URL.Path {
	case "/":
		err = tmpl.Execute(w, page)
		if err != nil {
			log.Printf("Template execution error: %v", err)
		}
	}
}

func main() {
	// Parse flags.
	flag.Usage = usage
	flag.Parse()

	// Parse and validate arguments (none).
	args := flag.Args()
	if len(args) != 0 {
		usage()
	}
	// Get the Gemini API key from the environment.
	if key := os.Getenv("API_KEY"); key != "" {
		apiKey = key
	}

	llm, err := googleai.New(context.Background(), googleai.WithAPIKey(apiKey))
	if err != nil {
		log.Fatal(err)
	}

	// Serve static files and handle API requests.
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/api/generate", func(w http.ResponseWriter, r *http.Request) { generateHandler(w, r, llm) })
	http.HandleFunc("/", indexHandler)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
