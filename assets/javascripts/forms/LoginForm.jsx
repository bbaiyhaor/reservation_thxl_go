import 'weui';
import { Button, ButtonArea, CellBody, CellFooter, CellHeader, CellsTitle, Form, FormCell, Icon, Input, Label } from 'react-weui';
import React, { PropTypes } from 'react';

export default class LoginForm extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      username: '',
      password: '',
      usernameWarn: false,
      passwordWarn: false,
    };
    this.handleChange = this.handleChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  handleChange(e, name) {
    this.setState({ [name]: e.target.value });
  }

  handleSubmit() {
    this.setState({
      usernameWarn: false,
      passwordWarn: false,
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
    this.props.handleSubmit(this.state.username, this.state.password);
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
        </Form>
        <ButtonArea direction="horizontal">
          <Button onClick={this.handleSubmit}>{this.props.submitText}</Button>
          {this.props.cancelText && <Button type="default" onClick={this.props.handleCancel}>{this.props.cancelText}</Button>}
        </ButtonArea>
      </div>
    );
  }
}

LoginForm.propTypes = {
  titleTip: PropTypes.string,
  usernameLabel: PropTypes.string.isRequired,
  usernamePlaceholder: PropTypes.string,
  passwordLabel: PropTypes.string.isRequired,
  passwordPlaceholder: PropTypes.string,
  submitText: PropTypes.string.isRequired,
  cancelText: PropTypes.string,
  handleSubmit: PropTypes.func.isRequired,
  handleCancel: PropTypes.func,
};

LoginForm.defaultProps = {
  titleTip: '',
  usernamePlaceholder: '',
  passwordPlaceholder: '',
  cancelText: '',
  handleCancel: undefined,
};
