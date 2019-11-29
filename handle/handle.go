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
				switch cfg.GetBaseConfig().GetFileType() {
				case cfg.VueCss:
					func() {
						pla, err := replace.NewVueReplace(path)
						if err != nil {
							fmt.Println(err.Error())
							return
						}
						defer pla.Done()
						cls := pla.FindClass(cfg.GetBaseConfig().GetReg())
						if viper.GetBool("debug") {
							fmt.Println("this handle print class is ", cls)
						}
						rcss := pla.GetRegexpCss(cls, cfg.GetRegexp().GetCommonReg(), cfg.GetRegexp().GetSingle())
						rcss = pla.Zoom(rcss, cfg.GetBaseConfig().GetOutUnit(), cfg.GetBaseConfig().GetNeedZoomUnit(), cfg.GetBaseConfig().GetKeyNeedZoom(), cfg.GetBaseConfig().GetZoom())
						if viper.GetBool("debug") {
							fmt.Println("after:", *rcss)
						}
						line, ok := cfg.GetBaseConfig().GetOldCssCommonSplit()
						old, pos, err := pla.GetOldCss(cfg.GetBaseConfig().GetOldCssReg(), cfg.GetBaseConfig().GetOldCssIndex(), ok, line)
						if err != nil {
							fmt.Println(err.Error())
							return
						}
						if viper.GetBool("debug") {
							fmt.Println("old:", *old)
						}
						if replace.IsSame(rcss, old) {
							// the same,do nothing
							return
						}
						newPos := pla.Replace(old, rcss, pos)
						if viper.GetBool("debug") {
							fmt.Println("new pos:", *newPos)
						}
						fmt.Println(pla.Save(newPos, pos))
					}()
				case cfg.RNCSS:
					func() {
						pla, err := replace.NewRnReplace(path)
						if err != nil {
							fmt.Println(err.Error())
							return
						}
						defer pla.Done()
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
						line, ok := cfg.GetBaseConfig().GetOldCssCommonSplit()
						old, pos, err := pla.GetOldCss(cfg.GetBaseConfig().GetOldCssReg(), cfg.GetBaseConfig().GetOldCssIndex(), ok, line)
						if err != nil {
							fmt.Println(err.Error())
							return
						}
						if viper.GetBool("debug") {
							fmt.Println("old:", *old)
						}
						if replace.IsSame(rcss, old) {
							// the same,do nothing
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
		}
	}()
}
