/**
 * Created by shudi on 2016/10/22.
 */
import React from 'react';
import { hashHistory } from 'react-router';
import { Panel, PanelHeader, PanelBody } from '#react-weui';
import 'weui';

import LoginForm from '#forms/LoginForm';
import PageBottom from '#coms/PageBottom';
import { AlertDialog, LoadingHud } from '#coms/Huds';
import { User } from '#models/Models';

export default class TeacherLoginPage extends React.Component {
  constructor(props) {
    super(props);
    this.onLogin = this.onLogin.bind(this);
  }

  onLogin(username, password) {
    this.loading.show('正在加载中');
    User.teacherLogin(username, password, () => {
      this.loading.hide();
      hashHistory.push('reservation');
    }, (error) => {
      this.loading.hide();
      setTimeout(() => {
        this.alert.show('登录失败', error, '好的');
      }, 500);
    });
  }

  render() {
    return (
      <div>
        <Panel access>
          <PanelHeader style={{ fontSize: '18px' }}>咨询师登录</PanelHeader>
          <PanelBody>
            <LoginForm
              titleTip="请输入工号、密码登录"
              usernameLabel="工号"
              usernamePlaceholder="请输入工号"
              passwordLabel="密码"
              passwordPlaceholder="请输入密码"
              submitText="登录"
              handleSubmit={this.onLogin}
            />
            <div style={{ color: '#999999', padding: '10px 20px', textAlign: 'center', fontSize: '13px' }}>
              账号密码遇到任何问题，请与管理员联络。
            </div>
          </PanelBody>
        </Panel>
        <LoadingHud ref={(loading) => { this.loading = loading; }} />
        <AlertDialog ref={(alert) => { this.alert = alert; }} />
        <PageBottom
          styles={{ color: '#999999', textAlign: 'center', backgroundColor: 'white', fontSize: '14px' }}
          contents={['清华大学学生心理发展指导中心', '联系方式：010-62782007']}
          height="55px"
        />
      </div>
    );
  }
}
