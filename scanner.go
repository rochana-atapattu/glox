package main

import(
	"strconv"
	"fmt"
)


type Scanner struct{
	source string
	tokens []*Token

	start int
	current int
	line int

	keywords map[string]TokenType
}

func NewScanner(source string) *Scanner{
	return &Scanner{
		source: source,
		tokens: []*Token{},
		start: 0,
		current: 0,
		line: 1,
		keywords: map[string]TokenType{
   		     "and":    AND,
   		     "class":  CLASS,
   		     "else":   ELSE,
   		     "false":  FALSE,
   		     "for":    FOR,
   		     "fun":    FUN,
   		     "if":     IF,
   		     "nil":    NIL,
   		     "or":     OR,
   		     "print":  PRINT,
   		     "return": RETURN,
   		     "super":  SUPER,
   		     "this":   THIS,
   		     "true":   TRUE,
   		     "var":    VAR,
   		     "while":  WHILE,
   		 },
	}
}


func (s *Scanner) scanTokens() []*Token {
	
	for{
		if (s.isAtEnd()){
			break
		}
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append( s.tokens, NewToken( EOF, "", nil, s.line))
	return s.tokens
}

func (s *Scanner) isAtEnd() bool{
	return s.current >= len(s.source)
}

func (s *Scanner) scanToken() {

	switch c := s.advance(); c{
	case '(': 
		s.addToken(LEFT_PAREN)
        case ')': 
		s.addToken(RIGHT_PAREN)
        case '{': 
		s.addToken(LEFT_BRACE)
        case '}': 
		s.addToken(RIGHT_BRACE)
        case ',': 
		s.addToken(COMMA)
        case '.': 
		s.addToken(DOT)
        case '-': 
		s.addToken(MINUS)
        case '+': 
		s.addToken(PLUS)
        case ';': 
		s.addToken(SEMICOLON)
        case '*': 
		s.addToken(STAR)
	case '!':
		if s.match('=')	{
			s.addToken(BANG_EQUAL)
		} else{
			s.addToken(BANG)
		}
	case '=':
		if s.match('=')	{
			s.addToken(EQUAL_EQUAL)
		} else{
			s.addToken(EQUAL)
		}
	case '<':
		if s.match('=')	{
			s.addToken(LESS_EQUAL)
		} else{
			s.addToken(LESS)
		}
	case '>':
		if s.match('=')	{
			s.addToken(GREATER_EQUAL)
		} else{
			s.addToken(GREATER)
		}
	case '/':
		if s.match('/'){
			for {
				if s.peek() == '\n' && s.isAtEnd(){
					break
				}
				s.advance()
			}
		} else{
			s.addToken(SLASH)
		}
	case ' ':
	case '\r':
	case '\t':
	case '\n':
		s.line++
	case '"':
		s.string()
	default:
		if s.isDigit(c){
			s.number()
		} else if s.isAlpha(c){
			s.identifier()
		} else{
			error(s.line, "Unexpected character.")	
		}
	}
}

func (s *Scanner) advance() rune{
	c := rune(s.source[s.current])
	s.current++
	return c
}

func (s *Scanner) addToken(tokenType TokenType){
	s.addTokenWithLiteral(tokenType,nil)
}

func (s *Scanner) addTokenWithLiteral(tokenType TokenType, literal any){
	text := s.source[s.start:s.current]
	s.tokens = append (s.tokens, NewToken(tokenType, text, literal, s.line))
}

func (s *Scanner) match (expected rune) bool{
	if (s.isAtEnd()){
		return false
	}
	if rune(s.source[s.current]) != expected{
		return false
	}
	s.current++
	return true

}

func (s *Scanner) peek() rune{
	if s.isAtEnd(){
		return '\x00'
	}
	return rune(s.source[s.current])
}

func (s *Scanner) string(){
	for{
		if s.peek() != '"' && !s.isAtEnd(){
			if s.peek() == '\n'{
				s.line++
			}
			s.advance()
		}
	}
	if s.isAtEnd(){	
		error(s.line, "Unterminated string.")
		return
	}

	s.advance()
	
	text := s.source[s.start + 1:s.current - 1]
	s.addTokenWithLiteral(STRING, text)
}

func (s *Scanner) isDigit(c rune) bool{
	return c >= '0' && c <= '9'
}

func (s *Scanner) number() {
    for s.isDigit(s.peek()) {
        s.advance()
    }

    // Look for a fractional part.
    if s.peek() == '.' && s.isDigit(s.peekNext()) {
        // Consume the "."
        s.advance()

        for s.isDigit(s.peek()) {
            s.advance()
        }
    }

    value, err := strconv.ParseFloat(s.source[s.start:s.current], 64)
    if err != nil {
        // Handle parsing error if necessary
        fmt.Println("Error parsing number:", err)
    }

    s.addTokenWithLiteral(NUMBER, value)
}

func (s *Scanner) peekNext() rune {
    if s.current+1 >= len(s.source) {
        return '\x00' // Null character equivalent in Go
    }
    return rune(s.source[s.current+1])
}
func (s *Scanner) identifier() {
    	for s.isAlphaNumeric(s.peek()) {
        	s.advance()
    	}

    	text := s.source[s.start:s.current]
	tokenType, found := s.keywords[text]
	if !found {
    		tokenType = IDENTIFIER
	}
	s.addToken(tokenType)
}

func (s *Scanner)isAlpha(c rune) bool {
    return (c >= 'a' && c <= 'z') ||
           (c >= 'A' && c <= 'Z') ||
           c == '_'
}

func (s *Scanner)isAlphaNumeric(c rune) bool {
    return s.isAlpha(c) || s.isDigit(c)
}

