package types

type FileVo struct {
	Name string `json:"name,omitempty"`
	Url  string `json:"url,omitempty"`
	Tag  string `json:"tag,omitempty"`
	Ext  string `json:"ext,omitempty"`
	Size int64  `json:"size,omitempty"`
}

type FileDto struct {
	Url    string `json:"url"`
	Uid    string `json:"uid,omitempty"`
	Name   string `json:"name"`
	Height int64  `json:"height,omitempty"`
	Width  int64  `json:"width,omitempty"`
	Size   int64  `json:"size,omitempty"`
	Tag    string `json:"tag,omitempty"`
}
