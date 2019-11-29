# ycss
Only one configuration is needed, you can automatically complete your own style of CSS code!
# Required
* go version>=1.13
# How To Run
* go build -o ycss main.go
* ./ycss
# Config
## Features
* write config will work immediately
## Field
    {
         "debug": true, // debug mod
         "type": "rn", // rn|vue
         "common": ["./res/regexp/common/rn.reg"],// reg dir
         "single": ["./res/regexp/single/rn.reg"],// single dir
         "outUnit": "upx", // out unit,will work in vue
         "zoom": 1.4, // value zoom,if px need to rem;
         "needZoomUnit": "px|rem", // vaule unit will zoom,if not match will do nothing 
         "reg": ["GetStyle\\(\"([^\"]+)\""],// how to find your class,eq:class="w-15"
         "watchDir": ["./res/sample/rn"],// watch dir to do
         "oldCssReg": "/\\* Automatic generation Start \\*/([^/]+)/\\*", // both vue and rn do as <xxxx start> (auto code) <xxx end>,if not match do nothing;
         "keyNeedZoom": [ // will work in rn,if value need zoom ,please set it
             "width",
             "height"
         ]
     }
