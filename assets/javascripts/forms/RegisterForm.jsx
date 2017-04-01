import 'weui';
import { Agreement, Button, ButtonArea, CellBody, CellFooter, CellHeader, CellsTitle, Form, FormCell, Icon, Input, Label, Select } from '#react-weui';
import React, { PropTypes } from 'react';
import { Link } from 'react-router-dom';

export default class RegisterForm extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      username: '',
      password: '',
      confirmPassword: '',
      fullname: '',
      gender: '',
      birthday: '',
      school: '',
      grade: '',
      currentAddress: '',
      familyAddress: '',
      mobile: '',
      email: '',
      experienceTime: '',
      experienceLocation: '',
      experienceTeacher: '',
      fatherAge: '',
      fatherJob: '',
      fatherEdu: '',
      motherAge: '',
      motherJob: '',
      motherEdu: '',
      parentMarriage: '',
      protocolChecked: true,
      usernameWarn: false,
      passwordWarn: false,
      confirmPasswordWarn: false,
      fullnameWarn: false,
      genderWarn: false,
      birthdayWarn: false,
      schoolWarn: false,
      gradeWarn: false,
      currentAddressWarn: false,
      familyAddressWarn: false,
      mobileWarn: false,
      emailWarn: false,
    };
    this.handleChange = this.handleChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  handleChange(e) {
    const name = e.target.name;
    const value = e.target.type === 'checkbox' ? e.target.checked : e.target.value;
    this.setState({ [name]: value });
  }

  handleSubmit() {
    this.setState({
      usernameWarn: false,
      passwordWarn: false,
      confirmPasswordWarn: false,
      fullnameWarn: false,
      genderWarn: false,
      birthdayWarn: false,
      schoolWarn: false,
      gradeWarn: false,
      currentAddressWarn: false,
      familyAddressWarn: false,
      mobileWarn: false,
      emailWarn: false,
    });
    if (this.state.username === '') {
      this.setState({ usernameWarn: true });
      this.usernameInput.focus();
      return;
    }
    if (this.state.password === '') {
      this.setState({ passwordWarn: true });
      this.passwordInput.focus();
      return;
    }
    if (this.state.confirmPassword === '') {
      this.setState({ confirmPasswordWarn: true });
      this.confirmPasswordInput.focus();
      return;
    }
    if (this.state.fullname === '') {
      this.setState({ fullnameWarn: true });
      this.fullnameInput.focus();
      return;
    }
    if (this.state.gender === '') {
      this.setState({ genderWarn: true });
      this.genderSelect.focus();
      return;
    }
    if (this.state.birthday === '') {
      this.setState({ birthdayWarn: true });
      this.birthdayInput.focus();
      return;
    }
    if (this.state.school === '') {
      this.setState({ schoolWarn: true });
      this.schoolInput.focus();
      return;
    }
    if (this.state.grade === '') {
      this.setState({ gradeWarn: true });
      this.gradeInput.focus();
      return;
    }
    if (this.state.currentAddress === '') {
      this.setState({ currentAddressWarn: true });
      this.currentAddressInput.focus();
      return;
    }
    if (this.state.familyAddress === '') {
      this.setState({ familyAddressWarn: true });
      this.familyAddressInput.focus();
      return;
    }
    if (this.state.mobile === '') {
      this.setState({ mobileWarn: true });
      this.mobileInput.focus();
      return;
    }
    if (this.state.email === '') {
      this.setState({ emailWarn: true });
      this.emailInput.focus();
      return;
    }
    if (!this.state.protocolChecked) {
      this.props.showAlert && this.props.showAlert('注册失败', '请阅读并同意咨询协议', '好的');
      return;
    }
    if (this.state.password !== this.state.confirmPassword) {
      this.setState({
        passwordWarn: true,
        confirmPasswordWarn: true,
      });
      this.props.showAlert && this.props.showAlert('注册失败', '两次密码不一致，请重新输入', '好的');
      this.setState({
        password: '',
        confirmPassword: '',
      });
      return;
    }
    this.props.handleSubmit(this.state.username, this.state.password, this.state.fullname, this.state.gender, this.state.birthday, this.state.school, this.state.grade, this.state.currentAddress, this.state.familyAddress, this.state.mobile, this.state.email, this.state.experienceTime, this.state.experienceLocation, this.state.experienceTeacher, this.state.fatherAge, this.state.fatherJob, this.state.fatherEdu, this.state.motherAge, this.state.motherJob, this.state.motherEdu, this.state.parentMarriage);
  }

  render() {
    return (
      <div>
        {this.props.titleTip && this.props.titleTip !== '' &&
          <CellsTitle>{this.props.titleTip}</CellsTitle>
        }
        <Form className="weui-cells_form">
          <FormCell warn={this.state.usernameWarn}>
            <CellHeader>
              <Label>{this.props.usernameLabel}<span style={{ color: 'red' }}>*</span></Label>
            </CellHeader>
            <CellBody>
              <Input
                name="username"
                ref={(usernameInput) => { this.usernameInput = usernameInput; }}
                type="tel"
                placeholder={this.props.usernamePlaceholder}
                value={this.state.username}
                onChange={this.handleChange}
              />
            </CellBody>
            {this.state.usernameWarn &&
              <CellFooter>
                <Icon size="small" value="warn" />
              </CellFooter>
            }
          </FormCell>
          <FormCell warn={this.state.passwordWarn}>
            <CellHeader>
              <Label>{this.props.passwordLabel}<span style={{ color: 'red' }}>*</span></Label>
            </CellHeader>
            <CellBody>
              <Input
                name="password"
                ref={(passwordInput) => { this.passwordInput = passwordInput; }}
                type="password"
                placeholder={this.props.passwordPlaceholder}
                value={this.state.password}
                onChange={this.handleChange}
              />
            </CellBody>
            {this.state.passwordWarn &&
              <CellFooter>
                <Icon size="small" value="warn" />
              </CellFooter>
            }
          </FormCell>
          <FormCell warn={this.state.confirmPasswordWarn}>
            <CellHeader>
              <Label>{this.props.confirmPasswordLabel}<span style={{ color: 'red' }}>*</span></Label>
            </CellHeader>
            <CellBody>
              <Input
                name="confirmPassword"
                ref={(confirmPasswordInput) => { this.confirmPasswordInput = confirmPasswordInput; }}
                type="password"
                placeholder={this.props.confirmPasswordPlaceholder}
                value={this.state.confirmPassword}
                onChange={this.handleChange}
              />
            </CellBody>
            {this.state.confirmPasswordWarn &&
              <CellFooter>
                <Icon size="small" value="warn" />
              </CellFooter>
            }
          </FormCell>
          <FormCell warn={this.state.fullnameWarn}>
            <CellHeader>
              <Label>姓名<span style={{ color: 'red' }}>*</span></Label>
            </CellHeader>
            <CellBody>
              <Input
                name="fullname"
                ref={(fullnameInput) => { this.fullnameInput = fullnameInput; }}
                type="input"
                placeholder="请输入姓名"
                value={this.state.fullname}
                onChange={this.handleChange}
              />
            </CellBody>
            {this.state.fullnameWarn &&
              <CellFooter>
                <Icon size="small" value="warn" />
              </CellFooter>
            }
          </FormCell>
          <FormCell warn={this.state.genderWarn} select selectPos="after">
            <CellHeader>
              <Label>性别<span style={{ color: 'red' }}>*</span></Label>
            </CellHeader>
            <CellBody>
              <Select
                name="gender"
                ref={(genderSelect) => { this.genderSelect = genderSelect; }}
                value={this.state.gender}
                onChange={this.handleChange}
              >
                <option value="">请选择</option>
                <option value="男">男</option>
                <option value="女">女</option>
              </Select>
            </CellBody>
            {this.state.genderWarn &&
              <CellFooter style={{ marginRight: '25px' }}>
                <Icon size="small" value="warn" />
              </CellFooter>
            }
          </FormCell>
          <FormCell warn={this.state.birthdayWarn}>
            <CellHeader>
              <Label>出生日期<span style={{ color: 'red' }}>*</span></Label>
            </CellHeader>
            <CellBody>
              <Input
                name="birthday"
                ref={(birthdayInput) => { this.birthdayInput = birthdayInput; }}
                type="date"
                placeholder="请输入出生日期"
                value={this.state.birthday}
                onChange={this.handleChange}
              />
            </CellBody>
            {this.state.birthdayWarn &&
              <CellFooter>
                <Icon size="small" value="warn" />
              </CellFooter>
            }
          </FormCell>
          <FormCell warn={this.state.schoolWarn}>
            <CellHeader>
              <Label>院系<span style={{ color: 'red' }}>*</span></Label>
            </CellHeader>
            <CellBody>
              <Input
                name="school"
                ref={(schoolInput) => { this.schoolInput = schoolInput; }}
                type="input"
                placeholder="请输入院系"
                value={this.state.school}
                onChange={this.handleChange}
              />
            </CellBody>
            {this.state.schoolWarn &&
              <CellFooter>
                <Icon size="small" value="warn" />
              </CellFooter>
            }
          </FormCell>
          <FormCell warn={this.state.gradeWarn}>
            <CellHeader>
              <Label>年级<span style={{ color: 'red' }}>*</span></Label>
            </CellHeader>
            <CellBody>
              <Input
                name="grade"
                ref={(gradeInput) => { this.gradeInput = gradeInput; }}
                type="input"
                placeholder="请输入年级"
                value={this.state.grade}
                onChange={this.handleChange}
              />
            </CellBody>
            {this.state.gradeWarn &&
              <CellFooter>
                <Icon size="small" value="warn" />
              </CellFooter>
            }
          </FormCell>
          <FormCell warn={this.state.currentAddressWarn}>
            <CellHeader>
              <Label>现在住址<span style={{ color: 'red' }}>*</span></Label>
            </CellHeader>
            <CellBody>
              <Input
                name="currentAddress"
                ref={(currentAddressInput) => { this.currentAddressInput = currentAddressInput; }}
                type="input"
                placeholder="请输入现在住址"
                value={this.state.currentAddress}
                onChange={this.handleChange}
              />
            </CellBody>
            {this.state.currentAddressWarn &&
              <CellFooter>
                <Icon size="small" value="warn" />
              </CellFooter>
            }
          </FormCell>
          <FormCell warn={this.state.familyAddressWarn}>
            <CellHeader>
              <Label>家庭住址<span style={{ color: 'red' }}>*</span></Label>
            </CellHeader>
            <CellBody>
              <Input
                name="familyAddress"
                ref={(familyAddressInput) => { this.familyAddressInput = familyAddressInput; }}
                type="input"
                placeholder="请输入家庭住址"
                value={this.state.familyAddress}
                onChange={this.handleChange}
              />
            </CellBody>
            {this.state.familyAddressWarn &&
              <CellFooter>
                <Icon size="small" value="warn" />
              </CellFooter>
            }
          </FormCell>
          <FormCell warn={this.state.mobileWarn}>
            <CellHeader>
              <Label>联系电话<span style={{ color: 'red' }}>*</span></Label>
            </CellHeader>
            <CellBody>
              <Input
                name="mobile"
                ref={(mobileInput) => { this.mobileInput = mobileInput; }}
                type="tel"
                placeholder="请输入联系电话"
                value={this.state.mobile}
                onChange={this.handleChange}
              />
            </CellBody>
            {this.state.mobileWarn &&
              <CellFooter>
                <Icon size="small" value="warn" />
              </CellFooter>
            }
          </FormCell>
          <FormCell warn={this.state.emailWarn}>
            <CellHeader>
              <Label>邮箱<span style={{ color: 'red' }}>*</span></Label>
            </CellHeader>
            <CellBody>
              <Input
                name="email"
                ref={(emailInput) => { this.emailInput = emailInput; }}
                type="input"
                placeholder="请输入邮箱"
                value={this.state.email}
                onChange={this.handleChange}
              />
            </CellBody>
            {this.state.emailWarn &&
              <CellFooter>
                <Icon size="small" value="warn" />
              </CellFooter>
            }
          </FormCell>
          <CellsTitle>过往咨询经历</CellsTitle>
          <FormCell>
            <CellHeader>
              <Label>时间</Label>
            </CellHeader>
            <CellBody>
              <Input
                name="experienceTime"
                type="input"
                placeholder="请输入咨询时间"
                value={this.state.experienceTime}
                onChange={this.handleChange}
              />
            </CellBody>
          </FormCell>
          <FormCell>
            <CellHeader>
              <Label>地点</Label>
            </CellHeader>
            <CellBody>
              <Input
                name="experienceLocation"
                type="input"
                placeholder="请输入咨询地点"
                value={this.state.experienceLocation}
                onChange={this.handleChange}
              />
            </CellBody>
          </FormCell>
          <FormCell>
            <CellHeader>
              <Label>咨询师姓名</Label>
            </CellHeader>
            <CellBody>
              <Input
                name="experienceTeacher"
                type="input"
                placeholder="请输入咨询师姓名"
                value={this.state.experienceTeacher}
                onChange={this.handleChange}
              />
            </CellBody>
          </FormCell>
          <CellsTitle>家庭情况</CellsTitle>
          <FormCell>
            <CellHeader>
              <Label>父亲年龄</Label>
            </CellHeader>
            <CellBody>
              <Input
                name="fatherAge"
                type="input"
                placeholder="请输入父亲年龄"
                value={this.state.fatherAge}
                onChange={this.handleChange}
              />
            </CellBody>
          </FormCell>
          <FormCell>
            <CellHeader>
              <Label>父亲职业</Label>
            </CellHeader>
            <CellBody>
              <Input
                name="fatherJob"
                type="input"
                placeholder="请输入父亲职业"
                value={this.state.fatherJob}
                onChange={this.handleChange}
              />
            </CellBody>
          </FormCell>
          <FormCell>
            <CellHeader>
              <Label>父亲学历</Label>
            </CellHeader>
            <CellBody>
              <Input
                name="fatherEdu"
                type="input"
                placeholder="请输入父亲学历"
                value={this.state.fatherEdu}
                onChange={this.handleChange}
              />
            </CellBody>
          </FormCell>
          <FormCell>
            <CellHeader>
              <Label>母亲年龄</Label>
            </CellHeader>
            <CellBody>
              <Input
                name="motherAge"
                type="input"
                placeholder="请输入母亲年龄"
                value={this.state.motherAge}
                onChange={this.handleChange}
              />
            </CellBody>
          </FormCell>
          <FormCell>
            <CellHeader>
              <Label>母亲职业</Label>
            </CellHeader>
            <CellBody>
              <Input
                name="motherJob"
                type="input"
                placeholder="请输入母亲职业"
                value={this.state.motherJob}
                onChange={this.handleChange}
              />
            </CellBody>
          </FormCell>
          <FormCell>
            <CellHeader>
              <Label>母亲学历</Label>
            </CellHeader>
            <CellBody>
              <Input
                name="motherEdu"
                type="input"
                placeholder="请输入母亲学历"
                value={this.state.motherEdu}
                onChange={this.handleChange}
              />
            </CellBody>
          </FormCell>
          <FormCell select selectPos="after">
            <CellHeader>
              <Label>父母婚姻状况</Label>
            </CellHeader>
            <CellBody>
              <Select
                name="parentMarriage"
                value={this.state.parentMarriage}
                onChange={this.handleChange}
              >
                <option value="">请选择</option>
                <option value="良好">良好</option>
                <option value="一般">一般</option>
                <option value="离婚">离婚</option>
                <option value="再婚">再婚</option>
              </Select>
            </CellBody>
          </FormCell>
        </Form>
        {this.props.protocol ?
          <Agreement
            name="protocolChecked"
            defaultChecked
            value="protoclChecked"
            onChange={this.handleChange}
          >
            {this.props.protocolPrefix}
            <Link to={this.props.protocolLink} style={{ pointerEvents: 'all' }}>{this.props.protocol}</Link>
            {this.props.protocolSuffix}
          </Agreement> : null
        }
        <ButtonArea direction="horizontal">
          <Button onClick={this.handleSubmit}>{this.props.submitText}</Button>
          <Button type="default" onClick={this.props.handleCancel}>{this.props.cancelText}</Button>
        </ButtonArea>
        <div style={{ height: '10px' }} />
      </div>
    );
  }
}

RegisterForm.propTypes = {
  titleTip: PropTypes.string,
  usernameLabel: PropTypes.string.isRequired,
  usernamePlaceholder: PropTypes.string,
  passwordLabel: PropTypes.string.isRequired,
  passwordPlaceholder: PropTypes.string,
  confirmPasswordLabel: PropTypes.string.isRequired,
  confirmPasswordPlaceholder: PropTypes.string,
  protocol: PropTypes.string,
  protocolPrefix: PropTypes.string,
  protocolSuffix: PropTypes.string,
  protocolLink: PropTypes.string,
  submitText: PropTypes.string.isRequired,
  cancelText: PropTypes.string.isRequired,
  handleSubmit: PropTypes.func.isRequired,
  handleCancel: PropTypes.func.isRequired,
  showAlert: PropTypes.func,
};

RegisterForm.defaultProps = {
  titleTip: '',
  usernamePlaceholder: '',
  passwordPlaceholder: '',
  confirmPasswordPlaceholder: '',
  protocol: '',
  protocolPrefix: '',
  protocolSuffix: '',
  protocolLink: '',
  showAlert: undefined,
};
