import 'weui';
import { Button, ButtonArea, CellBody, CellFooter, CellHeader, CellsTitle, Form, FormCell, Icon, Input, Label } from 'react-weui';
import React, { PropTypes } from 'react';

export default class ChangePasswordForm extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      oldPassword: '',
      newPassword: '',
      newPasswordConfirm: '',
      oldPasswordWarn: false,
      newPasswordWarn: false,
      newPasswordConfirmWarn: false,
    };
    this.handleChange = this.handleChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  handleChange(e, name) {
    this.setState({ [name]: e.target.value });
  }

  handleSubmit() {
    this.setState({
      oldPasswordWarn: false,
      newPasswordWarn: false,
      newPasswordConfirmWarn: false,
    });
    if (this.state.oldPassword === '') {
      this.setState({ oldPasswordWarn: true });
      this.oldPasswordInput.focus();
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
      this.props.showAlert && this.props.showAlert('注册失败', '两次密码不一致，请重新输入', '好的');
      this.setState({
        newPassword: '',
        newPasswordConfirm: '',
      });
      return;
    }
    this.props.handleSubmit(this.state.oldPassword, this.state.newPassword);
  }

  render() {
    return (
      <div>
        {this.props.titleTip && this.props.titleTip !== '' &&
          <CellsTitle>{this.props.titleTip}</CellsTitle>
        }
        <Form>
          <FormCell warn={this.state.oldPasswordWarn}>
            <CellHeader>
              <Label>{this.props.oldPasswordLabel}</Label>
            </CellHeader>
            <CellBody>
              <Input
                ref={(oldPasswordInput) => { this.oldPasswordInput = oldPasswordInput; }}
                type="password"
                placeholder={this.props.oldPasswordPlaceholder}
                value={this.state.oldPassword}
                onChange={(e) => { this.handleChange(e, 'oldPassword'); }}
              />
            </CellBody>
            {this.state.oldPasswordWarn &&
            <CellFooter>
              <Icon value="warn" />
            </CellFooter>
            }
          </FormCell>
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

ChangePasswordForm.propTypes = {
  titleTip: PropTypes.string,
  oldPasswordLabel: PropTypes.string.isRequired,
  oldPasswordPlaceholder: PropTypes.string,
  newPasswordLabel: PropTypes.string.isRequired,
  newPasswordPlaceholder: PropTypes.string,
  newPasswordConfirmLabel: PropTypes.string.isRequired,
  newPasswordConfirmPlaceholder: PropTypes.string,
  submitText: PropTypes.string.isRequired,
  cancelText: PropTypes.string,
  handleSubmit: PropTypes.func.isRequired,
  handleCancel: PropTypes.func,
  showAlert: PropTypes.func,
};

ChangePasswordForm.defaultProps = {
  titleTip: '',
  oldPasswordPlaceholder: '',
  newPasswordPlaceholder: '',
  newPasswordConfirmPlaceholder: '',
  cancelText: '',
  handleCancel: undefined,
  showAlert: undefined,
};
