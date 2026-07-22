package pypi

type Index struct {
	Meta     Meta      `json:"meta"`
	Projects []Project `json:"projects"`
}

type Meta struct {
	APIVersion string `json:"api-version"`
}

type Project struct {
	Name string `json:"name"`
}

type ProjectDetail struct {
	Meta     Meta     `json:"meta"`
	Name     string   `json:"name"`
	Files    []File   `json:"files"`
	Versions []string `json:"versions,omitempty"`
}

type File struct {
	Filename       string            `json:"filename"`
	URL            string            `json:"url"`
	Hashes         map[string]string `json:"hashes"`
	RequiresPython string            `json:"requires-python,omitempty"`
	Yanked         any               `json:"yanked,omitempty"` // bool or string reason
	Size           int64             `json:"size,omitempty"`
	UploadTime     string            `json:"upload-time,omitempty"`
}
