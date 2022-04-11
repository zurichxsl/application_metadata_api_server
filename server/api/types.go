package api

type Maintainer struct {
	Name string `json:"name" validate:"required"`

	Email string `json:"email" validate:"required,email"`
}

type Release struct {
	Name string `json:"name,omitempty"`

	Comment string `json:"comment,omitempty"`

	Author Maintainer `json:"author,omitempty"`
}

type App struct {
	Id Id `json:"id,omitempty"`

	Title string `json:"title" validate:"required"`

	Version string `json:"version" validate:"required"`

	Maintainers []Maintainer `json:"maintainers" validate:"required"`

	Company string `json:"company" validate:"required"`

	Website string `json:"website" validate:"required"`

	Source string `json:"source" validate:"required"`

	License string `json:"license" validate:"required"`

	Description string `json:"description"`

	Labels map[string]string `json:"labels,omitempty"`

	Release Release `json:"release,omitempty"`
}

type Id string

func (a App) ID() Id {
	return a.Id
}
