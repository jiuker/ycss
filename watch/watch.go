package watch

import (
	"fmt"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
)

type fileOp int

const (
	OpCreate     = fileOp(1)
	OpWrite      = fileOp(2)
	OpToolHandle = fileOp(3)
)

type FileWatch interface {
	Watch(pathCh chan string)
}
type fileWatch struct {
	watcher     *fsnotify.Watcher
	watcherDirs []string
	// path:*FileOp
	changeList sync.Map
	isListen   bool
}
type FileOp struct {
	op          fileOp
	ignoreCount int64
}

func NeedIgnore(filePath string) {
	fW.changeList.Store(filePath, &FileOp{
		op:          OpToolHandle,
		ignoreCount: 2,
	})
}
func (f *fileWatch) Watch(pathCh chan string) {
	go func() {
		for {
			select {
			case evt := <-f.watcher.Events:
				if strings.LastIndex(evt.Name, "~") != len(evt.Name)-1 && (evt.Op == fsnotify.Create || evt.Op == fsnotify.Write) {
					value, ok := f.changeList.Load(evt.Name)
					if !ok {
						// create
						pathCh <- evt.Name
						if evt.Op == fsnotify.Create {
							f.changeList.Store(evt.Name, &FileOp{
								op:          OpCreate,
								ignoreCount: 0,
							})
						} else {
							// write
							f.changeList.Store(evt.Name, &FileOp{
								op:          OpWrite,
								ignoreCount: 1,
							})
						}
						continue
					}
					val := value.(*FileOp)
					if evt.Op == fsnotify.Write {
						if val.ignoreCount == 0 {
							val.ignoreCount++
							pathCh <- evt.Name
						} else {
							val.ignoreCount--
						}
						f.changeList.Store(evt.Name, val)
					} else {
						// create?
						f.changeList.Store(evt.Name, &FileOp{
							op:          OpCreate,
							ignoreCount: 0,
						})
						pathCh <- evt.Name
						continue
					}
				}
			case err := <-f.watcher.Errors:
				panic(err)
			}
		}
	}()
}

var fW = &fileWatch{}

func NowWatch(dirs []string, pathCh chan string) FileWatch {
	var err error
	if fW.watcher != nil {
		// first remove watch dir
		for _, v := range fW.watcherDirs {
			fW.watcher.Remove(v)
		}
	} else {
		fW.watcher, err = fsnotify.NewWatcher()
		if err != nil {
			fmt.Println(err.Error())
			return nil
		}
	}
	if fW.isListen == false {
		fW.Watch(pathCh)
	}
	fW.watcherDirs = dirs
	for _, v := range fW.watcherDirs {
		fW.watcher.Add(v)
	}
	return fW
}
