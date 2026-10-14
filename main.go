// Command site builds spinney.dev from the markdown in content/, set in
// templates/layout.html: content/home.md becomes public/index.html, and
// every other page, such as content/projects.md, public/projects/index.html,
// reached as projects/. Run it from the site folder:
//
//	go run .            # regenerate the pages in public/
//	go run . -serve     # the same, then serve public/ on 127.0.0.1:8000,
//	                    # rebuilding on every page load
//	go run . -serve -addr 127.0.0.1:9000
//
// A page is a markdown file with its title between two --- lines on top:
//
//	---
//	title: projects
//	---
//
//	- [canopy](https://canopy.spinney.dev/): a small web viewer for git
package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

// page is what the layout gets.
type page struct {
	Name     string // the file's name without .md; "home" is the home page
	Title    string
	Phonetic string // how the title is said, if given, as /ˈspɪni/
	Body     template.HTML
}

var md = goldmark.New(
	goldmark.WithExtensions(extension.GFM),
	goldmark.WithRendererOptions(html.WithUnsafe()),
)

func main() {
	out := flag.String("out", "public", "the site folder")
	serve := flag.Bool("serve", false, "serve the site after building it, and rebuild on every page load")
	addr := flag.String("addr", "127.0.0.1:8000", "address to serve on, with -serve")
	flag.Parse()
	log.SetFlags(0)

	if err := build(*out); err != nil {
		log.Fatal(err)
	}
	if !*serve {
		return
	}
	files := http.FileServer(http.Dir(*out))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if p := r.URL.Path; strings.HasSuffix(p, "/") || strings.HasSuffix(p, ".html") {
			if err := build(*out); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		files.ServeHTTP(w, r)
	})
	log.Printf("serving %s on http://%s/", *out, *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

// build writes a page into out for every markdown file in content/.
func build(out string) error {
	t, err := template.ParseFiles("templates/layout.html")
	if err != nil {
		return err
	}
	srcs, err := filepath.Glob("content/*.md")
	if err != nil {
		return err
	}
	for _, src := range srcs {
		p, err := parse(src)
		if err != nil {
			return err
		}
		dir := p.Name
		if dir == "home" {
			dir = ""
		}
		var b bytes.Buffer
		if err := t.ExecuteTemplate(&b, "layout", p); err != nil {
			return err
		}
		path := filepath.Join(out, dir, "index.html")
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(path, b.Bytes(), 0o644); err != nil {
			return err
		}
	}
	return nil
}

// parse reads one markdown file: its front matter, then its body.
func parse(src string) (page, error) {
	raw, err := os.ReadFile(src)
	if err != nil {
		return page{}, err
	}
	p := page{Name: strings.TrimSuffix(filepath.Base(src), ".md")}
	parts := strings.SplitN(strings.ReplaceAll(string(raw), "\r\n", "\n"), "---\n", 3)
	if len(parts) != 3 || parts[0] != "" {
		return page{}, fmt.Errorf("%s: no front matter: start with ---, title: ..., ---", src)
	}
	for _, line := range strings.Split(parts[1], "\n") {
		k, v, _ := strings.Cut(line, ":")
		switch strings.TrimSpace(k) {
		case "title":
			p.Title = strings.TrimSpace(v)
		case "phonetic":
			p.Phonetic = strings.TrimSpace(v)
		}
	}
	if p.Title == "" {
		return page{}, fmt.Errorf("%s: no title in the front matter", src)
	}
	var b bytes.Buffer
	if err := md.Convert([]byte(parts[2]), &b); err != nil {
		return page{}, fmt.Errorf("%s: %w", src, err)
	}
	p.Body = template.HTML(b.String())
	return p, nil
}
