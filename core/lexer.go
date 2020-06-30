package core

import "strings"

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

type keyword string

const (
	selectKeyword     keyword = "select"
	fromKeyword       keyword = "from"
	asKeyword         keyword = "as"
	tableKeyword      keyword = "table"
	createKeyword     keyword = "create"
	dropKeyword       keyword = "drop"
	insertKeyword     keyword = "insert"
	intoKeyword       keyword = "into"
	valuesKeyword     keyword = "values"
	intKeyword        keyword = "int"
	textKeyword       keyword = "text"
	boolKeyword       keyword = "boolean"
	whereKeyword      keyword = "where"
	andKeyword        keyword = "and"
	orKeyword         keyword = "or"
	trueKeyword       keyword = "true"
	falseKeyword      keyword = "false"
	uniqueKeyword     keyword = "unique"
	indexKeyword      keyword = "index"
	onKeyword         keyword = "on"
	primarykeyKeyword keyword = "primary key"
	nullKeyword       keyword = "null"
)

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
	keywordKind
	boolKind
	nullKind
	stringKind
	identifierKind
)

// LexParse parse input to a list of tokens.
func LexParse(input string) ([]*Token, error) {
	var tokens []*Token
	lexers := []lexerImpl{keywordLexer, symbolLexer, numberLexer}

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

func keywordLexer(source string, ic cursor) (*Token, cursor, bool) {
	cur := ic
	keywords := []keyword{
		selectKeyword,
		insertKeyword,
		valuesKeyword,
		tableKeyword,
		createKeyword,
		dropKeyword,
		whereKeyword,
		fromKeyword,
		intoKeyword,
		textKeyword,
		boolKeyword,
		intKeyword,
		andKeyword,
		orKeyword,
		asKeyword,
		trueKeyword,
		falseKeyword,
		uniqueKeyword,
		indexKeyword,
		onKeyword,
		primarykeyKeyword,
		nullKeyword,
	}

	var options []string
	for _, k := range keywords {
		options = append(options, string(k))
	}

	match := pickFrom(source, ic, options)
	if match == "" {
		return nil, ic, false
	}

	cur.pointer = ic.pointer + uint(len(match))
	cur.loc.col = ic.loc.col + uint(len(match))

	kind := keywordKind
	if match == string(trueKeyword) || match == string(falseKeyword) {
		kind = boolKind
	}

	if match == string(nullKeyword) {
		kind = nullKind
	}

	return &Token{
		value: match,
		kind:  kind,
		loc:   ic.loc,
	}, cur, true
}

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
	lowerSource := strings.ToLower(source)
	cur := ic

	for _, candidate := range candidates {
		canBe := true
		for _, candidateLetter := range candidate {
			currentLetter := lowerSource[cur.pointer]

			if rune(currentLetter) != candidateLetter {
				canBe = false
				break
			}

			cur.pointer++
		}

		if canBe {
			return candidate
		}

		cur.pointer = ic.pointer
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

// lexCharacterDelimited looks through a source string starting at the
// given cursor to find a start- and end- delimiter. The delimiter can
// be escaped be preceeding the delimiter with itself
func lexCharacterDelimited(source string, ic cursor, delimiter byte) (*Token, cursor, bool) {
	cur := ic

	if len(source[cur.pointer:]) == 0 {
		return nil, ic, false
	}

	if source[cur.pointer] != delimiter {
		return nil, ic, false
	}

	cur.loc.col++
	cur.pointer++

	var value []byte
	for ; cur.pointer < uint(len(source)); cur.pointer++ {
		c := source[cur.pointer]

		if c == delimiter {
			// SQL escapes are via double characters, not backslash.
			if cur.pointer+1 >= uint(len(source)) || source[cur.pointer+1] != delimiter {
				cur.pointer++
				cur.loc.col++
				return &Token{
					value: string(value),
					loc:   ic.loc,
					kind:  stringKind,
				}, cur, true
			}
			value = append(value, delimiter)
			cur.pointer++
			cur.loc.col++
		}

		value = append(value, c)
		cur.loc.col++
	}

	return nil, ic, false
}

func identifierLexer(source string, ic cursor) (*Token, cursor, bool) {
	// Handle separately if is a double-quoted identifier
	if token, newCursor, ok := lexCharacterDelimited(source, ic, '"'); ok {
		return token, newCursor, true
	}

	cur := ic

	c := source[cur.pointer]
	// Other characters count too, big ignoring non-ascii for now
	isAlphabetical := (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z')
	if !isAlphabetical {
		return nil, ic, false
	}
	cur.pointer++
	cur.loc.col++

	value := []byte{c}
	for ; cur.pointer < uint(len(source)); cur.pointer++ {
		c = source[cur.pointer]

		// Other characters count too, big ignoring non-ascii for now
		isAlphabetical := (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z')
		isNumeric := c >= '0' && c <= '9'
		if isAlphabetical || isNumeric || c == '$' || c == '_' {
			value = append(value, c)
			cur.loc.col++
			continue
		}

		break
	}

	if len(value) == 0 {
		return nil, ic, false
	}

	return &Token{
		// Unquoted identifiers are case-insensitive
		value: strings.ToLower(string(value)),
		loc:   ic.loc,
		kind:  identifierKind,
	}, cur, true
}
