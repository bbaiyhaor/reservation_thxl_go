/**
 * Created by shudi on 2016/10/22.
 */
import React from 'react';
import {hashHistory} from 'react-router';
import {Panel, PanelHeader, Toast, Dialog} from 'react-weui';
const {Alert} = Dialog;
import 'weui';

import {LoginForm} from '#coms/user-form';
import PageBottom from '#coms/page-bottom';
import {User} from '#models/models';

export default class StudentLoginPage extends React.Component {
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
        this.onLogin = this.onLogin.bind(this);
        this.toRegister = this.toRegister.bind(this);
        this.showAlert = this.showAlert.bind(this);
    }

    onLogin(username, password) {
        this.setState({showLoading: true});
        User.studentLogin(username, password, (payload) => {
            setTimeout(() => {
                this.setState({showLoading: false});
                console.log(payload);
                // hashHistory.push('');
            }, 500);
        }, (status) => {
            this.setState({showLoading: false});
            setTimeout(() => {
                this.showAlert('登录失败', status);
            }, 500);
        });
    }

    toRegister() {
        hashHistory.push('register');
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
                    <LoginForm titleTip="请输入学号、密码登录（与info账号不同）"
                               usernameLabel="学号"
                               usernameType="tel"
                               usernamePlaceholder="请输入学号"
                               passwordLabel="密码"
                               passwordType="password"
                               passwordPlaceholder="请输入密码"
                               submitText="登录"
                               cancelText="没有账户"
                               onSubmit={this.onLogin}
                               onCancel={this.toRegister}/>
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