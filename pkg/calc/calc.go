/*
Парсер и калькулятор математических выражений.
Поддерживаются выражения в скобках и базовые бинарные операторы '+', '-', '/', '*'.
Поддерживаются операции над целыми и дробными числами.
*/
package calc

import (
	"bytes"
	"fmt"
	"strconv"
	"unicode"
)

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

var binOps = map[Token]struct{}{
	Plus:  {},
	Minus: {},
	Div:   {},
	Mul:   {},
}

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

func NewSolver(parser *Parser) (*Solver, error) {
	l, err := parser.getNextLexeme()
	if err != nil {
		return nil, err
	}
	return &Solver{
		parser:        parser,
		currentLexeme: l,
	}, nil
}

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
	return &Const{token: Int, value: int64(number)}, nil
}

func (p *Parser) getNextLexeme() (Lexeme, error) {
	switch r := p.currentRune; {
	case unicode.IsDigit(r):
		l, err := p.readNumber()
		if err != nil {
			return nil, fmt.Errorf("number parsing error: %w", err)
		}
		p.skipWhitespace()
		return l, nil
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

type Solver struct {
	parser        *Parser
	currentLexeme Lexeme
}

func (s *Solver) Term() (*genericConst, error) {
	result, err := s.Expr()
	if err != nil {
		return nil, err
	}
	if !s.isBinOp() && s.currentLexeme.Token() != EOF && s.currentLexeme.Token() != Rparen {
		v, _ := s.currentLexeme.Value()
		return nil, fmt.Errorf("binary operator or EOF expectd, got %v", v)
	}
	for s.currentLexeme.Token() == Plus || s.currentLexeme.Token() == Minus {
		switch s.currentLexeme.Token() {
		case Plus:
			err = s.consume(Plus)
			right, err := s.Expr()
			if err != nil {
				return nil, err
			}
			result = result.Add(right)
		case Minus:
			err = s.consume(Minus)
			right, err := s.Expr()
			if err != nil {
				return nil, err
			}
			result = result.Sub(right)
		}
	}
	return result, nil
}

func (s *Solver) consume(t Token) error {
	var err error
	if s.currentLexeme.Token() == t {
		s.currentLexeme, err = s.parser.getNextLexeme()
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Solver) Expr() (*genericConst, error) {
	result, err := s.Factor()
	if err != nil {
		return nil, err
	}
	if !s.isBinOp() && s.currentLexeme.Token() != EOF && s.currentLexeme.Token() != Rparen {
		v, _ := s.currentLexeme.Value()
		return nil, fmt.Errorf("binary operator or EOF expectd, got %v", v)
	}
	for s.currentLexeme.Token() == Mul || s.currentLexeme.Token() == Div {
		switch s.currentLexeme.Token() {
		case Mul:
			err = s.consume(Mul)
			v, err := s.Factor()
			if err != nil {
				return nil, err
			}
			result = result.Mul(v)
		case Div:
			err = s.consume(Div)
			v, err := s.Factor()
			if err != nil {
				return nil, err
			}
			result = result.Div(v)
		}
	}
	return result, err
}

func (s *Solver) isBinOp() bool {
	_, ok := binOps[s.currentLexeme.Token()]
	return ok
}

func (s *Solver) Factor() (*genericConst, error) {
	switch s.currentLexeme.Token() {
	case Lparen:
		err := s.consume(Lparen)
		if err != nil {
			return nil, err
		}
		r, err := s.Term()
		if err != nil {
			return nil, err
		}
		err = s.consume(Rparen)
		if err != nil {
			return nil, err
		}
		return r, nil
	case Int:
		v, err := s.currentLexeme.Value()
		if err != nil {
			return nil, err
		}
		vv := &genericConst{value: v, typ: gInt}
		s.currentLexeme, err = s.parser.getNextLexeme()
		return vv, err
	case Real:
		v, err := s.currentLexeme.Value()
		if err != nil {
			return nil, err
		}
		vv := &genericConst{value: v, typ: gFloat}
		s.currentLexeme, err = s.parser.getNextLexeme()
		return vv, err
	default:
		return nil, fmt.Errorf("unexpected token %v", s.currentLexeme.Token())
	}
}

func (s *Solver) Solve() (interface{}, error) {
	r, err := s.Term()
	if err != nil {
		return nil, err
	}
	if s.parser.currentRune != -1 {
		return nil, fmt.Errorf("parsing error: expected EOF got %s", string(s.parser.currentRune))
	}
	return r.Value(), nil
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
