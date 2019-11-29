package css

import (
	"regexp"

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
func NewRnCssUint(reg *regexp.Regexp, cssVal string) Uint {
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
