/**
 * Created by shudi on 2016/10/22.
 */
import React, { PropTypes } from 'react';
import { Link } from 'react-router';
import { CellsTitle, Form, FormCell, CellHeader, Label, CellBody, Input, CellFooter, Icon, ButtonArea, Button, Checkbox } from '#react-weui';
import 'weui';

const propTypes = {
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

export default class RegisterForm extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      username: '',
      password: '',
      confirmPassword: '',
      protocolChecked: true,
      usernameWarn: false,
      passwordWarn: false,
      confirmPasswordWarn: false,
      protocolWarn: false,
    };
    this.handleChange = this.handleChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  handleChange(e, name) {
    const value = e.target.value;
    if (name && name !== '') {
      this.setState({ [name]: value });
    } else {
      this.setState(prevState => ({
        [value]: !prevState[value],
      }));
    }
  }

  handleSubmit() {
    this.setState({
      usernameWarn: false,
      passwordWarn: false,
      confirmPasswordWarn: false,
      protocolWarn: false,
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
    if (!this.state.protocolChecked) {
      this.setState({ protocolWarn: true });
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
    this.props.handleSubmit(this.state.username, this.state.password);
  }

  render() {
    return (
      <div>
        {this.props.titleTip && this.props.titleTip !== '' &&
          <CellsTitle>{this.props.titleTip}</CellsTitle>
        }
        <Form checkbox={this.props.protocol && true} className="weui_cells_form">
          <FormCell warn={this.state.usernameWarn}>
            <CellHeader>
              <Label>{this.props.usernameLabel}</Label>
            </CellHeader>
            <CellBody>
              <Input
                ref={(usernameInput) => { this.usernameInput = usernameInput; }}
                type="tel"
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
          <FormCell warn={this.state.passwordWarn}>
            <CellHeader>
              <Label>{this.props.passwordLabel}</Label>
            </CellHeader>
            <CellBody>
              <Input
                ref={(passwordInput) => { this.passwordInput = passwordInput; }}
                type="password"
                placeholder={this.props.passwordPlaceholder}
                value={this.state.password}
                onChange={(e) => { this.handleChange(e, 'password'); }}
              />
            </CellBody>
            {this.state.passwordWarn &&
              <CellFooter>
                <Icon value="warn" />
              </CellFooter>
            }
          </FormCell>
          <FormCell warn={this.state.confirmPasswordWarn}>
            <CellHeader>
              <Label>{this.props.confirmPasswordLabel}</Label>
            </CellHeader>
            <CellBody>
              <Input
                ref={(confirmPasswordInput) => { this.confirmPasswordInput = confirmPasswordInput; }}
                type="password"
                placeholder={this.props.confirmPasswordPlaceholder}
                value={this.state.confirmPassword}
                onChange={(e) => { this.handleChange(e, 'confirmPassword'); }}
              />
            </CellBody>
            {this.state.confirmPasswordWarn &&
              <CellFooter>
                <Icon value="warn" />
              </CellFooter>
            }
          </FormCell>
          {this.props.protocol ?
            <FormCell checkbox warn={this.state.protocolWarn}>
              <CellHeader>
                <Checkbox
                  defaultChecked
                  value="protocolChecked"
                  onChange={this.handleChange}
                />
              </CellHeader>
              <CellBody>
                {this.props.protocolPrefix}
                <Link to={this.props.protocolLink}>{this.props.protocol}</Link>
                {this.props.protocolSuffix}
              </CellBody>
              {this.state.protocolWarn &&
              <CellFooter>
                <Icon value="warn" />
              </CellFooter>
              }
            </FormCell> : null
          }
        </Form>
        <ButtonArea direction="horizontal">
          <Button onClick={this.handleSubmit}>{this.props.submitText}</Button>
          <Button type="default" onClick={this.props.handleCancel}>{this.props.cancelText}</Button>
        </ButtonArea>
      </div>
    );
  }
}

RegisterForm.propTypes = propTypes;
