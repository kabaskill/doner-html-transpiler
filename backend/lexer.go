package main

import (
	"fmt"
	"strings"
	"unicode"
)

// Security constants
const (
	MAX_INPUT_SIZE  = 100 * 1024 // 1MB limit
	MAX_TOKEN_LENGTH = 1000       // Prevent extremely long tokens
)

// TokenType represents the type of token
type TokenType int

const (
	TOKEN_UNKNOWN TokenType = iota
	TOKEN_TAG_OPEN           // <
	TOKEN_TAG_CLOSE          // >
	TOKEN_TAG_CLOSE_SLASH    // />
	TOKEN_TAG_END            // </
	TOKEN_TEXT               // plain text content
	TOKEN_TAG_NAME           // tag name
	TOKEN_ATTR_NAME          // attribute name
	TOKEN_ATTR_VALUE         // attribute value
	TOKEN_EQUALS             // =
	TOKEN_QUOTE              // " or '
	TOKEN_EOF                // end of file
)

// Token represents a lexical token
type Token struct {
	Type     TokenType
	Value    string
	Position int
}

// Lexer tokenizes German HTML input
type Lexer struct {
	input           []rune  // Use rune slice for proper UTF-8 handling
	position        int
	current         rune
	insideTag       bool
	afterTagName    bool    // Track if we just read a tag name
	afterEquals     bool    // Track if we just read an equals sign
}

// NewLexer creates a new lexer instance
func NewLexer(input string) *Lexer {
	runes := []rune(input)
	l := &Lexer{input: runes, insideTag: false, afterTagName: false, afterEquals: false}
	l.readChar()
	return l
}

// readChar reads the next character and advances position
func (l *Lexer) readChar() {
	if l.position >= len(l.input) {
		l.current = 0 // EOF
	} else {
		l.current = l.input[l.position]
	}
	l.position++
}

// peekChar returns the next character without advancing position
func (l *Lexer) peekChar() rune {
	if l.position >= len(l.input) {
		return 0
	}
	return l.input[l.position]
}

// skipWhitespace skips whitespace characters
func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(l.current) {
		l.readChar()
	}
}

// readString reads a quoted string
func (l *Lexer) readString(quote rune) string {
	var result []rune
	for {
		l.readChar()
		if l.current == quote || l.current == 0 {
			break
		}
		result = append(result, l.current)
	}
	return string(result)
}

// readIdentifier reads an identifier (tag name or attribute name)
func (l *Lexer) readIdentifier() string {
	position := l.position - 1
	for unicode.IsLetter(l.current) || unicode.IsDigit(l.current) || l.current == '_' || l.current == '-' {
		l.readChar()
	}
	return string(l.input[position : l.position-1])
}

// readText reads plain text content until a '<' is encountered
func (l *Lexer) readText() string {
	position := l.position - 1
	originalPos := position
	
	for l.current != '<' && l.current != 0 {
		l.readChar()
		// Security: Break large text into smaller chunks
		if (l.position - originalPos) > MAX_TOKEN_LENGTH {
			break
		}
	}
	
	text := l.input[position : l.position-1]
	l.position-- // Go back one position so we don't skip the '<'
	l.readChar()
	return strings.TrimSpace(string(text))
}

// readUnquotedValue reads an unquoted attribute value
func (l *Lexer) readUnquotedValue() string {
	position := l.position - 1
	// Read until whitespace or tag end characters
	for l.current != 0 && l.current != '>' && l.current != '/' && !unicode.IsSpace(l.current) {
		l.readChar()
	}
	return string(l.input[position : l.position-1])
}

// NextToken returns the next token from the input
func (l *Lexer) NextToken() Token {
	var tok Token
	
	// If we're not inside a tag and we encounter text content
	if !l.insideTag && l.current != '<' && l.current != 0 {
		tok.Type = TOKEN_TEXT
		tok.Position = l.position - 1
		tok.Value = l.readText()
		return tok
	}

	l.skipWhitespace()

	switch l.current {
	case '<':
		if l.peekChar() == '/' {
			l.readChar() // consume '/'
			l.readChar() // move to next char
			l.insideTag = true
			l.afterTagName = false
			l.afterEquals = false
			tok = Token{Type: TOKEN_TAG_END, Value: "</", Position: l.position - 2}
		} else {
			l.insideTag = true
			l.afterTagName = false
			l.afterEquals = false
			tok = Token{Type: TOKEN_TAG_OPEN, Value: "<", Position: l.position - 1}
			l.readChar()
		}
	case '>':
		l.insideTag = false
		l.afterTagName = false
		l.afterEquals = false
		tok = Token{Type: TOKEN_TAG_CLOSE, Value: ">", Position: l.position - 1}
		l.readChar()
	case '/':
		if l.peekChar() == '>' {
			l.readChar() // consume '>'
			l.readChar() // move to next char
			l.insideTag = false
			l.afterTagName = false
			l.afterEquals = false
			tok = Token{Type: TOKEN_TAG_CLOSE_SLASH, Value: "/>", Position: l.position - 2}
		} else {
			tok = Token{Type: TOKEN_UNKNOWN, Value: string(l.current), Position: l.position - 1}
			l.readChar()
		}
	case '=':
		l.afterEquals = true
		tok = Token{Type: TOKEN_EQUALS, Value: "=", Position: l.position - 1}
		l.readChar()
	case '"':
		l.afterEquals = false
		tok.Type = TOKEN_ATTR_VALUE
		tok.Position = l.position - 1
		tok.Value = l.readString('"')
		l.readChar() // consume closing quote
	case '\'':
		l.afterEquals = false
		tok.Type = TOKEN_ATTR_VALUE
		tok.Position = l.position - 1
		tok.Value = l.readString('\'')
		l.readChar() // consume closing quote
	case 0:
		tok = Token{Type: TOKEN_EOF, Value: "", Position: l.position}
	default:
		if unicode.IsLetter(l.current) {
			if l.afterEquals {
				// This is an unquoted attribute value
				tok.Type = TOKEN_ATTR_VALUE
				l.afterEquals = false
			} else if l.afterTagName {
				// This is an attribute name
				tok.Type = TOKEN_ATTR_NAME
			} else {
				// This is a tag name
				tok.Type = TOKEN_TAG_NAME
				l.afterTagName = true
			}
			tok.Position = l.position - 1
			
			if tok.Type == TOKEN_ATTR_VALUE {
				tok.Value = l.readUnquotedValue()
			} else {
				tok.Value = l.readIdentifier()
			}
		} else {
			tok = Token{Type: TOKEN_UNKNOWN, Value: string(l.current), Position: l.position - 1}
			l.readChar()
		}
	}

	return tok
}

// TokenTypeString returns a string representation of the token type
func (t TokenType) String() string {
	switch t {
	case TOKEN_TAG_OPEN:
		return "TAG_OPEN"
	case TOKEN_TAG_CLOSE:
		return "TAG_CLOSE"
	case TOKEN_TAG_CLOSE_SLASH:
		return "TAG_CLOSE_SLASH"
	case TOKEN_TAG_END:
		return "TAG_END"
	case TOKEN_TEXT:
		return "TEXT"
	case TOKEN_TAG_NAME:
		return "TAG_NAME"
	case TOKEN_ATTR_NAME:
		return "ATTR_NAME"
	case TOKEN_ATTR_VALUE:
		return "ATTR_VALUE"
	case TOKEN_EQUALS:
		return "EQUALS"
	case TOKEN_QUOTE:
		return "QUOTE"
	case TOKEN_EOF:
		return "EOF"
	default:
		return "UNKNOWN"
	}
}

// String returns a string representation of the token
func (t Token) String() string {
	return fmt.Sprintf("Token{Type: %s, Value: %q, Position: %d}", t.Type, t.Value, t.Position)
}

// Tokenize converts the input string into a slice of tokens with security checks
func Tokenize(input string) ([]Token, error) {
	// Security: Check input size to prevent DoS attacks
	if len(input) > MAX_INPUT_SIZE {
		return nil, fmt.Errorf("input too large: %d bytes exceeds limit of %d", 
			len(input), MAX_INPUT_SIZE)
	}
	
	lexer := NewLexer(input)
	var tokens []Token
	
	for {
		token := lexer.NextToken()
		
		// Security: Check token value length to prevent memory exhaustion
		if len(token.Value) > MAX_TOKEN_LENGTH {
			return nil, fmt.Errorf("token too long: %d characters exceeds limit of %d at position %d", 
				len(token.Value), MAX_TOKEN_LENGTH, token.Position)
		}
		
		tokens = append(tokens, token)
		
		if token.Type == TOKEN_EOF {
			break
		}
		
		// Security: Prevent infinite loops by checking token count
		if len(tokens) > MAX_INPUT_SIZE/10 { // Reasonable token count limit
			return nil, fmt.Errorf("too many tokens: %d exceeds safety limit", len(tokens))
		}
	}
	
	return tokens, nil
}

// SecureNewLexer creates a new lexer with input validation
func SecureNewLexer(input string) (*Lexer, error) {
	if len(input) > MAX_INPUT_SIZE {
		return nil, fmt.Errorf("input too large: %d bytes exceeds limit of %d", 
			len(input), MAX_INPUT_SIZE)
	}
	return NewLexer(input), nil
}
