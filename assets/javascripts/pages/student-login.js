/**
 * Created by shudi on 2016/10/22.
 */
import React from 'react';
import {hashHistory} from 'react-router';
import {Panel, PanelHeader, PanelBody} from 'react-weui';
import 'weui';

import {LoginForm} from '#coms/user-form';
import PageBottom from '#coms/page-bottom';
import {AlertDialog, LoadingHud} from '#coms/huds';
import {User} from '#models/models';

export default class StudentLoginPage extends React.Component {
    constructor(props) {
        super(props);
        this.onLogin = this.onLogin.bind(this);
        this.toRegister = this.toRegister.bind(this);
    }

    onLogin(username, password) {
        this.refs['loading'].show('正在加载中');
        User.studentLogin(username, password, () => {
            this.refs['loading'].hide();
            hashHistory.push('reservation');
        }, (status) => {
            this.refs['loading'].hide();
            setTimeout(() => {
                this.refs['alert'].show('登录失败', status, '好的');
            }, 500);
        });
    }

    toRegister() {
        hashHistory.push('register');
    }

    render() {
        return (
            <div>
                <Panel access>
                    <PanelHeader style={{fontSize: "18px"}}>学生登录</PanelHeader>
                    <PanelBody>
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
                    </PanelBody>
                </Panel>
                <LoadingHud ref="loading"/>
                <AlertDialog ref="alert"/>
                <PageBottom style={{color: "#999999", textAlign: "center", backgroundColor: "white", fontSize: "14px"}}
                            contents={["清华大学学生心理发展指导中心", "联系方式：010-62782007"]}
                            height="55px"/>
            </div>
        );
    }
}