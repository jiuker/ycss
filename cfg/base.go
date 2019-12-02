package cfg

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/jiuker/ycss/watch"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// base cfg
// like dir or other
type Config interface {
	Debug() bool
	GetOutUnit() string
	GetCommonRegexpPath() []string
	GetSinglePath() []string
	GetFileType() FileType
	GetZoom() float64
	GetNeedZoomUnit() string
	GetReg() []*regexp.Regexp
	GetOldCssReg() *regexp.Regexp
	GetWatchDir() []string
	GetKeyNeedZoom() []string
	GetStatic() *sync.Map
	GetOutPath() string
}
type cfg struct {
	// should like "px","rem"
	OutUnit      string
	CommonPath   []string
	SinglePath   []string
	CssType      FileType
	rwLocker     sync.RWMutex
	notifyCount  int64
	Zoom         float64
	NeedZoomUnit string
	Reg          []*regexp.Regexp
	WatchDir     []string
	IsDebug      bool
	OldCssReg    *regexp.Regexp
	KeyNeedZoom  []string
	StaticMap    *sync.Map
	OutPath      string
}

func (c *cfg) GetOutPath() string {
	return c.OutPath
}

func (c *cfg) GetStatic() *sync.Map {
	c.rwLocker.RLock()
	defer c.rwLocker.RUnlock()
	return c.StaticMap
}

func (c *cfg) GetKeyNeedZoom() []string {
	c.rwLocker.RLock()
	defer c.rwLocker.RUnlock()
	return c.KeyNeedZoom
}

func (c *cfg) GetOldCssReg() *regexp.Regexp {
	c.rwLocker.RLock()
	defer c.rwLocker.RUnlock()
	return c.OldCssReg
}

func (c *cfg) GetNeedZoomUnit() string {
	c.rwLocker.RLock()
	defer c.rwLocker.RUnlock()
	return c.NeedZoomUnit
}

func (c *cfg) Debug() bool {
	c.rwLocker.RLock()
	defer c.rwLocker.RUnlock()
	return c.IsDebug
}
func (c *cfg) GetWatchDir() []string {
	c.rwLocker.RLock()
	defer c.rwLocker.RUnlock()
	return c.WatchDir
}

func (c *cfg) GetReg() []*regexp.Regexp {
	c.rwLocker.RLock()
	defer c.rwLocker.RUnlock()
	return c.Reg
}

func (c *cfg) GetZoom() float64 {
	c.rwLocker.RLock()
	defer c.rwLocker.RUnlock()
	return c.Zoom
}

func (c *cfg) GetCommonRegexpPath() []string {
	c.rwLocker.RLock()
	defer c.rwLocker.RUnlock()
	return c.CommonPath
}

func (c *cfg) GetSinglePath() []string {
	c.rwLocker.RLock()
	defer c.rwLocker.RUnlock()
	return c.SinglePath
}

func (c *cfg) GetOutUnit() string {
	c.rwLocker.RLock()
	defer c.rwLocker.RUnlock()
	return c.OutUnit
}
func (c *cfg) GetFileType() FileType {
	c.rwLocker.RLock()
	defer c.rwLocker.RUnlock()
	return c.CssType
}

var _cfg = &cfg{}

func GetBaseConfig() Config {
	_cfg.rwLocker.RLock()
	defer _cfg.rwLocker.RUnlock()
	return _cfg
}

// set base config file path
func SetBasePath(path string) {
	viper.SetConfigType("json")
	viper.AddConfigPath(path)
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		refreshCfg()
	})
	// start read config
	viper.ReadInConfig()
	refreshCfg()
}

var ChangeFilePath = make(chan string, 50)

func refreshCfg() {
	atomic.AddInt64(&_cfg.notifyCount, 1)
	if atomic.LoadInt64(&_cfg.notifyCount)%2 == 0 {
		return
	}
	_cfg.rwLocker.Lock()
	defer func() {
		// reload data
		_reg.Parse()
		// reload watch
		watch.NowWatch(GetBaseConfig().GetWatchDir(), ChangeFilePath)
	}()
	defer _cfg.rwLocker.Unlock()
	typ := viper.GetString("type")
	if strings.ToLower(typ) == "vue" {
		_cfg.CssType = VueCss
	} else if strings.ToLower(typ) == "rn" {
		_cfg.CssType = RNCSS
	} else {
		fmt.Println("set type error!must be vue or rn!")
		return
	}
	_cfg.CommonPath = viper.GetStringSlice("common")
	_cfg.SinglePath = viper.GetStringSlice("single")
	_cfg.OutUnit = viper.GetString("outUnit")
	_cfg.Zoom = viper.GetFloat64("zoom")
	_cfg.NeedZoomUnit = viper.GetString("needZoomUnit")
	_cfg.WatchDir = viper.GetStringSlice("watchDir")
	_cfg.IsDebug = viper.GetBool("debug")
	_cfg.OldCssReg = regexp.MustCompile(fmt.Sprintf(`%s`, viper.GetString("oldCssReg")))
	_cfg.KeyNeedZoom = viper.GetStringSlice("keyNeedZoom")
	_cfg.StaticMap = &sync.Map{}
	for k, v := range viper.GetStringMap("static") {
		_cfg.StaticMap.Store(k, v)
	}
	_cfg.OutPath = viper.GetString("outPath")
	// clear it regexp
	_cfg.Reg = []*regexp.Regexp{}
	for _, v := range viper.GetStringSlice("reg") {
		_cfg.Reg = append(_cfg.Reg, regexp.MustCompile(v))
	}
	outData, err := json.Marshal(_cfg)
	if err == nil {
		fmt.Println(string(outData))
		fmt.Println("reg is:", _cfg.Reg)
	}
}
