package models

type Spec struct {
	Path string `json:"path"`

	// other fields... :)
}

func (s Spec) GetPath() string {
	return s.Path
}
