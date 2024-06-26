package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/dosadczuk/go-pdftotext"
)

func main() {
	cmd, err := pdftotext.NewCommand(
		pdftotext.WithEncoding("UTF-8"),
		pdftotext.WithModeLayout(),
		pdftotext.WithMargin(20, 20, 20, 20),
		pdftotext.WithNoPageBreak(),
	)
	if err != nil {
		log.Fatal(err)
	}

	out, err := cmd.Run(context.Background(), "./example.pdf")
	if err != nil {
		log.Panic(err)
	}

	txt, err := io.ReadAll(out)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println(string(txt))
}
