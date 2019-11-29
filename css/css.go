package css

import (
	"regexp"

	"github.com/gogf/gf/container/garray"
)

// consistent css class
// data structure should like "key":Unit
type Css interface {
	Set(uint Uint)
	// get all data reg:Unit
	GetAllData() *garray.Array
}

// consistent css attrValue
type Uint interface {
	Reg() *regexp.Regexp
	Val() interface{}
}
