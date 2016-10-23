/**
 * Created by shudi on 2016/10/22.
 */
import React from 'react';
import {hashHistory} from 'react-router';
import {Panel, PanelHeader, Toast} from 'react-weui';
import 'weui';

import {RegisterForm} from '#coms/user-form';
import PageBottom from '#coms/page-bottom';
import {AlertDialog, LoadingHud} from '#coms/huds';
import {User} from '#models/models';

export default class StudentRegisterPage extends React.Component {
    constructor(props) {
        super(props);
        this.onRegister = this.onRegister.bind(this);
        this.toLogin = this.toLogin.bind(this);
        this.alert = this.alert.bind(this);
    }

    onRegister(username, password) {
        this.refs['loading'].show('正在加载中');
        User.studentRegister(username, password, () => {
            this.refs['loading'].hide();
            hashHistory.push('reservation');
        }, (status) => {
            this.refs['loading'].hide();
            setTimeout(() => {
                this.refs['alert'].show('注册失败', status, '好的');
            }, 500);
        });
    }

    toLogin() {
        hashHistory.push('login');
    }

    alert(title, msg, label) {
        this.refs['alert'].show(title, msg, label);
    }

    render() {
        return (
            <div>
                <Panel access={true}>
                    <PanelHeader style={{fontSize: "18px"}}>学生注册</PanelHeader>
                    <RegisterForm titleTip="请用学号注册（密码与info账号不同）"
                                  usernameLabel="学号"
                                  usernameType="tel"usernamePlaceholder="请输入学号"
                                  passwordLabel="密码"
                                  passwordType="password"
                                  passwordPlaceholder="请输入密码"
                                  confirmPasswordLabel="确认密码"
                                  confirmPasswordType="password"
                                  confirmPasswordPlaceholder="请确认密码"
                                  protocol="咨询协议"
                                  protocolPrefix="我已阅读并同意"
                                  protocolLink="protocol"
                                  submitText="注册"
                                  cancelText="已有账户"
                                  onSubmit={this.onRegister}
                                  onCancel={this.toLogin}
                                  showAlert={this.alert}/>
                    <div style={{color: "#999999", padding: "10px 20px", textAlign: "center", fontSize: "13px"}}>
                        账号密码遇到任何问题，请在工作时间（周一至周五8:00-12:00；13:00-17:00）致电010-62782007。
                    </div>
                </Panel>
                <LoadingHud ref="loading"/>
                <AlertDialog ref="alert"/>
                <PageBottom style={{color: "#999999", textAlign: "center", backgroundColor: "white", fontSize: "14px"}}
                            contents={["清华大学学生心理发展指导中心", "联系方式：010-62782007"]}
                            height="50px"/>
            </div>
        );
    }
}