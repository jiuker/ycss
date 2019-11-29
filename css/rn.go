package css

import (
	"regexp"
	"sync"
)

type RnCss struct {
	// map[string]Unit
	data *sync.Map
}

func (r *RnCss) Set(classNameReg *regexp.Regexp, uint Uint) {
	r.data.Store(classNameReg, uint)
}

func (r *RnCss) Get(className string) Uint {
	data, ok := r.data.Load(className)
	if !ok {
		return nil
	}
	return data.(Uint)
}

func (r *RnCss) MergeSelf(target Css) {
	target.GetAllData().Range(func(key, value interface{}) bool {
		r.data.Store(key, value)
		return true
	})
}
func (r *RnCss) GetAllData() *sync.Map {
	return r.data
}

func NewRnCss() Css {
	return &RnCss{
		data: &sync.Map{},
	}
}
func NewRnCssUint(cssVal string) Uint {
	return &RnCssUint{
		val: cssVal,
	}
}

type RnCssUint struct {
	val interface{}
}

func (v *RnCssUint) Val() interface{} {
	return v.val
}
