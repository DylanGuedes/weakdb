package core

// location inside input source
type location struct {
	line uint
	col  uint
}

// current pointer of the lexer
type cursor struct {
	pointer uint
	loc     location
}

// Token represents an AST token.
type Token struct {
	value string
	kind  tokenKind
	loc   location
}

type symbol string

const (
	semicolonSymbol  symbol = ";"
	asteriskSymbol   symbol = "*"
	commaSymbol      symbol = ","
	leftParenSymbol  symbol = "("
	rightParenSymbol symbol = ")"
	eqSymbol         symbol = "="
	neqSymbol        symbol = "<>"
	neqSymbol2       symbol = "!="
	concatSymbol     symbol = "||"
	plusSymbol       symbol = "+"
	ltSymbol         symbol = "<"
	lteSymbol        symbol = "<="
	gtSymbol         symbol = ">"
	gteSymbol        symbol = ">="
)

type tokenKind uint

const (
	numericKind tokenKind = iota
	symbolKind
)

// LexParse parse input to a list of tokens.
func LexParse(input string) ([]*Token, error) {
	var tokens []*Token
	lexers := []lexerImpl{numberLexer}

	cur := cursor{}

	for cur.pointer < uint(len(input)) {
		for _, lexer := range lexers {
			token, newCursor, ok := lexer(input, cur)

			if ok {
				cur = newCursor
				if token != nil {
					tokens = append(tokens, token)
				}
			}
		}
	}

	return tokens, nil
}

type lexerImpl func(string, cursor) (*Token, cursor, bool)

func symbolLexer(source string, ic cursor) (*Token, cursor, bool) {
	c := source[ic.pointer]
	cur := ic
	// Will get overwritten later if not an ignored syntax
	cur.pointer++
	cur.loc.col++

	switch c {
	// Syntax that should be thrown away
	case '\n':
		cur.loc.line++
		cur.loc.col = 0
		fallthrough
	case '\t':
		fallthrough
	case ' ':
		return nil, cur, true
	}

	// Syntax that should be kept
	symbols := []symbol{
		eqSymbol,
		neqSymbol,
		neqSymbol2,
		ltSymbol,
		lteSymbol,
		gtSymbol,
		gteSymbol,
		concatSymbol,
		plusSymbol,
		commaSymbol,
		leftParenSymbol,
		rightParenSymbol,
		semicolonSymbol,
		asteriskSymbol,
	}

	var options []string
	for _, s := range symbols {
		options = append(options, string(s))
	}

	match := pickFrom(source, ic, options)
	// Unknown character
	if match == "" {
		return nil, ic, false
	}

	cur.pointer = ic.pointer + uint(len(match))
	cur.loc.col = ic.loc.col + uint(len(match))

	return &Token{
		value: match,
		loc:   ic.loc,
		kind:  symbolKind,
	}, cur, true
}

func pickFrom(source string, ic cursor, candidates []string) string {
	cur := ic

	for _, candidate := range candidates {
		currentLetter := source[cur.pointer]
		canBe := true
		for _, candidateLetter := range candidate {
			if rune(currentLetter) != candidateLetter {
				canBe = false
				break
			}

			cur.pointer++
		}

		if canBe {
			return candidate
		}
	}

	return ""
}

func numberLexer(input string, ptr cursor) (*Token, cursor, bool) {
	val := ""

	for _, letter := range input {
		isNumeric := letter >= '0' && letter <= '9'
		if !isNumeric {
			break
		}

		val = val + string(letter)

		ptr.pointer++
		ptr.loc.col++
	}

	if val == "" {
		return nil, ptr, false
	}

	return &Token{
		value: val,
		loc:   ptr.loc,
		kind:  numericKind,
	}, ptr, true
}
