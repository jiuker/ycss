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
            <div styles={GetStyle("w-110 h-25 tX-12")}>123</div>
        )
    }
}
function GetStyle(className) {
   return styles[md5]
}

const styles = StyleSheet.create({
    /* Automatic generation Start */
    "container": {
        flex: 1,
        justifyContent: 'center',
        alignItems: 'center',
        backgroundColor: '#F5FCFF',
    },
    /* Automatic generation End */
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
});