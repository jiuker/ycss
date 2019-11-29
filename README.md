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
# Reg
## Common
Common is an intermediary mechanism, which can be understood as a container
means:
  w-15-h-20{w-15 h-20 pl-10}
* key is the regexp
* value $1,$2 mean regexp match value,w-($1)-h-($2)
## Single
Single is the most basic style expression
means:
   h-20{height:20px}
### Vue
* key is the regexp
* value $1,$2 mean regexp match value,w-($1)-h-($2)
### RN
Why do like -1,-2?
-1,-2 is a special value and can keep the original data type
* key is the regexp
* value -1,-2 mean regexp match value,w-(-1)-h-(-2),-1,-2,-3,-4,-5,-6 also can work

# QQ
* 941057162
