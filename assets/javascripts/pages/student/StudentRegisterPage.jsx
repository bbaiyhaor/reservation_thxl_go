/**
 * Created by shudi on 2016/10/22.
 */
import React from 'react';
import { hashHistory } from 'react-router';
import { Panel, PanelHeader, PanelBody } from '#react-weui';
import 'weui';

import RegisterForm from '#forms/RegisterForm';
import PageBottom from '#coms/PageBottom';
import { AlertDialog, LoadingHud } from '#coms/Huds';
import { User } from '#models/Models';

export default class StudentRegisterPage extends React.Component {
  static toLogin() {
    hashHistory.push('login');
  }

  constructor(props) {
    super(props);
    this.onRegister = this.onRegister.bind(this);
    this.showAlert = this.showAlert.bind(this);
  }

  onRegister(username, password) {
    this.loading.show('正在加载中');
    User.studentRegister(username, password, () => {
      this.loading.hide();
      hashHistory.push('reservation');
    }, (error) => {
      this.loading.hide();
      setTimeout(() => {
        this.alert.show('注册失败', error, '好的');
      }, 500);
    });
  }

  showAlert(title, msg, label) {
    this.alert.show(title, msg, label);
  }

  render() {
    return (
      <div>
        <Panel access>
          <PanelHeader style={{ fontSize: '18px' }}>学生注册</PanelHeader>
          <PanelBody>
            <RegisterForm
              titleTip="请用学号注册（密码与info账号不同）"
              usernameLabel="学号"
              usernamePlaceholder="请输入学号"
              passwordLabel="密码"
              passwordPlaceholder="请输入密码"
              confirmPasswordLabel="确认密码"
              confirmPasswordPlaceholder="请确认密码"
              protocol="咨询协议"
              protocolPrefix="我已阅读并同意"
              protocolLink="protocol"
              submitText="注册"
              cancelText="已有账户"
              handleSubmit={this.onRegister}
              handleCancel={StudentRegisterPage.toLogin}
              showAlert={this.showAlert}
            />
            <div style={{ color: '#999999', padding: '10px 20px', textAlign: 'center', fontSize: '13px' }}>
            账号密码遇到任何问题，请在工作时间（周一至周五8:00-12:00；13:00-17:00）致电010-62782007。
            </div>
          </PanelBody>
        </Panel>
        <LoadingHud ref={(loading) => { this.loading = loading; }} />
        <AlertDialog ref={(alert) => { this.alert = alert; }} />
        <PageBottom
          styles={{ color: '#999999', textAlign: 'center', backgroundColor: 'white', fontSize: '14px' }}
          contents={['清华大学学生心理发展指导中心', '联系方式：010-62782007']}
          height="55  px"
        />
      </div>
    );
  }
}
