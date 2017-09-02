import 'weui';
import { AlertDialog, LoadingHud } from '#coms/Huds';
import { Panel, PanelBody, PanelHeader } from 'react-weui';
import ChangePasswordForm from '#forms/ChangePasswordForm';
import PageBottom from '#coms/PageBottom';
import PropTypes from 'prop-types';
import React from 'react';
import { User } from '#models/Models';

export default class TeacherChangePasswordPage extends React.Component {
  constructor(props) {
    super(props);
    this.handleCancel = this.handleCancel.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
    this.showAlert = this.showAlert.bind(this);
  }

  handleCancel() {
    this.props.history.goBack();
  }

  handleSubmit(oldPassword, newPassword) {
    this.loading.show('正在加载中');
    User.updateSession(() => {
      setTimeout(() => {
        User.teacherPasswordChange(User.username, oldPassword, newPassword, () => {
          this.loading.hide();
          this.alert.show('更改成功', '您已成功更改密码，请重新登录', '好的', () => {
            User.logout((data) => {
              if (data.redirect_url) {
                window.location.href = data.redirect_url;
              }
            });
          });
        }, (error) => {
          this.loading.hide();
          this.alert.show('更改失败', error, '好的');
        });
      }, 500);
    }, () => {
      this.props.history.push('/login');
    });
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
            <ChangePasswordForm
              titleTip="请输入原密码和新密码"
              oldPasswordLabel="原密码"
              oldPasswordPlaceholder="请输入原密码"
              newPasswordLabel="新密码"
              newPasswordPlaceholder="请输入新密码"
              newPasswordConfirmLabel="确认密码"
              newPasswordConfirmPlaceholder="请确认新密码"
              submitText="确认更改"
              cancelText="取消"
              handleSubmit={this.handleSubmit}
              handleCancel={this.handleCancel}
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

TeacherChangePasswordPage.propTypes = {
  history: PropTypes.object.isRequired,
};
