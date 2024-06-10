package main

import (
	"fmt"
	"io"
	"log"

	"github.com/dosadczuk/pdftotext"
)

func main() {
	cmd := pdftotext.NewCommand(
		pdftotext.WithEncoding("UTF-8"),
		pdftotext.WithModeLayout(),
		pdftotext.WithMargin(20, 20, 20, 20),
		pdftotext.WithNoPageBreak(),
	)

	out, err := cmd.Run("./example.pdf")
	if err != nil {
		log.Panic(err)
	}

	txt, err := io.ReadAll(out)
	if err != nil {
		log.Panic(err)
	}

	fmt.Print(string(txt))
}
