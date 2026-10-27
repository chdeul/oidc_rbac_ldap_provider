package obj

type Project struct {
	Pid         string
	Identifiant string
	NomComplet  string
	Description string
}

func NewProject(id, ident, nom, desc string) *Project {
	return &Project{Pid: id, Identifiant: ident, NomComplet: nom, Description: desc}
}
