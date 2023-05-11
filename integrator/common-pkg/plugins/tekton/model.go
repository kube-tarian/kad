package tekton

type tekton struct{}

// func NewTekton() (*tektonspec, error) {
// 	return &tektonspec{}, nil
// }

type Tektonspec struct {
	//	Path string
	//	Filepath string
	Hostname string
	Giturl   string
	Name     string
	Mail     string
}
