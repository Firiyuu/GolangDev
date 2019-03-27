package system

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"pinehq/models"
)

var tmpl *template.Template

//LoadTemplates loads templates from views directory
func LoadTemplates() {
	tmpl = template.New("").Funcs(template.FuncMap{
		"isActiveLink":        isActiveLink,
		"stringInSlice":       stringInSlice,
		"formatDateTime":      formatDateTime,
		"now":                 now,
		"activeUserEmail":     activeUserEmail,
		"activeUserName":      activeUserName,
		"activeUserID":        activeUserID,
		"isUserAuthenticated": isUserAuthenticated,
		"signUpEnabled":       signUpEnabled,
		"noescape":            noescape,

	})
	fn := func(path string, f os.FileInfo, err error) error {
		if f.IsDir() != true && strings.HasSuffix(f.Name(), ".gohtml") {
			var err error
			tmpl, err = tmpl.ParseFiles(path)
			if err != nil {
				return err
			}
		}
		return nil
	}

	if err := filepath.Walk("views", fn); err != nil {
		panic(err)
	}
}

//GetTemplates returns preloaded templates
func GetTemplates() *template.Template {
	return tmpl
}

//isActiveLink checks uri against currently active (uri, or nil) and returns "active" if they are equal
func isActiveLink(c *gin.Context, uri string) string {
	if c != nil && c.Request.RequestURI == uri {
		return "active"
	}
	return ""
}

//formatDateTime prints timestamp in human format
func formatDateTime(t time.Time) string {
	return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
}

//stringInSlice returns true if value is in list slice
func stringInSlice(value string, list []string) bool {
	for i := range list {
		if value == list[i] {
			return true
		}
	}
	return false
}

//now returns current timestamp
func now() time.Time {
	return time.Now()
}

//activeUserEmail returns currently authenticated user email
func activeUserEmail(c *gin.Context) string {
	if c != nil {
		u, _ := c.Get("User")
		if user, ok := u.(*models.User); ok {
			return user.Email
		}
	}
	return ""
}

//activeUserName returns currently authenticated user name
func activeUserName(c *gin.Context) string {
	if c != nil {
		u, _ := c.Get("User")
		if user, ok := u.(*models.User); ok {
			return user.Name
		}
	}
	return ""
}

//activeUserID returns currently authenticated user ID
func activeUserID(c *gin.Context) uint64 {
	if c != nil {
		u, _ := c.Get("User")
		if user, ok := u.(*models.User); ok {
			return user.ID
		}
	}
	return 0
}

//isUserAuthenticated returns true is user is authenticated
func isUserAuthenticated(c *gin.Context) bool {
	if c != nil {
		u, _ := c.Get("User")
		if _, ok := u.(*models.User); ok {
			return true
		}
	}
	return false
}

//signUpEnabled returns true if sign up is enabled by config
func signUpEnabled(c *gin.Context) bool {
	if c != nil {
		se, _ := c.Get("SignupEnabled")
		if enabled, ok := se.(bool); ok {
			return enabled
		}
	}
	return false
}

//noescape unescapes html content
func noescape(content string) template.HTML {
	return template.HTML(content)
}

