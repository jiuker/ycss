package handle

import (
	"fmt"
	"time"

	"github.com/jiuker/ycss/filePath"
	"github.com/jiuker/ycss/log"
	"github.com/jiuker/ycss/replace"

	"github.com/jiuker/ycss/cfg"
)

func StartHandle() {
	go func() {
		for {
			select {
			case path := <-cfg.ChangeFilePath:
				go func() {
					// too fast will error?
					time.Sleep(time.Millisecond * 300)
					log.Log("the url will handle", path)
					var pla replace.Replace
					var err error
					f := filePath.NewFilePath(path)
					switch cfg.GetBaseConfig().GetFileType() {
					case cfg.VueCss:
						if f.GetFileType() != ".vue" && f.GetFileType() != ".nvue" && f.GetFileType() != ".html" && f.GetFileType() != ".htm" {
							return
						}
						pla, err = replace.NewVueReplace(f)
						if err != nil {
							fmt.Println(err.Error())
							return
						}
						defer pla.Done()
					case cfg.RNCSS:
						if f.GetFileType() != ".js" && f.GetFileType() != ".jsx" {
							return
						}
						pla, err = replace.NewRnReplace(f)
						if err != nil {
							fmt.Println(err.Error())
							return
						}
						defer pla.Done()
					}
					cls := pla.FindClass(cfg.GetBaseConfig().GetReg())
					log.Log("this handle print class is ", cls)
					// cls is nil,may error
					if len(cls) == 0 {
						return
					}
					rcss := pla.GetRegexpCss(cls, cfg.GetRegexp().GetCommonReg(), cfg.GetRegexp().GetSingle())
					log.Log("zoom befer:", *rcss)
					rcss = pla.Zoom(rcss, cfg.GetBaseConfig().GetOutUnit(), cfg.GetBaseConfig().GetNeedZoomUnit(), cfg.GetBaseConfig().GetKeyNeedZoom(), cfg.GetBaseConfig().GetZoom())
					log.Log("zoom after:", *rcss)
					old, pos, err := pla.GetOldCss(cfg.GetBaseConfig().GetOldCssReg())
					if err != nil {
						fmt.Println(err.Error())
						return
					}
					log.Log("old:", *old)
					if replace.IsSame(rcss, old) {
						log.Log("just is same,do nothing!")
						return
					}
					newPos := pla.Replace(old, rcss, pos)
					log.Log("new pos:", *newPos)
					fmt.Println(pla.Save(newPos, pos))
				}()
			}
		}
	}()
}
