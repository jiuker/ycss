package css

import (
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/gogf/gf/container/garray"
)

type RnCss struct {
	// []unit
	data *garray.Array
}

func (r *RnCss) Set(unit Uint) {
	r.data.Append(unit)
}

func (r *RnCss) GetAllData() *garray.Array {
	return r.data
}

func NewRnCss() Css {
	return &RnCss{
		data: garray.NewArray(),
	}
}
func NewRnCssUint(reg *regexp.Regexp, cssVal string, staticMap *sync.Map) Uint {
	cssVal = strings.TrimSpace(cssVal)
	// replace static value @key
	staticMap.Range(func(key, value interface{}) bool {
		cssVal = strings.Replace(cssVal, "@"+key.(string), fmt.Sprintf("%v", value), -1)
		if strings.Contains(cssVal, "@") {
			return true
		}
		return false
	})
	return &RnCssUint{
		val: cssVal,
		reg: reg,
	}
}

type RnCssUint struct {
	val interface{}
	reg *regexp.Regexp
}

func (v *RnCssUint) Reg() *regexp.Regexp {
	return v.reg
}

func (v *RnCssUint) Val() interface{} {
	return v.val
}
