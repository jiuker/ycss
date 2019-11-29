package css

import (
	"regexp"
	"sync"
)

// consistent css class
// data structure should like "key":Unit
type Css interface {
	Set(classNameReg *regexp.Regexp, uint Uint)
	Get(className string) Uint
	MergeSelf(target Css)
	// get all data reg:Unit
	GetAllData() *sync.Map
}

// consistent css attrValue
type Uint interface {
	Val() interface{}
}
