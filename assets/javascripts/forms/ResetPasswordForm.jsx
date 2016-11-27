/**
 * Created by shudi on 2016/11/6.
 */
/* eslint consistent-return: "off" */
import React, { PropTypes } from 'react';
import { CellsTitle, Form, FormCell, CellHeader, Label, CellBody, Input, CellFooter, Icon, ButtonArea, Button } from '#react-weui';
import 'weui';

import VerifyCodeInput from '#coms/VerifyCodeInput';
import { User } from '#models/Models';

const propTypes = {
  titleTip: PropTypes.string,
  usernameLabel: PropTypes.string.isRequired,
  usernamePlaceholder: PropTypes.string.isRequired,
  fullnameLabel: PropTypes.string.isRequired,
  fullnamePlaceholder: PropTypes.string.isRequired,
  mobileLabel: PropTypes.string.isRequired,
  mobilePlaceholder: PropTypes.string.isRequired,
  verifyCodeLabel: PropTypes.string.isRequired,
  verifyCodePlaceholder: PropTypes.string.isRequired,
  newPasswordLabel: PropTypes.string.isRequired,
  newPasswordPlaceholder: PropTypes.string.isRequired,
  newPasswordConfirmLabel: PropTypes.string.isRequired,
  newPasswordConfirmPlaceholder: PropTypes.string.isRequired,
  submitText: PropTypes.string.isRequired,
  cancelText: PropTypes.string,
  handleSubmit: PropTypes.func.isRequired,
  handleCancel: PropTypes.func,
  showAlert: PropTypes.func,
};

export default class ResetPasswordForm extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      username: '',
      fullname: '',
      mobile: '',
      newPassword: '',
      newPasswordConfirm: '',
      usernameWarn: false,
      fullnameWarn: false,
      mobileWarn: false,
      newPasswordWarn: false,
      newPasswordConfirmWarn: false,
    };
    this.handleChange = this.handleChange.bind(this);
    this.sendSms = this.sendSms.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  handleChange(e, name) {
    this.setState({ [name]: e.target.value });
  }

  sendSms() {
    this.setState({
      usernameWarn: false,
      fullnameWarn: false,
      mobileWarn: false,
    });
    if (this.state.username === '') {
      this.setState({ usernameWarn: true });
      this.usernameInput.focus();
      return;
    }
    if (this.state.fullname === '') {
      this.setState({ fullnameWarn: true });
      this.fullnameInput.focus();
      return;
    }
    if (this.state.mobile === '') {
      this.setState({ mobileWarn: true });
      this.mobileInput.focus();
      return;
    }
    User.teacherPasswordResetSms(this.state.username, this.state.fullname, this.state.mobile, () => {
      this.verifyCodeInput.restart();
    }, (error) => {
      this.props.showAlert('发送失败', error, '好的');
    });
  }

  handleSubmit() {
    this.setState({
      usernameWarn: false,
      fullnameWarn: false,
      mobileWarn: false,
      newPasswordWarn: false,
      newPasswordConfirmWarn: false,
    });
    this.verifyCodeInput.resetWarn();
    if (this.state.username === '') {
      this.setState({ usernameWarn: true });
      this.usernameInput.focus();
      return;
    }
    if (this.verifyCodeInput.getValue() === '') {
      this.verifyCodeInput.setInvalidate();
      this.verifyCodeInput.focus();
      return;
    }
    if (this.state.newPassword === '') {
      this.setState({ newPasswordWarn: true });
      this.newPasswordInput.focus();
      return;
    }
    if (this.state.newPasswordConfirm === '') {
      this.setState({ newPasswordConfirmWarn: true });
      this.newPasswordConfirmInput.focus();
      return;
    }
    if (this.state.newPassword !== this.state.newPasswordConfirm) {
      this.setState({
        newPasswordWarn: true,
        newPasswordConfirmWarn: true,
      });
      this.props.showAlert && this.props.showAlert('重置失败', '两次密码不一致，请重新输入', '好的');
      this.setState({
        newPassword: '',
        newPasswordConfirm: '',
      });
      return;
    }
    this.props.handleSubmit(this.state.username, this.state.newPassword, this.verifyCodeInput.getValue());
  }

  render() {
    return (
      <div>
        {this.props.titleTip && this.props.titleTip !== '' &&
          <CellsTitle>{this.props.titleTip}</CellsTitle>
        }
        <Form>
          <FormCell warn={this.state.usernameWarn}>
            <CellHeader>
              <Label>{this.props.usernameLabel}</Label>
            </CellHeader>
            <CellBody>
              <Input
                ref={(usernameInput) => { this.usernameInput = usernameInput; }}
                type="input"
                placeholder={this.props.usernamePlaceholder}
                value={this.state.username}
                onChange={(e) => { this.handleChange(e, 'username'); }}
              />
            </CellBody>
            {this.state.usernameWarn &&
            <CellFooter>
              <Icon value="warn" />
            </CellFooter>
            }
          </FormCell>
          <FormCell warn={this.state.fullnameWarn}>
            <CellHeader>
              <Label>{this.props.fullnameLabel}</Label>
            </CellHeader>
            <CellBody>
              <Input
                ref={(fullnameInput) => { this.fullnameInput = fullnameInput; }}
                type="input"
                placeholder={this.props.fullnamePlaceholder}
                value={this.state.fullname}
                onChange={(e) => { this.handleChange(e, 'fullname'); }}
              />
            </CellBody>
            {this.state.fullnameWarn &&
            <CellFooter>
              <Icon value="warn" />
            </CellFooter>
            }
          </FormCell>
          <FormCell warn={this.state.mobileWarn}>
            <CellHeader>
              <Label>{this.props.mobileLabel}</Label>
            </CellHeader>
            <CellBody>
              <Input
                ref={(mobileInput) => { this.mobileInput = mobileInput; }}
                type="input"
                placeholder={this.props.mobilePlaceholder}
                value={this.state.mobile}
                onChange={(e) => { this.handleChange(e, 'mobile'); }}
              />
            </CellBody>
            {this.state.mobileWarn &&
            <CellFooter>
              <Icon value="warn" />
            </CellFooter>
            }
          </FormCell>
          <VerifyCodeInput
            ref={(verifyCodeInput) => { this.verifyCodeInput = verifyCodeInput; }}
            countDown
            verifyCodeLabel={this.props.verifyCodeLabel}
            verifyCodePlaceholder={this.props.verifyCodePlaceholder}
            handleClick={this.sendSms}
          />
          <FormCell warn={this.state.newPasswordWarn}>
            <CellHeader>
              <Label>{this.props.newPasswordLabel}</Label>
            </CellHeader>
            <CellBody>
              <Input
                ref={(newPasswordInput) => { this.newPasswordInput = newPasswordInput; }}
                type="password"
                placeholder={this.props.newPasswordPlaceholder}
                value={this.state.newPassword}
                onChange={(e) => { this.handleChange(e, 'newPassword'); }}
              />
            </CellBody>
            {this.state.newPasswordWarn &&
            <CellFooter>
              <Icon value="warn" />
            </CellFooter>
            }
          </FormCell>
          <FormCell warn={this.state.newPasswordConfirmWarn}>
            <CellHeader>
              <Label>{this.props.newPasswordConfirmLabel}</Label>
            </CellHeader>
            <CellBody>
              <Input
                ref={(newPasswordConfirmInput) => { this.newPasswordConfirmInput = newPasswordConfirmInput; }}
                type="password"
                placeholder={this.props.newPasswordConfirmPlaceholder}
                value={this.state.newPasswordConfirm}
                onChange={(e) => { this.handleChange(e, 'newPasswordConfirm'); }}
              />
            </CellBody>
            {this.state.newPasswordConfirmWarn &&
            <CellFooter>
              <Icon value="warn" />
            </CellFooter>
            }
          </FormCell>
        </Form>
        <ButtonArea direction="horizontal">
          <Button onClick={this.handleSubmit}>{this.props.submitText}</Button>
          {this.props.cancelText && <Button type="default" onClick={this.props.handleCancel}>{this.props.cancelText}</Button>}
        </ButtonArea>
      </div>
    );
  }
}

ResetPasswordForm.propTypes = propTypes;
