package replace

import (
	"context"
	"errors"
	"fmt"
	"github.com/jiuker/ycss/cfg"
	"github.com/jiuker/ycss/filePath"
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
	inAndOutSame bool
	path         string
	bodyStr      string
	file         *os.File
	outFile      *os.File
	outbodyStr   string
	ctx          context.Context
	cancelFun    context.CancelFunc
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
	bodyCopy := v.GetOutFileBody()
	newWrite := strings.Replace(bodyCopy, *oldPos, *newPos, 2)
	if viper.GetBool("debug") {
		fmt.Println("will insert ", newWrite)
	}
	err := v.outFile.Truncate(0)
	if err != nil {
		return err
	}
	_, err = v.outFile.WriteAt([]byte(newWrite), 0)
	if err != nil {
		return err
	}
	err = v.outFile.Sync()
	if err != nil {
		return err
	}
	return nil
}

func (v *vueReplace) GetOldCss(reg *regexp.Regexp) (*string, *string, error) {
	if viper.GetBool("debug") {
		fmt.Println("outFileBody--------------", v.GetOutFileBody())
	}
	if !v.inAndOutSame {
		reg = regexp.MustCompile(strings.ReplaceAll(reg.String(), "Start", fmt.Sprintf(`Start\(%s\)`, v.path)))
	}
	mCssStr := reg.FindAllStringSubmatch(v.GetOutFileBody(), -1)
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
						if !strings.Contains(val, "$") {
							if viper.GetBool("debug") {
								fmt.Println(val, matchVal[0], i)
							}
							break
						}
					}
					if !strings.Contains(val, "$") {
						for _, _unit := range single.GetAllData().Range(0) {
							unit := _unit.(css.Uint)
							matchVal1 := unit.Reg().FindAllStringSubmatch(val, -1)
							val1 := unit.Val().(string)
							if len(matchVal1) == 1 {
								for i := 1; i < len(matchVal1[0]); i++ {
									// replace data
									val1 = strings.ReplaceAll(val1, fmt.Sprintf("$%v", i), matchVal1[0][i])
								}
								if !strings.Contains(val1, "$") {
									cssVal += val1
								}
								break
							}
						}
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
func (v *vueReplace) GetOutFileBody() string {
	if v.outbodyStr == "" {
		outbodyStr, err := ioutil.ReadAll(v.outFile)
		if err != nil {
			return ""
		}
		v.outbodyStr = string(outbodyStr)
	}
	return v.outbodyStr
}

func (v *vueReplace) Done() {
	v.cancelFun()
}

func (v *vueReplace) FindClass(reg []*regexp.Regexp) []string {
	return findAllCss(v.GetFileBody(), reg)
}

// new replace
func NewVueReplace(fP filePath.FilePath) (Replace, error) {
	fmt.Println("--------", fP.GetFilePath())
	path, same := fP.Format(cfg.GetBaseConfig().GetOutPath())
	if same {
		file, err := os.OpenFile(path, os.O_RDWR, 0x666)
		if err != nil {
			return nil, err
		}
		outFile, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0x666)
		if err != nil {
			file.Close()
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
			path:         fP.GetFilePath(),
			file:         file,
			outFile:      outFile,
			cancelFun:    cancelFun,
			ctx:          ctx,
			inAndOutSame: false,
		}, nil
	} else {
		file, err := os.OpenFile(fP.GetFilePath(), os.O_RDWR, 0x666)
		if err != nil {
			return nil, err
		}
		os.MkdirAll(fP.GetFileDir(), 0x666)
		outFile, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0x666)
		if err != nil {
			file.Close()
			return nil, err
		}
		ctx, cancelFun := context.WithTimeout(context.Background(), time.Duration(time.Second*5))
		go func() {
			for {
				select {
				case <-ctx.Done():
					fmt.Println("file close:", fP.GetFilePath())
					file.Close()
					fmt.Println("file close:", path)
					outFile.Close()
					return
				}
			}
		}()
		return &vueReplace{
			path:         fP.GetFilePath(),
			file:         file,
			outFile:      outFile,
			cancelFun:    cancelFun,
			ctx:          ctx,
			inAndOutSame: false,
		}, nil
	}
}
