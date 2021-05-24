package brief

import (
	"io"
	"text/scanner"
)

const (
	// Whitespace no newlines
	Whitespace = 1<<' ' | 1<<'\t'
	// Newline no space
	Newline = 1<<'\n' | 1<<'\r'
	// TabCount default
	TabCount = 4
)

// Scanner for brief language
type Scanner struct {
	scanner.Scanner
	LineStart bool
	TabCount  int
	Indent    int
	atStart   bool
}

// Init contents of the brief token scanner
func (s *Scanner) Init(src io.Reader, tabsize int) *Scanner {
	s.Scanner.Init(src)
	s.Whitespace = Whitespace
	s.atStart = true
	s.TabCount = tabsize
	s.Mode = scanner.ScanIdents | scanner.ScanStrings | scanner.ScanRawStrings | scanner.ScanComments | scanner.ScanFloats
	return s
}

// Scan next token from input
func (s *Scanner) Scan() rune {
	// reads whitespace and if newline count indent
	s.countIndent()
	return s.Scanner.Scan()
}

func (s *Scanner) countIndent() {
	ch := s.Peek()
	if s.atStart {
		s.LineStart = true
		s.atStart = false
		// indent on first line
		if s.Whitespace&(1<<uint(ch)) != 0 {
			s.readIndent(ch)
		}
		return
	}
	// skip whitespace just as the text scanner does
	for s.Whitespace&(1<<uint(ch)) != 0 {
		s.Next()
		ch = s.Peek()
	}
	switch ch {
	case '\r', '\n':
		s.LineStart = true
		s.readIndent(ch)
	default:
		s.LineStart = false
	}
}

func (s *Scanner) readIndent(ch rune) {
	for {
		s.Indent = 0

		// trim all newlines
		for Newline&(1<<uint(ch)) != 0 {
			s.Next()
			ch = s.Peek()
		}
		// count whitespace
		for Whitespace&(1<<uint(ch)) != 0 {
			switch ch {
			case ' ':
				s.Indent++
			case '\t':
				s.Indent += s.TabCount
			}
			s.Next()
			ch = s.Peek()
		}

		// skip blank lines
		if Newline&(1<<uint(ch)) == 0 {
			break
		}
	}
}
