/**
 * Created by shudi on 2016/10/22.
 */
import React from 'react';
import {hashHistory} from 'react-router';
import {Panel, PanelHeader, Toast, Dialog} from 'react-weui';
const {Alert} = Dialog;
import 'weui';

import {RegisterForm} from '#coms/user-form';
import PageBottom from '#coms/page-bottom';
import {User} from '#models/models';

export default class StudentRegisterPage extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            showLoading: false,
            alertTitle: '',
            alertShow: false,
            alertMsg: '',
            alertButtons: [
                {
                    label: '好的',
                    onClick: this.hideAlert.bind(this),
                }
            ],
        };
        this.onRegister = this.onRegister.bind(this);
        this.toLogin = this.toLogin.bind(this);
        this.showAlert = this.showAlert.bind(this);
    }

    onRegister(username, password) {
        this.setState({showLoading: true});
        User.studentRegister(username, password, (payload) => {
            setTimeout(() => {
                this.setState({showLoading: false});
                console.log(payload);
                // hashHistory.push('');
            }, 500);
        }, (status) => {
            this.setState({showLoading: false});
            setTimeout(() => {
                this.showAlert('注册失败', status);
            }, 500);
        });
    }

    toLogin() {
        hashHistory.push('login');
    }

    showAlert(title, msg) {
        this.setState({
            alertTitle: title,
            alertShow: true,
            alertMsg: msg,
        });
    }

    hideAlert() {
        this.setState({alertShow: false});
    }

    render() {
        return (
            <div>
                <Panel access={true}>
                    <PanelHeader style={{fontSize: "18px"}}>学生登录</PanelHeader>
                    <RegisterForm titleTip="请输入学号、密码登录（与info账号不同）"
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
                                  showAlert={this.showAlert}/>
                    <div style={{color: "#999999", padding: "10px 20px", textAlign: "center", fontSize: "13px"}}>
                        账号密码遇到任何问题，请在工作时间（周一至周五8:00-12:00；13:00-17:00）致电010-62782007。
                    </div>
                </Panel>
                <Toast icon="loading" show={this.state.showLoading}>正在加载中...</Toast>
                <Alert title={this.state.alertTitle} buttons={this.state.alertButtons} show={this.state.alertShow}>
                    {this.state.alertMsg}
                </Alert>
                <PageBottom style={{color: "#999999", textAlign: "center", backgroundColor: "white", fontSize: "14px"}}
                            contents={["清华大学学生心理发展指导中心", "联系方式：010-62782007"]}/>
            </div>
        );
    }
}