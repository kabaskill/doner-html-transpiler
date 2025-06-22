package main

import (
	"fmt"
	"strings"
)

// Parser parses tokens into an AST
type Parser struct {
	lexer        *Lexer
	currentToken Token
	peekToken    Token
	dictionary   *Dictionary
}

// NewParser creates a new parser instance
func NewParser(lexer *Lexer, dictionary *Dictionary) *Parser {
	p := &Parser{
		lexer:      lexer,
		dictionary: dictionary,
	}
	
	// Read two tokens, so currentToken and peekToken are both set
	p.nextToken()
	p.nextToken()
	
	return p
}

// nextToken advances to the next token
func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

// Parse parses the input and returns a Document AST
func (p *Parser) Parse() (*Document, error) {
	doc := &Document{Children: []Node{}}
	
	for p.currentToken.Type != TOKEN_EOF {
		node, err := p.parseNode()
		if err != nil {
			return nil, err
		}
		if node != nil {
			doc.Children = append(doc.Children, node)
		}
		p.nextToken()
	}
	
	return doc, nil
}

// parseNode parses a single node (element or text)
func (p *Parser) parseNode() (Node, error) {
	switch p.currentToken.Type {
	case TOKEN_TAG_OPEN:
		return p.parseElement()
	case TOKEN_TEXT:
		if strings.TrimSpace(p.currentToken.Value) == "" {
			return nil, nil // Skip empty text nodes
		}
		return &TextNode{Content: p.currentToken.Value}, nil
	default:
		return nil, fmt.Errorf("unexpected token: %s", p.currentToken)
	}
}

// parseElement parses an HTML element
func (p *Parser) parseElement() (*Element, error) {
	// Expect opening tag
	if p.currentToken.Type != TOKEN_TAG_OPEN {
		return nil, fmt.Errorf("expected '<', got %s", p.currentToken)
	}
	
	p.nextToken() // consume '<'
	
	// Get tag name
	if p.currentToken.Type != TOKEN_TAG_NAME {
		return nil, fmt.Errorf("expected tag name, got %s", p.currentToken)
	}
	
	germanTagName := p.currentToken.Value
	htmlTagName, exists := p.dictionary.TranslateTag(germanTagName)
	if !exists {
		htmlTagName = germanTagName // Keep original if no translation exists
	}
	
	element := &Element{
		TagName:    htmlTagName,
		Attributes: make(map[string]string),
		Children:   []Node{},
	}
	
	p.nextToken() // consume tag name
	
	// Parse attributes
	for p.currentToken.Type == TOKEN_ATTR_NAME {
		attr, err := p.parseAttribute()
		if err != nil {
			return nil, err
		}
		element.Attributes[attr.Name] = attr.Value
	}
	
	// Check for self-closing tag
	if p.currentToken.Type == TOKEN_TAG_CLOSE_SLASH {
		element.SelfClosing = true
		return element, nil
	}
	
	// Expect closing '>'
	if p.currentToken.Type != TOKEN_TAG_CLOSE {
		return nil, fmt.Errorf("expected '>' or '/>', got %s", p.currentToken)
	}
	
	p.nextToken() // consume '>'
	
	// Parse children until we find the closing tag
	for p.currentToken.Type != TOKEN_TAG_END && p.currentToken.Type != TOKEN_EOF {
		child, err := p.parseNode()
		if err != nil {
			return nil, err
		}
		if child != nil {
			element.Children = append(element.Children, child)
		}
		p.nextToken()
	}
	
	// Check if we hit EOF without finding closing tag
	if p.currentToken.Type == TOKEN_EOF {
		return nil, fmt.Errorf("unexpected end of input: missing closing tag for <%s>", htmlTagName)
	}
	
	// Parse closing tag
	if p.currentToken.Type == TOKEN_TAG_END {
		p.nextToken() // consume '</'
		
		if p.currentToken.Type != TOKEN_TAG_NAME {
			return nil, fmt.Errorf("expected closing tag name, got %s", p.currentToken)
		}
		
		closingTagName := p.currentToken.Value
		closingHtmlTagName, exists := p.dictionary.TranslateTag(closingTagName)
		if !exists {
			closingHtmlTagName = closingTagName
		}
		
		if closingHtmlTagName != htmlTagName {
			return nil, fmt.Errorf("mismatched closing tag: expected %s, got %s", htmlTagName, closingHtmlTagName)
		}
		
		p.nextToken() // consume closing tag name
		
		if p.currentToken.Type != TOKEN_TAG_CLOSE {
			return nil, fmt.Errorf("expected '>', got %s", p.currentToken)
		}
	}
	
	return element, nil
}

// Attribute represents an HTML attribute
type Attribute struct {
	Name  string
	Value string
}

// parseAttribute parses an HTML attribute
func (p *Parser) parseAttribute() (*Attribute, error) {
	if p.currentToken.Type != TOKEN_ATTR_NAME {
		return nil, fmt.Errorf("expected attribute name, got %s", p.currentToken)
	}
	
	germanAttrName := p.currentToken.Value
	htmlAttrName, exists := p.dictionary.TranslateAttribute(germanAttrName)
	if !exists {
		htmlAttrName = germanAttrName // Keep original if no translation exists
	}
	
	attr := &Attribute{Name: htmlAttrName}
	
	p.nextToken() // consume attribute name
	
	// Check if there's a value
	if p.currentToken.Type == TOKEN_EQUALS {
		p.nextToken() // consume '='
		
		if p.currentToken.Type == TOKEN_ATTR_VALUE {
			attr.Value = p.currentToken.Value
			p.nextToken() // consume attribute value
		} else {
			return nil, fmt.Errorf("expected attribute value, got %s", p.currentToken)
		}
	}
	
	return attr, nil
}
