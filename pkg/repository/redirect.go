package repository

type Redirect struct {
	Original     string `json:"original"`
	Substitution string `json:"substitution"`
}
