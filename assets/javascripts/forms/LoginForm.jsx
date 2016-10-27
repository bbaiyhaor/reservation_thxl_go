/**
 * Created by shudi on 2016/10/22.
 */
import React from 'react';
import { CellsTitle, Form, FormCell, CellHeader, Label, CellBody, Input, CellFooter, Icon, ButtonArea, Button } from 'react-weui';
import 'weui';

const propTypes = {
  titleTip: React.PropTypes.string,
  usernameLabel: React.PropTypes.string.isRequired,
  usernamePlaceholder: React.PropTypes.string,
  passwordLabel: React.PropTypes.string.isRequired,
  passwordPlaceholder: React.PropTypes.string,
  submitText: React.PropTypes.string.isRequired,
  cancelText: React.PropTypes.string.isRequired,
  onSubmit: React.PropTypes.func.isRequired,
  onCancel: React.PropTypes.func.isRequired,
};

class LoginForm extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      username: '',
      password: '',
      usernameWarn: false,
      passwordWarn: false,
    };
    this.onSubmit = this.onSubmit.bind(this);
  }

  onSubmit() {
    this.setState({
      usernameWarn: false,
      passwordWarn: false,
    });
    if (this.state.username === '') {
      this.setState({ usernameWarn: true });
      return;
    }
    if (this.state.password === '') {
      this.setState({ passwordWarn: true });
      return;
    }
    this.props.onSubmit(this.state.username, this.state.password);
  }

  handleChange(e, name) {
    this.setState({ [name]: e.target.value });
  }

  render() {
    return (
      <div>
        {
          (this.props.titleTip && this.props.titleTip !== '' &&
            <CellsTitle>{this.props.titleTip}</CellsTitle>)
        }
        <Form>
          <FormCell warn={this.state.usernameWarn}>
            <CellHeader>
              <Label>{this.props.usernameLabel}</Label>
            </CellHeader>
            <CellBody>
              <Input
                type="input"
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
        </Form>
        <ButtonArea direction="horizontal">
          <Button onClick={this.onSubmit}>{this.props.submitText}</Button>
          <Button type="default" onClick={this.props.onCancel}>{this.props.cancelText}</Button>
        </ButtonArea>
      </div>
        );
  }
}

LoginForm.propTypes = propTypes;

export default LoginForm;
