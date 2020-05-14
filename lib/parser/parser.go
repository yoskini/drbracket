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

type BracketKind rune

const (
	OpenRound    BracketKind = '('
	OpenSquare               = '['
	OpenBrace                = '{'
	ClosedRound              = ')'
	ClosedSquare             = ']'
	ClosedBrace              = '}'
)

type Bracket struct {
	Kind BracketKind
	Line int
	Col  int
}

type BracketParser struct {
	stack []Bracket
}

func NewBracketParser() *BracketParser {
	return &BracketParser{}
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

func (p *BracketParser) ParseLine(lineNum int, line string) error {
	for col, c := range line {
		switch c {
		case '(':
			p.stack = append(p.stack, Bracket{Kind: OpenRound, Line: lineNum, Col: col + 1})
		case '[':
			p.stack = append(p.stack, Bracket{Kind: OpenSquare, Line: lineNum, Col: col + 1})
		case '{':
			p.stack = append(p.stack, Bracket{Kind: OpenBrace, Line: lineNum, Col: col + 1})
		case ')':
			if len(p.stack) == 0 || p.stack[len(p.stack)-1].Kind != OpenRound {
				return fmt.Errorf("Unbalanced bracket. Found ')' at line %d, col %d", lineNum, col)
			}
			n := len(p.stack) - 1
			p.stack = p.stack[:n]
		case ']':
			if len(p.stack) == 0 || p.stack[len(p.stack)-1].Kind != OpenSquare {
				return fmt.Errorf("Unbalanced bracket. Found ']' at line %d, col %d", lineNum, col)
			}
			n := len(p.stack) - 1
			p.stack = p.stack[:n]
		case '}':
			if len(p.stack) == 0 || p.stack[len(p.stack)-1].Kind != OpenBrace {
				return fmt.Errorf("Unbalanced bracket. Found '}' at line %d, col %d", lineNum, col)
			}
			n := len(p.stack) - 1
			p.stack = p.stack[:n]
		default:
		}
	}
	return nil
}
