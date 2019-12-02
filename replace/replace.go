package replace

import (
	"regexp"
	"strings"
	"sync"

	"github.com/jiuker/ycss/css"
)

type Replace interface {
	GetFileBody() string
	GetOutFileBody() string
	FindClass([]*regexp.Regexp) []string
	GetRegexpCss([]string, *sync.Map, css.Css) *string
	// needZoomKey for rn
	Zoom(css *string, unit string, needZoomUint string, needZoomKey []string, zoom float64) *string
	GetOldCss(*regexp.Regexp) (*string, *string, error)
	Replace(old *string, new *string, pos *string) *string
	Save(newPos *string, old *string) error
	// close file
	Done()
}

var sameTestCh = `1234567890-=qwertyuiop[]asdfghjkl;'zxcvbnm,./!@#$%^&*()_+~'`

func IsSame(a *string, b *string) bool {
	for _, v := range sameTestCh {
		if strings.Count(*a, string(v)) != strings.Count(*b, string(v)) {
			return false
		}
	}
	return true
}

// find All css
func findAllCss(body string, reg []*regexp.Regexp) []string {
	var cssMap = map[string]bool{}
	var cssArray []string
	// make class in order
	for _, r := range reg {
		allClass := r.FindAllStringSubmatch(body, -1)
		for _, v1 := range allClass {
			if len(v1) != 2 {
				continue
			}
			v1[1] = strings.Replace(v1[1], "	", " ", -1)
			classSplit := strings.SplitN(v1[1], " ", -1)
			for _, v2 := range classSplit {
				if v2 != "" {
					cssArray = append(cssArray, v2)
				}
			}
		}
	}
	var ret []string
	for _, v1 := range cssArray {
		if !cssMap[v1] {
			ret = append(ret, v1)
			cssMap[v1] = true
		}
	}
	return ret
}
