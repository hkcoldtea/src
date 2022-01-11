package controller

import (
	"net/http"
	"path"
	"strconv"
	"strings"

	"config"
	"model"
	"font"
	"inventory"

	"github.com/hkcoldtea/src/projects/web.go/server"
)

func cssHandler(c *server.Context) {
	family := c.GetAttribute("family")
	if family == "" {
		// TODO: Add logging.
		BadRequest(c)
		return
	}
	sFormat := c.GetAttribute("format")
	if sFormat == "" {
		// Font format not specified,
		// default to WOFF.
		sFormat = "woff"
	}
	sDisplay := c.GetAttribute("display")
	if sDisplay != "swap" {
		sDisplay = ""
	}
	format := font.NOF
	format.FromString(sFormat)
	if format == font.NOF {
		// TODO: Add logging.
		BadRequest(c)
		return
	}
	queries := Queries(family, format)
	if len(queries) == 0 {
		// TODO: Add logging.
		BadRequest(c)
		return
	}

	cfg := config.Config()

	var templateData []*model.FontFace
	var templateName string
	for _, query := range queries {
		i := model.GetFontInventory()
		fnt := i.Query(*query)
		if fnt == nil {
			// TODO: Add logging.
			BadRequest(c)
			return
		}
		fontFace := new(model.FontFace)
		fontFace.SetBaseUrl(cfg.BaseUrl)
		fontFace.SetDisplay(sDisplay)
		fontFace.FromFont(*fnt)

		templateData = append(templateData, fontFace)
		templateName = fnt.Format.String() + ".css.tmpl"
	}
	maxAge := strconv.FormatUint(cfg.Flags.CcMaxAge, 10)
	c.SetHeader("Cache-Control", "max-age="+maxAge)
	c.SetHeader("Content-Type", "text/css; charset=utf-8")
	c.HTML(http.StatusOK, templateName, templateData)
}

func fontsHandler(c *server.Context) {
	URLPath := c.Param("filepath")
	sFormat := path.Ext(URLPath)
	if sFormat != "" {
		sFormat = sFormat[1:]
	}
	format := font.NOF
	format.FromString(sFormat)
	if format == font.NOF {
		// TODO: Add logging.
		BadRequest(c)
		return
	}
	family := c.Param("family")

	var f font.Font
	f.Format = format
	cfg := config.Config()
	maxAge := strconv.FormatUint(cfg.Flags.CcMaxAge, 10)
	c.SetHeader("Cache-Control", "max-age="+maxAge)
	c.SetHeader("Content-Type", f.MimeType())
	http.ServeFile(c.Writer, c.Req, path.Join("fonts", family, URLPath))
}

// Queries builds and returns a slice of pointers to inventory queries from the
// given family form value and font format format.
func Queries(family string, format font.Format) []*inventory.Query {
	var queries []*inventory.Query
	if strings.HasPrefix(family, "|") || strings.HasPrefix(family, ":") ||
		strings.HasPrefix(family, ",") {
		return queries
	}
	families := strings.Split(family, "|")
	for _, f := range families {
		// Contains font family name at index 0
		// and specified styles at index 1.
		familyStyles := strings.Split(f, ":")
		switch len(familyStyles) {
		case 1:
			q := new(inventory.Query)
			q.RowKey = familyStyles[0]
			// Weight and style are not specified, default
			// weight to 400 and style to normal.
			q.ColumnKey = "400normal"
			if q.RowKey == "" || q.ColumnKey == "" {
				continue
			}
			q.ColumnKey = format.String() + q.ColumnKey
			queries = append(queries, q)
		case 2:
			styles := strings.Split(familyStyles[1], ",")
			for _, s := range styles {
				q := new(inventory.Query)
				q.RowKey = familyStyles[0]
				if _, err := strconv.Atoi(s); err == nil {
					// Only weight is specified,
					// default style to normal.
					q.ColumnKey = s + "normal"
				} else {
					q.ColumnKey = s
				}
				if q.RowKey == "" || q.ColumnKey == "" {
					continue
				}
				q.ColumnKey = format.String() + q.ColumnKey
				queries = append(queries, q)
			}
		}
	}
	return queries
}
