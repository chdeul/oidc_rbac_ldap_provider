package obj

type Right struct {
	Type  string   // User , Group
	Name  string   // Identifiant of user
	Tools []string // roles in each tools
}

func NewRight(typ, name string, tools []string) *Right {
	return &Right{Type: typ, Name: name, Tools: tools}
}
