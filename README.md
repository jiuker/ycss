# ycss
Only one configuration is needed, you can automatically complete your own style of CSS code!
You just need to write class, CSS is generated by us!

![Image text](https://github.com/jiuker/ycss/blob/master/res/vue.gif)

![Image text](https://github.com/jiuker/ycss/blob/master/res/react-native.gif)
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
         ],
         "static": { // work at @key,code is newCssUint
              "red": "#ff0000",
              "red-1": "#ff1100",
              "height": 12
         },
         "outPath": "@FILEDIR/@FILENAME.@FILETYPE" // if static ,will write in a file,dir/name.type is constant
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
# example
CSS is automatically generated and can be configured! Can achieve a frame effect, this is your own frame!
The default rule is demo, and you can write your own structure.
## /res/sample
    <template>
        <div class="bc-ff1123"></div>
        <div class="bc-000-112-231 br-nr bp-c bs-c bs-10-15"></div>
        <div class="b-1-001 br-1-123 o-1-000121 c-fff ls-12 lh-20"></div>
        <div class="ta-c ta-r ta-l"></div>
        <div class="fs-20 fw-100"></div>
        <div class="m-1010 p-0505 h-10 w-20 h10 w10"></div>
        <div class="maxh-23 maxw-10 minh-10 minw-22"></div>
        <div class="p-f p-a p-r d-b t-2 b-1 l-3 r-40 va-m zi-205"></div>
        <div class="mt-10 ml-10 mr-10 mb-10"></div>
        <div class="pt-10 pl-10 pr-10 pb-10 br-1"></div>
        <div class="d-f fd-r ai-c jc-c ai-c fw-nw f-21 test1"></div>
    </template>
    <style>
        .test{
            width: 10px;
        }
        /* Automatic generation Start */
    .bc-ff1123{background-color:#ff1123;}
    .bc-000-112-231{background-color:rgb(000,112,231);}
    .br-nr{background-repeat:no-repeat;}
    .bp-c{background-position:center;}
    .bs-c{background-size:cover;}
    .bs-10-15{background-size:20px 30px;}
    .b-1-001{border:2px solid #001;}
    .br-1-123{border-right:2px solid #123;}
    .o-1-000121{outline:#000121 dotted 2px;}
    .c-fff{color:#fff;}
    .ls-12{letter-spacing:24px;}
    .lh-20{line-height:40px;}
    .ta-c{text-align:center;}
    .ta-r{text-align:right;}
    .ta-l{text-align:left;}
    .fs-20{font-size:40px;}
    .fw-100{font-weight:100;}
    .m-1010{margin:20px 20px;}
    .p-0505{padding:10px 10px;}
    .h-10{height:20px;}
    .w-20{width:40px;}
    .h10{height:10%;}
    .w10{width:10%;}
    .maxh-23{max-height:46px;}
    .maxw-10{max-width:20px;}
    .minh-10{min-height:20px;}
    .minw-22{min-width:44px;}
    .p-f{position:fixed;}
    .p-a{position:absolute;}
    .p-r{position:relative;}
    .d-b{display:block;}
    .t-2{top:4px;}
    .b-1{bottom:2px;}
    .l-3{left:6px;}
    .r-40{right:80px;}
    .va-m{vertical-align:middle;}
    .zi-205{z-index:205;}
    .mt-10{margin-top:20px;}
    .ml-10{margin-left:20px;}
    .mr-10{margin-right:20px;}
    .mb-10{margin-bottom:20px;}
    .pt-10{padding-top:20px;}
    .pl-10{padding-left:20px;}
    .pr-10{padding-right:20px;}
    .pb-10{padding-bottom:20px;}
    .br-1{border-radius:2px;}
    .d-f{display: -webkit-flex;
        display: flex;}
    .fd-r{flex-direction:row;}
    .ai-c{align-items:center;}
    .jc-c{justify-content:center;}
    .fw-nw{flex-wrap:nowrap;}
    .f-11{flex:11;}
    /* Automatic generation End */
    </style>
# QQGroup
* 941057162