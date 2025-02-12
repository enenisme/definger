package pkg

type Tags struct {
	Tags []Tag
}

// Tag 定义指纹的结构
type Tag struct {
	ID   string `json:"id"`
	Info Infos  `json:"info"`
	HTTP []HTTP `json:"http"`
}

// Infos 定义指纹的信息
type Infos struct {
	Name     string    `json:"name"`
	Author   string    `json:"author"`
	Tags     string    `json:"tags"`
	Severity string    `json:"severity"`
	Metadata Metadatas `json:"metadata"`
}

// Metadatas 定义指纹的元数据
type Metadatas struct {
	Product  string `json:"product"`
	Vendor   string `json:"vendor"`
	Verified bool   `json:"verified"`
}

// HTTP 定义指纹的HTTP请求
type HTTP struct {
	Method   string     `json:"method"`
	Path     []string   `json:"path"`
	Mode     string     `json:"mode,omitempty"`
	Matchers []Matchers `json:"matchers"`
}

// Matchers 定义指纹的匹配器
type Matchers struct {
	Type            string   `json:"type,omitempty"`
	Words           []string `json:"words,omitempty"`
	Part            string   `json:"part,omitempty"`
	Condition       string   `json:"condition,omitempty"`
	CaseInsensitive bool     `json:"case-insensitive,omitempty"`
	Hash            []string `json:"hash,omitempty"`
}
