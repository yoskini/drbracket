// Copyright (C) 2020 Fabio Del Vigna
//
// This file is part of drbracket.
//
// drbracket is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// drbracket is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with drbracket.  If not, see <http://www.gnu.org/licenses/>.

package parser

import (
	"fmt"
)

const (
	BracketOpenRound    = '('
	BracketOpenSquare   = '['
	BracketOpenBrace    = '{'
	BracketOpenAngular  = '<'
	BracketClosedRound  = ')'
	BracketClosedSquare = ']'
	BracketClosedBrace  = '}'
	BracketCloseAngular = '>'
)

func expectedOpen(b rune) rune {
	switch b {
	case BracketClosedRound:
		return BracketOpenRound
	case BracketClosedSquare:
		return BracketOpenSquare
	case BracketClosedBrace:
		return BracketOpenBrace
	case BracketCloseAngular:
		return BracketOpenAngular
	}
	panic("Unknown bracket kind")
}

type Bracket struct {
	Kind rune
	Line int
	Col  int
}

type BracketParser struct {
	stack []Bracket
}

func NewBracketParser() *BracketParser {
	return &BracketParser{
		stack: make([]Bracket, 0, 100),
	}
}

func (p *BracketParser) Empty() bool {
	return len(p.stack) == 0
}

func (p *BracketParser) Top() *Bracket {
	if n := len(p.stack); n > 0 {
		return &p.stack[n-1]
	}
	return nil
}

func (p *BracketParser) Pop() *Bracket {
	if n := len(p.stack); n > 0 {
		ret := &p.stack[n-1]
		p.stack = p.stack[:n-1]
		return ret
	}
	return nil
}

func (p *BracketParser) Push(b Bracket) {
	p.stack = append(p.stack, b)
}

func bracketError(found rune, lineFound, colFound int, expected rune, lineExpected, colExpected int) error {
	return fmt.Errorf("Unbalanced bracket. Found %c at line: %d, col: %d. Expected %c from line: %d, col: %d",
		found, lineFound, colFound, expected, lineExpected, colExpected)
}

func (p *BracketParser) ParseLine(lineNum int, line string) error {
	for col, c := range line {
		switch c {
		case BracketOpenRound:
			fallthrough
		case BracketOpenSquare:
			fallthrough
		case BracketOpenBrace:
			//fallthrough
			//case BracketOpenAngular:
			p.Push(Bracket{Kind: c, Line: lineNum, Col: col + 1})
		case BracketClosedRound:
			fallthrough
		case BracketClosedSquare:
			fallthrough
		case BracketClosedBrace:
			//fallthrough
			//case BracketCloseAngular:
			if b := p.Top(); p != nil && b.Kind != expectedOpen(c) {
				return bracketError(c, lineNum, col+1, b.Kind, b.Line, b.Col)
			}
			_ = p.Pop()
		default:
		}
	}
	return nil
}
