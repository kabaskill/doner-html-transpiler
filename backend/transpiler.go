package main

import (
	"fmt"
	"strings"
)

// Transpiler handles the conversion from German HTML to standard HTML
type Transpiler struct {
	dictionary *Dictionary
}

// NewTranspiler creates a new transpiler instance
func NewTranspiler() *Transpiler {
	return &Transpiler{
		dictionary: NewDictionary(),
	}
}

// Transpile converts German HTML to standard HTML
func (t *Transpiler) Transpile(input string) (string, error) {
	// Create lexer
	lexer := NewLexer(input)
	
	// Create parser
	parser := NewParser(lexer, t.dictionary)
	
	// Parse into AST
	document, err := parser.Parse()
	if err != nil {
		return "", fmt.Errorf("parsing error: %w", err)
	}
	
	// Convert AST back to HTML string
	result := document.String()
	
	// Pretty print the result
	return t.formatHTML(result), nil
}

// formatHTML provides basic formatting for the HTML output
func (t *Transpiler) formatHTML(html string) string {
	// First, clean up the HTML and add newlines between tags
	html = strings.ReplaceAll(html, "><", ">\n<")
	
	// Split into lines and process each line
	lines := strings.Split(html, "\n")
	var formatted strings.Builder
	indent := 0
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		// Check if this is a closing tag
		isClosingTag := strings.HasPrefix(line, "</")
		// Check if this is a self-closing tag or contains both opening and closing
		isSelfClosing := strings.HasSuffix(line, "/>") || 
			(strings.Contains(line, "</") && strings.Contains(line, ">") && !isClosingTag)
		
		// For closing tags, decrease indent before printing
		if isClosingTag {
			indent--
			if indent < 0 {
				indent = 0
			}
		}
		
		// Add indentation
		for i := 0; i < indent; i++ {
			formatted.WriteString("  ")
		}
		formatted.WriteString(line)
		formatted.WriteString("\n")
		
		// For opening tags (that are not self-closing), increase indent after printing
		if strings.HasPrefix(line, "<") && !isClosingTag && !isSelfClosing {
			indent++
		}
	}
	
	return formatted.String()
}

// GetSupportedTags returns a map of supported German tags to HTML tags
func (t *Transpiler) GetSupportedTags() map[string]string {
	return t.dictionary.tags
}

// GetSupportedAttributes returns a map of supported German attributes to HTML attributes
func (t *Transpiler) GetSupportedAttributes() map[string]string {
	return t.dictionary.attributes
}
