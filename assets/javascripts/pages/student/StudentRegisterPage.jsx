import 'weui';
import { AlertDialog, LoadingHud } from '#coms/Huds';
import { Panel, PanelBody, PanelHeader } from '#react-weui';
import React, { PropTypes } from 'react';
import PageBottom from '#coms/PageBottom';
import RegisterForm from '#forms/RegisterForm';
import { User } from '#models/Models';

export default class StudentRegisterPage extends React.Component {
  constructor(props) {
    super(props);
    this.onRegister = this.onRegister.bind(this);
    this.toLogin = this.toLogin.bind(this);
    this.showAlert = this.showAlert.bind(this);
  }

  onRegister(username, password, fullname, gender, birthday, school, grade, currentAddress, familyAddress, mobile, email, experienceTime, experienceLocation, experienceTeacher, fatherAge, fatherJob, fatherEdu, motherAge, motherJob, motherEdu, parentMarriage) {
    const payload = {
      username,
      password,
      student_fullname: fullname,
      student_gender: gender,
      student_birthday: birthday,
      student_school: school,
      student_grade: grade,
      student_current_address: currentAddress,
      student_family_address: familyAddress,
      student_mobile: mobile,
      student_email: email,
      student_experience_time: experienceTime,
      student_experience_location: experienceLocation,
      student_experience_teacher: experienceTeacher,
      student_father_age: fatherAge,
      student_father_job: fatherJob,
      student_father_edu: fatherEdu,
      student_mother_age: motherAge,
      student_mother_job: motherJob,
      student_mother_edu: motherEdu,
      student_parent_marriage: parentMarriage,
    };
    this.loading.show('正在加载中');
    User.studentRegister(payload, () => {
      this.loading.hide();
      this.alert.show('注册成功', '请用学号和密码登录', '好的', () => {
        this.props.history.push('/login');
      });
    }, (error) => {
      this.loading.hide();
      setTimeout(() => {
        this.alert.show('注册失败', error, '好的');
      }, 500);
    });
  }

  toLogin() {
    this.props.history.push('/login');
  }

  showAlert(title, msg, label) {
    this.alert.show(title, msg, label);
  }

  render() {
    return (
      <div>
        <Panel>
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
              protocolLink="/protocol"
              submitText="注册"
              cancelText="已有账户"
              handleSubmit={this.onRegister}
              handleCancel={this.toLogin}
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
          height="55px"
        />
      </div>
    );
  }
}

StudentRegisterPage.propTypes = {
  history: PropTypes.object.isRequired,
};
