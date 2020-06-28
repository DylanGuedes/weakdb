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

type tokenKind uint

const (
	numericKind tokenKind = iota
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
