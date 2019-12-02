package handle

import (
	"fmt"

	"github.com/jiuker/ycss/replace"

	"github.com/spf13/viper"

	"github.com/jiuker/ycss/cfg"
)

func StartHandle() {
	go func() {
		for {
			select {
			case path := <-cfg.ChangeFilePath:
				if viper.GetBool("debug") {
					fmt.Println("the url will handle", path)
				}
				var pla replace.Replace
				var err error
				func() {
					switch cfg.GetBaseConfig().GetFileType() {
					case cfg.VueCss:
						pla, err = replace.NewVueReplace(path)
						if err != nil {
							fmt.Println(err.Error())
							return
						}
						defer pla.Done()
					case cfg.RNCSS:
						pla, err = replace.NewRnReplace(path)
						if err != nil {
							fmt.Println(err.Error())
							return
						}
						defer pla.Done()
					}
					cls := pla.FindClass(cfg.GetBaseConfig().GetReg())
					if viper.GetBool("debug") {
						fmt.Println("this handle print class is ", cls)
					}
					rcss := pla.GetRegexpCss(cls, cfg.GetRegexp().GetCommonReg(), cfg.GetRegexp().GetSingle())
					if viper.GetBool("debug") {
						fmt.Println("zoom befer:", *rcss)
					}
					rcss = pla.Zoom(rcss, cfg.GetBaseConfig().GetOutUnit(), cfg.GetBaseConfig().GetNeedZoomUnit(), cfg.GetBaseConfig().GetKeyNeedZoom(), cfg.GetBaseConfig().GetZoom())
					if viper.GetBool("debug") {
						fmt.Println("zoom after:", *rcss)
					}
					old, pos, err := pla.GetOldCss(cfg.GetBaseConfig().GetOldCssReg())
					if err != nil {
						fmt.Println(err.Error())
						return
					}
					if viper.GetBool("debug") {
						fmt.Println("old:", *old)
					}
					if replace.IsSame(rcss, old) {
						if viper.GetBool("debug") {
							fmt.Println("just is same,do nothing!")
						}
						return
					}
					newPos := pla.Replace(old, rcss, pos)
					if viper.GetBool("debug") {
						fmt.Println("new pos:", *newPos)
					}
					fmt.Println(pla.Save(newPos, pos))
				}()
			}
		}
	}()
}
