import React from 'react';

export default class App extends React.Component {
    constructor(props) {
        super(props);
        // 初始化state
    }

    componentDidMount() {
        // 数据请求
    }

    componentWillReceiveProps() {
        // 在组件接收到一个新的 prop (更新后)时被调用。这个方法在初始化render时不会被调用
    }

    componentWillUnmount() {
        // 销毁长链接等本组件占用资源的操作
    }

    render() {
        // 注意render一个组件
        return (
            <div styles={GetStyle("w-20-h-18 tX-16 c-r")}>123</div>
        )
    }
}
function GetStyle(className) {
   return styles[md5]
}

const styles = StyleSheet.create({
    welcome: {
        fontSize: 20,
        textAlign: 'center',
        margin: 10,
    },
    instructions: {
        textAlign: 'center',
        color: '#333333',
        marginBottom: 5,
    },
    /* Automatic generation Start */

	"c-r": {
		"color": "#ff0000"
	},
	"tX-16": {
		"transform": [
			{
				"translateX": 16
			},
			{
				"translateY": 16
			}
		],
		"width": 22.4
	},
	"w-20-h-30": {
		"height": 42,
		"width": 28
	}

/* Automatic generation End */
});