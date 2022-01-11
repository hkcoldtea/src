package model

import (
	"fmt"

	"font"
)

// FontFace represents a single @font-face CSS rule.
type FontFace struct {
	Base64Data string
	Family     string
	Format     string
	MimeType   string
	Style      string
	Weight     int
	Url        string
	BaseUrl    string
	Display    bool
	Range      string
}

// FromFont initializes ff from the given font. If it is called on an already
// initialized font face, changes the font face according to the given font.
// Returns an error if the initialization fails.
func (ff *FontFace) FromFont(f font.Font) {
	mimeType := f.MimeType()
	if mimeType == "" {
		// Should not happen.
		fmt.Errorf("can't determine font MIME type")
	}
	ff.Url = ff.BaseUrl + f.Path
	ff.Family = f.Family
	ff.Format = f.Format.String()
	ff.MimeType = mimeType
	ff.Style = f.Style
	ff.Weight = f.Weight
	ff.Range = f.Range
}

func (ff *FontFace) SetBaseUrl(u string) {
	ff.BaseUrl = u
}

func (ff *FontFace) SetDisplay(u string) {
	if u == "" {
		ff.Display = false
		return
	}
	ff.Display = true
}
