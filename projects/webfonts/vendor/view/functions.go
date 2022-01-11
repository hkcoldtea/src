package view

import (
	"fmt"
	"html"
	"html/template"
	"strings"
	"time"

	"config"

	"github.com/hkcoldtea/src/projects/web.go/util"
)

var (
	// The functions available for use in the templates.
	TemplateFuncs = map[string]interface{}{
		"set": func(viewArgs map[string]interface{}, key string, value interface{}) template.JS {
			viewArgs[key] = value
			return template.JS("")
		},
		"append": func(viewArgs map[string]interface{}, key string, value interface{}) template.JS {
			if viewArgs[key] == nil {
				viewArgs[key] = []interface{}{value}
			} else {
				viewArgs[key] = append(viewArgs[key].([]interface{}), value)
			}
			return template.JS("")
		},

		"firstof": func(args ...interface{}) interface{} {
			for _, val := range args {
				switch val.(type) {
				case nil:
					continue
				case string:
					if val == "" {
						continue
					}
					return val
				default:
					return val
				}
			}
			return nil
		},

		"radio": func(name, f, val string) template.HTML {
			checked := ""
			if f == val {
				checked = " checked"
			}
			return template.HTML(fmt.Sprintf(`<input type="radio" name="%s" value="%s"%s>`,
				html.EscapeString(name), html.EscapeString(val), checked))
		},

		"checkbox": func(name, f, val string) template.HTML {
			checked := ""
			if f == val {
				checked = " checked"
			}
			return template.HTML(fmt.Sprintf(`<input type="checkbox" name="%s" value="%s"%s>`,
				html.EscapeString(name), html.EscapeString(val), checked))
		},

		// Pads the given string with &nbsp;'s up to the given width.
		"pad": func(str string, width int) template.HTML {
			if len(str) >= width {
				return template.HTML(html.EscapeString(str))
			}
			return template.HTML(html.EscapeString(str) + strings.Repeat("&nbsp;", width-len(str)))
		},

		// Replaces newlines with <br>
		"nl2br": func(text string) template.HTML {
			return template.HTML(strings.Replace(template.HTMLEscapeString(text), "\n", "<br>", -1))
		},

		// Skips sanitation on the parameter. Do not use with dynamic data.
		"raw": func(text string) template.HTML {
			return template.HTML(text)
		},

		"noescape": func(s string) template.HTML {
			return template.HTML(s)
		},

		"htmlSafe": func(html string) template.HTML {
			return template.HTML(html)
		},

		"datetime": func(date time.Time) string {
			return date.Format(time.RFC1123)
		},

		"even": func(a int) bool { return (a % 2) == 0 },
		"formatAsDate": util.FormatAsDate,
		"trim":         strings.TrimSpace,
		"Upper": func(s string) string {
			return strings.ToUpper(s)
		},
	}
)


/////////////////////
// Template functions
/////////////////////

func SiteURL() string {
	return config.Config().BaseUrl
}
