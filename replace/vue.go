package replace

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/spf13/viper"

	"github.com/jiuker/ycss/watch"

	"github.com/jiuker/ycss/css"
)

type vueReplace struct {
	path      string
	bodyStr   string
	file      *os.File
	ctx       context.Context
	cancelFun context.CancelFunc
}

func (v *vueReplace) Replace(old *string, new *string, pos *string) *string {
	oldCopy := *old
	newCopy := *new
	posCopy := *pos
	posCopy = strings.Replace(posCopy, oldCopy, "\n"+newCopy, -1)
	return &posCopy
}
func (v *vueReplace) Save(newPos, oldPos *string) error {
	watch.NeedIgnore(v.path)
	bodyCopy := v.GetFileBody()
	newWrite := strings.Replace(bodyCopy, *oldPos, *newPos, 2)
	if viper.GetBool("debug") {
		fmt.Println("will insert ", newWrite)
	}
	_, err := v.file.WriteString(newWrite)
	if err != nil {
		return err
	}
	err = v.file.Sync()
	if err != nil {
		return err
	}
	return nil
}

func (v *vueReplace) GetOldCss(reg *regexp.Regexp) (*string, *string, error) {
	mCssStr := reg.FindAllStringSubmatch(v.GetFileBody(), -1)
	if len(mCssStr) == 0 {
		return nil, nil, errors.New("no match old css")
	}
	if len(mCssStr[0]) != 2 {
		return nil, nil, errors.New("no match old css")
	}
	return &mCssStr[0][1], &mCssStr[0][0], nil
}

func (v *vueReplace) Zoom(in *string, unit string, needZoomUnit string, needZoomKey []string, zoom float64) *string {
	reg := regexp.MustCompile(fmt.Sprintf(`[0-9|.]{1,10}[ |	]{0,3}[%s]{1,5}`, needZoomUnit))
	dataReg := regexp.MustCompile(`\d{1,5}|\d{1,5}\.\d{1,5}`)
	unitReg := regexp.MustCompile(fmt.Sprintf(`%s`, needZoomUnit))
	_in := reg.ReplaceAllStringFunc(*in, func(s string) string {
		s = dataReg.ReplaceAllStringFunc(s, func(s1 string) string {
			sFloat, err := strconv.ParseFloat(s1, 64)
			if err != nil {
				return s1
			}
			return fmt.Sprintf("%v", sFloat*zoom)
		})
		s = unitReg.ReplaceAllStringFunc(s, func(s2 string) string {
			if s2 == "" {
				return ""
			}
			return unit
		})
		return s
	})
	in = &_in
	return in
}

func (v *vueReplace) GetRegexpCss(cls []string, common *sync.Map, single css.Css) *string {
	var genClass = ``
	for _, cl := range cls {
		common.Range(func(key, value interface{}) bool {
			matchVal := key.(*regexp.Regexp).FindAllStringSubmatch(cl, -1)
			if len(matchVal) == 1 {
				var cssVal = ""
				for _, val := range value.([]string) {
					for i := 1; i < len(matchVal[0]); i++ {
						// replace data
						val = strings.ReplaceAll(val, fmt.Sprintf("$%v", i), matchVal[0][i])
					}
					if !strings.Contains(val, "$") {
						single.GetAllData().Range(func(key1, value1 interface{}) bool {
							matchVal1 := key1.(*regexp.Regexp).FindAllStringSubmatch(val, -1)
							val1 := value1.(css.Uint).Val().(string)
							if len(matchVal1) == 1 {
								for i := 1; i < len(matchVal1[0]); i++ {
									// replace data
									val1 = strings.ReplaceAll(val1, fmt.Sprintf("$%v", i), matchVal1[0][i])
								}
								if !strings.Contains(val1, "$") {
									cssVal += val1
								}
								return false
							}
							return true
						})
					}
				}
				if cssVal != "" {
					genClass += fmt.Sprintf(".%s{%s}\n", cl, cssVal)
				}
				return false
			}
			return true
		})
	}
	return &genClass
}

func (v *vueReplace) GetFileBody() string {
	if v.bodyStr == "" {
		bodyByt, err := ioutil.ReadAll(v.file)
		if err != nil {
			return ""
		}
		v.bodyStr = string(bodyByt)
	}
	return v.bodyStr
}

func (v *vueReplace) Done() {
	v.cancelFun()
}

func (v *vueReplace) FindClass(reg []*regexp.Regexp) []string {
	return findAllCss(v.GetFileBody(), reg)
}

// new replace
func NewVueReplace(path string) (Replace, error) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_SYNC, 0x666)
	if err != nil {
		return nil, err
	}
	ctx, cancelFun := context.WithTimeout(context.Background(), time.Duration(time.Second*5))
	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("file close:", path)
				file.Close()
				return
			}
		}
	}()
	return &vueReplace{
		path:      path,
		file:      file,
		cancelFun: cancelFun,
		ctx:       ctx,
	}, nil
}
