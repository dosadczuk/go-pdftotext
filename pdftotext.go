// Package pdftotext is a wrapper around Xpdf command line tool `pdftotext`.
//
// What is `pdftotext`?
//
//	Pdftotext converts Portable Document Format (PDF) file to plain text.
//
// Reference: https://www.xpdfreader.com/pdftotext-man.html
package pdftotext

import (
	"bytes"
	"io"
	"os/exec"
	"strconv"
)

// ----------------------------------------------------------------------------
// -- `pdftotext`
// ----------------------------------------------------------------------------

type command struct {
	path string
	args []string
}

// NewCommand creates new `pdftotext` command.
func NewCommand(opts ...option) *command {
	cmd := &command{path: "/usr/bin/pdftotext"}
	for _, opt := range opts {
		opt(cmd)
	}

	return cmd
}

// Run executes prepared `pdftotext` command.
func (c *command) Run(inpath string) (io.Reader, error) {
	cmd := exec.Command(c.path, append(c.args, inpath, "-")...)

	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(out), nil
}

// String returns a human-readable description of the command.
func (c *command) String() string {
	return exec.Command(c.path, append(c.args, "<inpath>")...).String()
}

// ----------------------------------------------------------------------------
// -- `pdftotext` options
// ----------------------------------------------------------------------------

type option func(*command)

// Set custom location for `pdftotext` executable.
func WithCustomPath(path string) option {
	return func(c *command) {
		c.path = path
	}
}

// Read config-file in place of ~/.xpdfrc or the system-wide config file.
func WithCustomConfig(path string) option {
	return func(c *command) {
		c.args = append(c.args, "-cfg", path)
	}
}

// Specifies the first page to convert.
func WithPageFrom(page uint64) option {
	return func(c *command) {
		c.args = append(c.args, "-f", strconv.FormatUint(page, 10))
	}
}

// Specifies the last page to convert.
func WithPageTo(page uint64) option {
	return func(c *command) {
		c.args = append(c.args, "-l", strconv.FormatUint(page, 10))
	}
}

// Specifies the range of pages to convert.
func WithPageRange(from, to uint64) option {
	return func(c *command) {
		WithPageFrom(from)
		WithPageTo(to)
	}
}

// Maintain (as best as possible) the original physical layout of the text.
func WithModeLayout() option {
	return func(c *command) {
		c.args = append(c.args, "-layout")
	}
}

// Similar to `WithModeLayout`, but optimized for simple one-column pages.
//
// This mode will do a better job of maintaining horizontal spacing, but it
// will only work properly with a single column of text.
func WithModeSimple() option {
	return func(c *command) {
		c.args = append(c.args, "-simple")
	}
}

// Similar to `WithModeSimple` but handles slightly rotated text better.
//
// Only works for pages with a single column of text.
func WithModeSimple2() option {
	return func(c *command) {
		c.args = append(c.args, "-simple2")
	}
}

// Table mode is similar to physical layout mode, but optimized for tabular
// data, with the goal of keeping rows and columns aligned (at the expense of
// inserting extra whitespace).
//
// If the `WithCharFixedWidth` option is given, character spacing within each
// line will be determined by the specified character pitch.
func WithModeTable() option {
	return func(c *command) {
		c.args = append(c.args, "-table")
	}
}

// Line printer mode uses a strict fixed-character-pitch and -height layout.
// The page is broken into a grid, and characters are placed into that grid.
//
// If the grid spacing is too small for the actual characters, the result is
// extra whitespace. If the grid spacing is too large, the result is missing
// whitespace.
//
// Use `WithCharFixedWidth` and `WithLineFixedSpacing` to specify grid spacing.
// If one or both are not given on the command line, it will attempt to compute
// appropriate value(s).
func WithModeLinePrinter() option {
	return func(c *command) {
		c.args = append(c.args, "-lineprinter")
	}
}

// Keep the text in content stream order.
//
// Depending on how the PDF file was generated, this may or may not be useful.
func WithModeRaw() option {
	return func(c *command) {
		c.args = append(c.args, "-raw")
	}
}

// Specify the character pitch (width), in points.
//
// Works only with `WithModeLayout`, `WithModeTable` and `WithModeLinePrinter`.
func WithCharFixedWidth(width uint64) option {
	return func(c *command) {
		c.args = append(c.args, "-fixed", strconv.FormatUint(width, 10))
	}
}

// Specify the line spacing, in points.
//
// Works only with `WithModeLinePrinter`.
func WithLineFixedSpacing(spacing uint64) option {
	return func(c *command) {
		c.args = append(c.args, "-linespacing", strconv.FormatUint(spacing, 10))
	}
}

// Text which is hidden because of clipping is removed before doing layout,
// and then added back in.
//
// This can be helpful for tables where clipped (invisible) text would overlap
// the next column.
func WithTextClipping() option {
	return func(c *command) {
		c.args = append(c.args, "-clip")
	}
}

// Diagonal text, i.e., text that is not close to one of the 0, 90, 180, or 270
// degree axes, is discarded.
//
// This is useful to skip watermarks drawn on top of body text, etc.
func WithNoTextDiagonal() option {
	return func(c *command) {
		c.args = append(c.args, "-nodiag")
	}
}

// Sets the encoding to use for text output.
//
// The name must be defined with the unicodeMap command (see xpdfrc(5)).
// The encoding name is case-sensitive. This defaults to "Latin1".
//
// Available options: `pdftotext -listencodings`.
func WithEncoding(name string) option {
	return func(c *command) {
		c.args = append(c.args, "-enc", name)
	}
}

// Sets the end-of-line convention to use for text output.
//
// Available options: "unix", "dos", "mac".
func WithEndOfLine(kind string) option {
	return func(c *command) {
		c.args = append(c.args, "-eol", kind)
	}
}

// Don’t insert a page breaks (form feed character) at the end of each page.
func WithNoPageBreak() option {
	return func(c *command) {
		c.args = append(c.args, "-nopgbrk")
	}
}

// Insert a Unicode byte order marker (BOM) at the start of the text output.
func WithByteOrderMarker() option {
	return func(c *command) {
		c.args = append(c.args, "-bom")
	}
}

// Specifies the left margin, in points.
//
// Text in the left margin (i.e., within that many points of the left edge
// of the page) is discarded.
func WithMarginLeft(margin uint64) option {
	return func(c *command) {
		c.args = append(c.args, "-marginl", strconv.FormatUint(margin, 10))
	}
}

// Specifies the right margin, in points.
//
// Text in the right margin (i.e., within that many points of the right edge
// of the page) is discarded.
func WithMarginRight(margin uint64) option {
	return func(c *command) {
		c.args = append(c.args, "-marginr", strconv.FormatUint(margin, 10))
	}
}

// Specifies the top margin, in points.
//
// Text in the top margin (i.e., within that many points of the top edge
// of the page) is discarded.
func WithMarginTop(margin uint64) option {
	return func(c *command) {
		c.args = append(c.args, "-margint", strconv.FormatUint(margin, 10))
	}
}

// Specifies the bottom margin, in points.
//
// Text in the bottom margin (i.e., within that many points of the bottom edge
// of the page) is discarded.
func WithMarginBottom(margin uint64) option {
	return func(c *command) {
		c.args = append(c.args, "-marginb", strconv.FormatUint(margin, 10))
	}
}

// Specifies the margins, in points.
func WithMargin(t, r, b, l uint64) option {
	return func(c *command) {
		WithMarginTop(t)
		WithMarginRight(r)
		WithMarginBottom(b)
		WithMarginLeft(l)
	}
}

// Specify the owner password for the PDF file.
//
// Providing this will bypass all security restrictions.
func WithOwnerPassword(password string) option {
	return func(c *command) {
		c.args = append(c.args, "-opw", password)
	}
}

// Specify the user password for the PDF file.
func WithUserPassword(password string) option {
	return func(c *command) {
		c.args = append(c.args, "-upw", password)
	}
}
