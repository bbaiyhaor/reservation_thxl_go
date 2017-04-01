import 'weui';
import { AlertDialog, LoadingHud } from '#coms/Huds';
import { Panel, PanelBody, PanelHeader } from '#react-weui';
import React, { PropTypes } from 'react';
import LoginForm from '#forms/LoginForm';
import PageBottom from '#coms/PageBottom';
import { User } from '#models/Models';

export default class TeacherLoginPage extends React.Component {
  constructor(props) {
    super(props);
    this.onLogin = this.onLogin.bind(this);
    this.toResetPassword = this.toResetPassword.bind(this);
  }

  onLogin(username, password) {
    this.loading.show('正在加载中');
    User.teacherLogin(username, password, () => {
      this.loading.hide();
      this.props.history.push('/reservation');
    }, (error) => {
      this.loading.hide();
      setTimeout(() => {
        this.alert.show('登录失败', error, '好的');
      }, 500);
    });
  }

  toResetPassword() {
    this.props.history.push('/password/reset');
  }

  render() {
    return (
      <div>
        <Panel>
          <PanelHeader style={{ fontSize: '18px' }}>咨询师登录</PanelHeader>
          <PanelBody>
            <LoginForm
              titleTip="请输入工号、密码登录"
              usernameLabel="工号"
              usernamePlaceholder="请输入工号"
              passwordLabel="密码"
              passwordPlaceholder="请输入密码"
              submitText="登录"
              cancelText="忘记密码"
              handleSubmit={this.onLogin}
              handleCancel={this.toResetPassword}
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

TeacherLoginPage.propTypes = {
  history: PropTypes.object.isRequired,
};
