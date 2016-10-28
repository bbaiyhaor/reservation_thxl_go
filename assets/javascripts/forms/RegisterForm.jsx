/**
 * Created by shudi on 2016/10/22.
 */
import React from 'react';
import { Link } from 'react-router';
import { CellsTitle, Form, FormCell, CellHeader, Label, CellBody, Input, CellFooter, Icon, ButtonArea, Button, Checkbox } from 'react-weui';
import 'weui';

const propTypes = {
  titleTip: React.PropTypes.string,
  usernameLabel: React.PropTypes.string.isRequired,
  usernamePlaceholder: React.PropTypes.string,
  passwordLabel: React.PropTypes.string.isRequired,
  passwordPlaceholder: React.PropTypes.string,
  confirmPasswordLabel: React.PropTypes.string.isRequired,
  confirmPasswordPlaceholder: React.PropTypes.string,
  protocol: React.PropTypes.string,
  protocolPrefix: React.PropTypes.string,
  protocolSuffix: React.PropTypes.string,
  protocolLink: React.PropTypes.string,
  submitText: React.PropTypes.string.isRequired,
  cancelText: React.PropTypes.string.isRequired,
  handleSubmit: React.PropTypes.func.isRequired,
  handleCancel: React.PropTypes.func.isRequired,
  showAlert: React.PropTypes.func,
};

class RegisterForm extends React.Component {
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
      return;
    }
    if (this.state.password === '') {
      this.setState({ passwordWarn: true });
      return;
    }
    if (this.state.confirmPassword === '') {
      this.setState({ confirmPasswordWarn: true });
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
        {
          this.props.titleTip && this.props.titleTip !== '' &&
            <CellsTitle>{this.props.titleTip}</CellsTitle>
        }
        <Form checkbox={this.props.protocol && true} className="weui_cells_form">
          <FormCell warn={this.state.usernameWarn}>
            <CellHeader>
              <Label>{this.props.usernameLabel}</Label>
            </CellHeader>
            <CellBody>
              <Input
                type="tel"
                placeholder={this.props.usernamePlaceholder}
                value={this.state.username}
                onChange={(e) => { this.handleChange(e, 'username'); }}
              />
            </CellBody>
            {
              this.state.usernameWarn &&
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
                type="password"
                placeholder={this.props.passwordPlaceholder}
                value={this.state.password}
                onChange={(e) => { this.handleChange(e, 'password'); }}
              />
            </CellBody>
            {
              this.state.passwordWarn &&
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
                type="password"
                placeholder={this.props.confirmPasswordPlaceholder}
                value={this.state.confirmPassword}
                onChange={(e) => { this.handleChange(e, 'confirmPassword'); }}
              />
            </CellBody>
            {
              this.state.confirmPasswordWarn &&
                <CellFooter>
                  <Icon value="warn" />
                </CellFooter>
            }
          </FormCell>
          {
            this.props.protocol ?
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
                {
                  this.state.protocolWarn &&
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

export default RegisterForm;
