package core

type astKind uint

const (
	selectKind astKind = iota
	createTableKind
)

type statement struct {
	Kind astKind
}

// Ast is the AST representation of the given command
type Ast struct {
	Statements []*statement
}
