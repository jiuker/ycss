package replace

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/jiuker/ycss/watch"

	"github.com/spf13/viper"

	"github.com/jiuker/ycss/css"
)

type rnReplace struct {
	path      string
	bodyStr   string
	file      *os.File
	ctx       context.Context
	cancelFun context.CancelFunc
}

func (v *rnReplace) GetFileBody() string {
	if v.bodyStr == "" {
		bodyByt, err := ioutil.ReadAll(v.file)
		if err != nil {
			return ""
		}
		v.bodyStr = string(bodyByt)
	}
	return v.bodyStr
}

func (v *rnReplace) FindClass(reg []*regexp.Regexp) []string {
	return findAllCss(v.GetFileBody(), reg)
}

func (v *rnReplace) GetRegexpCss(cls []string, common *sync.Map, single css.Css) *string {
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
						for _, _unit := range single.GetAllData().Range(0) {
							unit := _unit.(css.Uint)
							matchVal1 := unit.Reg().FindAllStringSubmatch(val, -1)
							val1 := unit.Val().(string)
							if len(matchVal1) == 1 {
								for i := 1; i < len(matchVal1[0]); i++ {
									// replace data
									val1 = strings.ReplaceAll(val1, fmt.Sprintf("-%v", i), matchVal1[0][i])
								}
								if !strings.Contains(val1, "-") {
									cssVal += strings.TrimSpace(val1) + ","
								}
								break
							}
						}
					}
				}
				if cssVal != "" {
					if string(cssVal[len(cssVal)-1]) == "," {
						cssVal = cssVal[0 : len(cssVal)-1]
					}
					genClass += fmt.Sprintf(`"%s":{%s},`, cl, cssVal)
				}
				return false
			}
			return true
		})
	}
	if genClass != "" {
		if string(genClass[len(genClass)-1]) == "," {
			genClass = genClass[0 : len(genClass)-1]
		}
		genClass = fmt.Sprintf(`{%s}`, genClass)
	}
	return &genClass
}

func (v *rnReplace) Zoom(css *string, unit string, needZoomUint string, keyNeedZoom []string, zoom float64) *string {
	// unit,needZoomUint is not effect
	var data interface{}
	err := json.Unmarshal([]byte(*css), &data)
	if err != nil {
		return css
	}
	dataByte, err := json.MarshalIndent(walkToSet(data, "", keyNeedZoom, zoom), "", "	")
	if err != nil {
		return css
	}
	str := string(dataByte)
	if viper.GetBool("debug") {
		fmt.Println("zoom ing...", str)
	}
	// rcss dont need outline{}
	str = str[1:len(str)-1] + ","
	return &str
}
func walkToSet(data interface{}, key string, keyNeedZoom []string, zoom float64) interface{} {
	switch data.(type) {
	case float64:
		var needZoom = false
		for _, v := range keyNeedZoom {
			if v == key {
				needZoom = true
				break
			}
		}
		if needZoom {
			data = data.(float64) * zoom
		}
	case []interface{}:
		for k, _ := range data.([]interface{}) {
			data.([]interface{})[k] = walkToSet(data.([]interface{})[k], "", keyNeedZoom, zoom)
		}
	case map[string]interface{}:
		for k, _ := range data.(map[string]interface{}) {
			v := data.(map[string]interface{})[k]
			data.(map[string]interface{})[k] = walkToSet(v, k, keyNeedZoom, zoom)
		}
	}
	return data
}
func (v *rnReplace) GetOldCss(reg *regexp.Regexp) (*string, *string, error) {
	mCssStr := reg.FindAllStringSubmatch(v.GetFileBody(), -1)
	if len(mCssStr) == 0 {
		return nil, nil, errors.New("no match old css")
	}
	if len(mCssStr[0]) != 2 {
		return nil, nil, errors.New("no match old css")
	}
	return &mCssStr[0][1], &mCssStr[0][0], nil
}

func (v *rnReplace) Replace(old *string, new *string, pos *string) *string {
	oldCopy := *old
	newCopy := *new
	posCopy := *pos
	posCopy = strings.Replace(posCopy, oldCopy, "\n"+newCopy+"\n", -1)
	return &posCopy
}

func (v *rnReplace) Save(newPos *string, oldPos *string) error {
	watch.NeedIgnore(v.path)
	bodyCopy := v.GetFileBody()
	newWrite := strings.Replace(bodyCopy, *oldPos, *newPos, 2)
	if viper.GetBool("debug") {
		fmt.Println("will insert ", newWrite)
	}
	err := v.file.Truncate(0)
	if err != nil {
		return err
	}
	_, err = v.file.WriteAt([]byte(newWrite), 0)
	if err != nil {
		return err
	}
	err = v.file.Sync()
	if err != nil {
		return err
	}
	return nil
}

func (v *rnReplace) Done() {
	v.cancelFun()
}

// new replace
func NewRnReplace(path string) (Replace, error) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_SYNC, 0x666)
	if err != nil {
		return nil, err
	}
	ctx, cancelFun := context.WithTimeout(context.Background(), time.Duration(time.Second*5))
	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("file close!")
				file.Close()
				return
			}
		}
	}()
	return &rnReplace{
		path:      path,
		file:      file,
		cancelFun: cancelFun,
		ctx:       ctx,
	}, nil
}
