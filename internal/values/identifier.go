package values

type Identifier struct {
	identifier string
}

func IdentifierFromString(identifier string) Identifier {
	return Identifier{
		identifier: identifier,
	}
}

func (i Identifier) String() string {
	return i.identifier
}
