package config


// Flags holds the state of the HTTP flags.
type Flags struct {
	AcAllowOrigin  bool   // Cross-origin resource sharing toggle.
	CcMaxAge       uint64 // Cache-Control max-age value.
	Etag           bool   // Entity tags validation toggle.
	Gzip           bool   // Response gzip compression toggle.
	Version        string // Server version string.
}

type SiteConfig struct {
	Url            string
	Flags          Flags
	BaseUrl        string
}

func Config() *SiteConfig {
	cfg := &SiteConfig{
		Url:            ":8080",
		BaseUrl:        "http://www.example.com/",
	}
	return cfg
}
