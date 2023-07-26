package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/russross/blackfriday/v2"
)

const (
	ansiBold       = "\033[1m"
	ansiUnderlined = "\033[4m"
	ansiItalic     = "\033[3m"
	ansiReset      = "\033[0m"
)

type ANSIFormatRenderer struct {
	links map[int]string
}

func main() {
	if _, err := os.Stat("offline"); os.IsNotExist(err) {
		err := os.Mkdir("offline", 0755)
		if err != nil {
			fmt.Println("Error creating 'offline' directory:", err)
			os.Exit(1)
		}
	}

	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		fmt.Println("Please provide a URL")
		os.Exit(1)
	}

	for {
		renderer := getAndRenderMarkdown(args[0])

		fmt.Println("\nEnter a link number to open it, or 'q' to quit:")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "q" {
			os.Exit(0)
		}

		linkNumber, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Invalid input, please enter a number or 'q'")
			continue
		}

		link, ok := renderer.links[linkNumber]
		if !ok {
			fmt.Println("No such link number, try again")
			continue
		}

		args[0] = link
	}
}

func (r *ANSIFormatRenderer) RenderNode(w io.Writer, node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
	switch node.Type {
	case blackfriday.Text:
		switch node.Parent.Type {
		case blackfriday.Strong:
			if entering {
				io.WriteString(w, ansiBold)
			}
			w.Write(node.Literal)
			if !entering {
				io.WriteString(w, ansiReset)
			}
		case blackfriday.Emph:
			if entering {
				io.WriteString(w, ansiItalic)
			}
			w.Write(node.Literal)
			if !entering {
				io.WriteString(w, ansiReset)
			}
		default:
			w.Write(node.Literal)
		}
	case blackfriday.Link:
		if entering {
			linkNumber := len(r.links) + 1
			r.links[linkNumber] = string(node.LinkData.Destination)
			io.WriteString(w, fmt.Sprintf("[%d]", linkNumber))
		}
	case blackfriday.Heading:
		if node.HeadingData.Level == 1 {
			if entering {
				io.WriteString(w, ansiBold+ansiUnderlined)
			} else {
				io.WriteString(w, ansiReset+"\n\n")
			}
		} else if node.HeadingData.Level == 2 {
			if entering {
				io.WriteString(w, ansiUnderlined)
			} else {
				io.WriteString(w, ansiReset+"\n\n")
			}
		}
	case blackfriday.Item:
		if entering {
			io.WriteString(w, "- ")
		}
	case blackfriday.Paragraph:
		if !entering {
			io.WriteString(w, "\n")
		}
	}
	return blackfriday.GoToNext
}

func (r *ANSIFormatRenderer) RenderHeader(w io.Writer, ast *blackfriday.Node) {}
func (r *ANSIFormatRenderer) RenderFooter(w io.Writer, ast *blackfriday.Node) {}
func getAndRenderMarkdown(url string) *ANSIFormatRenderer {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching URL:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	var text strings.Builder
	for scanner.Scan() {
		text.WriteString(scanner.Text())
		text.WriteString("\n")
	}

	domainName := strings.Split(resp.Request.URL.Host, ":")[0]
	fileName := resp.Request.URL.Path
	if fileName == "/" {
		fileName = "index.md"
	} else if fileName == "/" || fileName == "" {
		fileName = "index.md"
	} else if !strings.HasSuffix(fileName, ".md") {
		fileName += ".md"
	}

	offlineDir := filepath.Join("offline", strings.TrimPrefix(domainName, "www."))
	if err := os.MkdirAll(offlineDir, 0755); err != nil {
		fmt.Println("Error creating offline directory:", err)
		os.Exit(1)
	}

	offlineFile := filepath.Join(offlineDir, fileName)
	if err := os.WriteFile(offlineFile, []byte(text.String()), 0644); err != nil {
		fmt.Println("Error writing offline file:", err)
		os.Exit(1)
	}

	renderer := &ANSIFormatRenderer{links: make(map[int]string)}
	md := blackfriday.New(blackfriday.WithRenderer(renderer))

	ast := md.Parse([]byte(text.String()))

	ast.Walk(func(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		return renderer.RenderNode(os.Stdout, node, entering)
	})

	baseURL := url[:strings.LastIndex(url, "/")+1]
	for i, link := range renderer.links {
		if !strings.Contains(link, "://") {
			newLink := baseURL + link
			renderer.links[i] = newLink
		}
	}

	fmt.Println("\n\n---\nLinks:")
	for i, link := range renderer.links {
		fmt.Printf("[%d] -> %s\n", i, link)
	}
	fmt.Println("---")

	return renderer
}
