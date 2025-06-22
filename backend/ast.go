package main

import (
	"fmt"
	"strings"
)

// Node represents a node in the AST
type Node interface {
	String() string
}

// Document represents the root of the HTML document
type Document struct {
	Children []Node
}

func (d *Document) String() string {
	var result strings.Builder
	for _, child := range d.Children {
		result.WriteString(child.String())
	}
	return result.String()
}

// Element represents an HTML element
type Element struct {
	TagName    string
	Attributes map[string]string
	Children   []Node
	SelfClosing bool
}

func (e *Element) String() string {
	var result strings.Builder
	
	// Opening tag
	result.WriteString("<")
	result.WriteString(e.TagName)
	
	// Attributes
	for name, value := range e.Attributes {
		result.WriteString(" ")
		result.WriteString(name)
		if value != "" {
			result.WriteString("=\"")
			result.WriteString(value)
			result.WriteString("\"")
		}
	}
	
	if e.SelfClosing {
		result.WriteString(" />")
		return result.String()
	}
	
	result.WriteString(">")
	
	// Children
	for _, child := range e.Children {
		result.WriteString(child.String())
	}
	
	// Closing tag
	result.WriteString("</")
	result.WriteString(e.TagName)
	result.WriteString(">")
	
	return result.String()
}

// TextNode represents a text node
type TextNode struct {
	Content string
}

func (t *TextNode) String() string {
	return t.Content
}

// CommentNode represents an HTML comment
type CommentNode struct {
	Content string
}

func (c *CommentNode) String() string {
	return fmt.Sprintf("<!-- %s -->", c.Content)
}

// ImageNode represents an image element
type ImageNode struct {
	Src string
	Alt string
}

func (i *ImageNode) String() string {
	if i.Alt != "" {
		return fmt.Sprintf(`<img src="%s" alt="%s" />`, i.Src, i.Alt)
	}
	return fmt.Sprintf(`<img src="%s" />`, i.Src)
}
