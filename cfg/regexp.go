package cfg

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/jiuker/ycss/css"

	"github.com/PuerkitoBio/goquery"
)

type FileType int

func (c FileType) ToString() string {
	switch c {
	case RNCSS:
		return "RN"
	case VueCss:
		return "Vue"
	}
	return "Unknown"
}

const (
	RNCSS  = FileType(1)
	VueCss = FileType(2)
)

// parse it to css
type Regexp interface {
	Parse() error
	GetCommonReg() *sync.Map
	GetSingle() css.Css
}
type cssRegexp struct {
	rwLocker sync.RWMutex
	// key(reg):arr[reg1 reg2]
	common *sync.Map
	// reg1(reg):map[cssUint cssUint]
	single css.Css
}

func (r *cssRegexp) GetCommonReg() *sync.Map {
	return r.common
}

func (r *cssRegexp) GetSingle() css.Css {
	return r.single
}

func (r *cssRegexp) Parse() error {
	fmt.Println("start parse css Regexp!")
	_reg.rwLocker.Lock()
	defer _reg.rwLocker.Unlock()
	c := GetBaseConfig()
	commonPath := c.GetCommonRegexpPath()
	singlePath := c.GetSinglePath()
	r.common = &sync.Map{}
	for _, path := range commonPath {
		file, err := os.OpenFile(path, os.O_RDONLY, 0x666)
		if err != nil {
			fmt.Println("open file  error", err.Error())
			continue
		}
		defer file.Close()
		doc, err := goquery.NewDocumentFromReader(file)
		if err != nil {
			fmt.Println("creat xml error", err.Error())
			continue
		}
		doc.Find("css").Each(func(i int, selection *goquery.Selection) {
			var className []string
			val, exists := selection.Attr("key")
			if exists {
				text := strings.TrimSpace(selection.Text())
				for _, v := range strings.SplitN(text, "\n", -1) {
					if v != "" {
						className = append(className, strings.TrimSpace(v))
					}
				}
			}
			r.common.Store(regexp.MustCompile(val), className)
		})
	}
	if c.Debug() {
		// Debug
		fmt.Println("common regexp --------------------------------")
		r.common.Range(func(key, value interface{}) bool {
			fmt.Println(key, value)
			return true
		})
	}
	// single is different than vue
	switch c.GetFileType() {
	case RNCSS:
		// rn
		r.single = css.NewRnCss()
	case VueCss:
		// vue
		r.single = css.NewVueCss()
	}
	for _, path := range singlePath {
		file, err := os.OpenFile(path, os.O_RDONLY, 0x666)
		if err != nil {
			fmt.Println("open file  error", err.Error())
			continue
		}
		defer file.Close()
		doc, err := goquery.NewDocumentFromReader(file)
		if err != nil {
			fmt.Println("creat xml error", err.Error())
			continue
		}
		doc.Find("css").Each(func(i int, selection *goquery.Selection) {
			val, exists := selection.Attr("key")
			if exists {
				switch c.GetFileType() {
				case RNCSS:
					// rn
					r.single.Set(css.NewRnCssUint(regexp.MustCompile(val), selection.Text(), _cfg.GetStatic()))
				case VueCss:
					// vue
					r.single.Set(css.NewVueCssUint(regexp.MustCompile(val), selection.Text(), _cfg.GetStatic()))
				}
			}
		})
	}
	if c.Debug() {
		// Debug
		fmt.Println("single regexp --------------------------------")
		for _, v := range r.single.GetAllData().Range(0) {
			data := v.(css.Uint)
			fmt.Println(data.Reg().String(), data.Val())
		}
	}
	return nil
}

var _reg = &cssRegexp{
	common: &sync.Map{},
	single: nil,
}

func GetRegexp() Regexp {
	_reg.rwLocker.RLock()
	defer _reg.rwLocker.RUnlock()
	return _reg
}
