package css

import (
	"regexp"
	"strings"

	"github.com/gogf/gf/container/garray"
)

type VueCss struct {
	// []Unit
	data *garray.Array
}

func (r *VueCss) Set(unit Uint) {
	r.data.Append(unit)
}

func (r *VueCss) GetAllData() *garray.Array {
	return r.data
}

func NewVueCss() Css {
	return &VueCss{
		data: garray.NewArray(),
	}
}
func NewVueCssUint(reg *regexp.Regexp, cssVal string) Uint {
	cssVal = strings.TrimSpace(cssVal)
	return &VueCssUint{
		val: cssVal,
		reg: reg,
	}
}

type VueCssUint struct {
	val string
	reg *regexp.Regexp
}

func (v *VueCssUint) Reg() *regexp.Regexp {
	return v.reg
}

func (v *VueCssUint) Val() interface{} {
	return v.val
}
