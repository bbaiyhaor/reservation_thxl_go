/**
 * Created by shudi on 2016/10/22.
 */
import React from 'react';
import {hashHistory} from 'react-router';
import {Panel, PanelHeader} from 'react-weui';
import 'weui';

import {RegisterForm} from '#coms/user-form';
import PageBottom from '#coms/page-bottom';

export default class StudentRegisterPage extends React.Component {
    constructor(props) {
        super(props);
        this.onRegister = this.onRegister.bind(this);
        this.toLogin = this.toLogin.bind(this);
    }

    onRegister() {

    }

    toLogin() {
        hashHistory.push('login');
    }

    render() {
        return (
            <div>
                <Panel access={true}>
                    <PanelHeader style={{fontSize: "18px"}}>学生注册</PanelHeader>
                    <RegisterForm titleTip="请用学号和密码注册（密码与info账号不同）"
                               names={["学号", "密码", "确认密码"]}
                               types={["tel", "number", "number"]}
                               placeholders={["请输入学号", "请输入密码", "请确认密码"]}
                               protocolNames={["咨询协议"]}
                               protocolLinks={["protocol"]}
                               protocolPrefix={["我已阅读并同意"]}
                               protocolSuffix={[""]}
                               submitText="注册"
                               cancelText="已有账户"
                               onSubmit={this.onSubmit}
                               onCancel={this.onCancel}/>
                    <div style={{color: "#999999", padding: "10px 20px", textAlign: "center", fontSize: "13px"}}>
                        账号密码遇到任何问题，请在工作时间（周一至周五8:00-12:00；13:00-17:00）致电010-62782007。
                    </div>
                </Panel>
                <PageBottom style={{color: "#999999", textAlign: "center", backgroundColor: "white", fontSize: "14px"}}
                            contents={["清华大学学生心理发展指导中心", "联系方式：010-62782007"]}/>
            </div>
        );
    }
}