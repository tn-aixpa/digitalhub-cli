package models

type Artifact struct {
	Spec Spec `json:"spec"`
}

func (a Artifact) GetSpec() Spec {
	return a.Spec
}
