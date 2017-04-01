import 'weui';
import { AlertDialog, LoadingHud } from '#coms/Huds';
import { Panel, PanelBody, PanelHeader } from '#react-weui';
import React, { PropTypes } from 'react';
import PageBottom from '#coms/PageBottom';
import ResetPasswordForm from '#forms/ResetPasswordForm';
import { User } from '#models/Models';

export default class TeacherResetPasswordPage extends React.Component {
  constructor(props) {
    super(props);
    this.toLogin = this.toLogin.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
    this.showAlert = this.showAlert.bind(this);
  }

  toLogin() {
    this.props.history.push('/login');
  }

  handleSubmit(username, newPassword, verifyCode) {
    this.loading.show('正在加载中');
    setTimeout(() => {
      User.teacherPasswordResetVerify(username, newPassword, verifyCode, () => {
        this.loading.hide();
        this.alert.show('重置成功', '您已成功重置密码，请重新登录', '好的', () => {
          this.props.history.push('/login');
        });
      }, (error) => {
        this.loading.hide();
        this.alert.show('重置失败', error, '好的');
      });
    }, 500);
  }

  showAlert(title, msg, label) {
    this.alert.show(title, msg, label);
  }

  render() {
    return (
      <div>
        <Panel>
          <PanelHeader style={{ fontSize: '18px' }}>咨询师更改密码</PanelHeader>
          <PanelBody>
            <ResetPasswordForm
              usernameLabel="工号"
              usernamePlaceholder="请输入工号"
              fullnameLabel="姓名"
              fullnamePlaceholder="请输入姓名"
              mobileLabel="手机号"
              mobilePlaceholder="请输入手机号"
              verifyCodeLabel="验证码"
              verifyCodePlaceholder="请输入验证码"
              newPasswordLabel="新密码"
              newPasswordPlaceholder="请输入新密码"
              newPasswordConfirmLabel="确认密码"
              newPasswordConfirmPlaceholder="请确认新密码"
              submitText="重置密码"
              cancelText="返回登录"
              handleSubmit={this.handleSubmit}
              handleCancel={this.toLogin}
              showAlert={this.showAlert}
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

TeacherResetPasswordPage.propTypes = {
  history: PropTypes.object.isRequired,
};
