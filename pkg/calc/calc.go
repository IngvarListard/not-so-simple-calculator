package calc

import (
	"bytes"
	"fmt"
	"strconv"
	"unicode"
)

type Token int

func (t Token) Priority() int {
	switch t {
	case Plus, Minus:
		return 3
	case Mul, Div:
		return 2
	}
	return -1
}

const (
	EOF Token = iota
	Int
	Real
	Plus
	Minus
	Mul
	Div
	Lparen
	Rparen

	Null rune = -1
)

type Parser struct {
	text        []rune
	pos         int
	currentRune rune
}

func (p *Parser) next() {
	p.pos++
	if p.pos >= len(p.text) {
		p.currentRune = Null
		return
	}
	p.currentRune = p.text[p.pos]
}

func (p *Parser) readNumber() (Lexeme, error) {
	var numBuf bytes.Buffer
	for unicode.IsDigit(p.currentRune) {
		numBuf.WriteRune(p.currentRune)
		p.next()
	}

	if p.currentRune == '.' {
		numBuf.WriteRune(p.currentRune)
		p.next()

		for unicode.IsDigit(p.currentRune) {
			numBuf.WriteRune(p.currentRune)
			p.next()
		}

		realNumber, err := strconv.ParseFloat(numBuf.String(), 64)
		if err != nil {
			return nil, fmt.Errorf("float parsing error: %v", err)
		}
		return &Const{token: Real, value: realNumber}, nil
	}

	number, _ := strconv.Atoi(numBuf.String())
	p.skipWhitespace()
	return &Const{token: Int, value: int64(number)}, nil
}

func (p *Parser) getNextLexeme() (Lexeme, error) {
	switch r := p.currentRune; {
	case unicode.IsDigit(r):
		return p.readNumber()
	case r == '(':
		p.next()
		p.skipWhitespace()
		return &Const{token: Lparen, value: r}, nil
	case r == ')':
		p.next()
		p.skipWhitespace()
		return &Const{token: Rparen, value: r}, nil
	case r == '+':
		p.next()
		p.skipWhitespace()
		return &Const{token: Plus, value: r}, nil
	case r == '-':
		p.next()
		p.skipWhitespace()
		return &Const{token: Minus, value: r}, nil
	case r == '*':
		p.next()
		p.skipWhitespace()
		return &Const{token: Mul, value: r}, nil
	case r == '/':
		p.next()
		p.skipWhitespace()
		return &Const{token: Div, value: r}, nil
	case r == Null:
		p.next()
		p.skipWhitespace()
		return &Const{token: EOF, value: Null}, nil
	default:
		return nil, fmt.Errorf("unexpected symbol occurance: %s", string(r))
	}
}

func (p *Parser) skipWhitespace() {
	for unicode.IsSpace(p.currentRune) {
		p.next()
	}
}

type Interpreter struct {
	parser        *Parser
	currentLexeme Lexeme
}

func (i *Interpreter) Term() (*genericConst, error) {
	result, err := i.Expr()
	if err != nil {
		return nil, err
	}
	for i.currentLexeme.Token() == Plus || i.currentLexeme.Token() == Minus {
		switch i.currentLexeme.Token() {
		case Plus:
			err = i.consume(Plus)
			right, err := i.Expr()
			if err != nil {
				return nil, err
			}
			result = result.Add(right)
		case Minus:
			err = i.consume(Minus)
			right, err := i.Expr()
			if err != nil {
				return nil, err
			}
			result = result.Sub(right)
		}
	}
	return result, nil
}

func (i *Interpreter) consume(t Token) error {
	var err error
	if i.currentLexeme.Token() == t {
		i.currentLexeme, err = i.parser.getNextLexeme()
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *Interpreter) Expr() (*genericConst, error) {
	result, err := i.Factor()
	if err != nil {
		return nil, err
	}
	for i.currentLexeme.Token() == Mul || i.currentLexeme.Token() == Div {
		switch i.currentLexeme.Token() {
		case Mul:
			err = i.consume(Mul)
			v, err := i.Factor()
			if err != nil {
				return nil, err
			}
			result = result.Mul(v)
		case Div:
			err = i.consume(Div)
			v, err := i.Factor()
			if err != nil {
				return nil, err
			}
			result = result.Div(v)
		}
	}
	return result, err
}

func (i *Interpreter) Factor() (*genericConst, error) {

	switch i.currentLexeme.Token() {
	case Lparen:
		err := i.consume(Lparen)
		if err != nil {
			return nil, err
		}
		r, err := i.Term()
		if err != nil {
			return nil, err
		}
		err = i.consume(Rparen)
		if err != nil {
			return nil, err
		}
		return r, nil
	case Int:
		v, err := i.currentLexeme.Value()
		if err != nil {
			return nil, err
		}
		vv := &genericConst{value: v, typ: gInt}
		i.currentLexeme, err = i.parser.getNextLexeme()
		return vv, err
	default:
		return nil, fmt.Errorf("unexpected token %v", i.currentLexeme.Token())
	}

}

type Lexeme interface {
	Token() Token
	Value() (interface{}, error)
}

type Const struct {
	token Token
	value interface{}
}

func (c *Const) Token() Token {
	return c.token
}

func (c *Const) Value() (interface{}, error) {
	return c.value, nil
}

func NewParser(text string) (*Parser, error) {
	if len(text) == 0 {
		return nil, fmt.Errorf("text is empty")
	}

	runes := []rune(text)
	parser := &Parser{
		text:        runes,
		pos:         0,
		currentRune: runes[0],
	}
	parser.skipWhitespace()
	return parser, nil
}

func NewInterpreter(parser *Parser) (*Interpreter, error) {
	l, err := parser.getNextLexeme()
	if err != nil {
		return nil, err
	}
	return &Interpreter{
		parser:        parser,
		currentLexeme: l,
	}, nil
}
