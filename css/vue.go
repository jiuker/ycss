package css

import (
	"regexp"
	"strings"
	"sync"
)

type VueCss struct {
	// map[string]Unit
	data *sync.Map
}

func (r *VueCss) Set(classNameReg *regexp.Regexp, uint Uint) {
	r.data.Store(classNameReg, uint)
}

func (r *VueCss) Get(className string) Uint {
	data, ok := r.data.Load(className)
	if !ok {
		return nil
	}
	return data.(Uint)
}

func (r *VueCss) MergeSelf(target Css) {
	target.GetAllData().Range(func(key, value interface{}) bool {
		r.data.Store(key, value)
		return true
	})
}
func (r *VueCss) GetAllData() *sync.Map {
	return r.data
}

func NewVueCss() Css {
	return &VueCss{
		data: &sync.Map{},
	}
}
func NewVueCssUint(cssVal string) Uint {
	cssVal = strings.TrimSpace(cssVal)
	return &VueCssUint{
		val: cssVal,
	}
}

type VueCssUint struct {
	val string
}

func (v *VueCssUint) Val() interface{} {
	return v.val
}
