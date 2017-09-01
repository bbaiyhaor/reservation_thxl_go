import 'weui';
import { Button, CellBody, CellFooter, CellHeader, FormCell, Icon, Input, Label } from 'react-weui';
import PropTypes from 'prop-types';
import React from 'react';

export default class VerifyCodeInput extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      verifyCode: '',
      verifyCodeWarn: false,
      firstTry: true,
      active: true,
      counter: 60,
    };
    this.tick = this.tick.bind(this);
    this.handleChange = this.handleChange.bind(this);
    this.handleClick = this.handleClick.bind(this);
  }

  componentWillUnmount() {
    clearInterval(this.interval);
  }

  getValue() {
    return this.state.verifyCode;
  }

  setInvalidate() {
    this.setState({
      verifyCodeWarn: true,
    });
  }

  resetWarn() {
    this.setState({
      verifyCodeWarn: false,
    });
  }

  focus() {
    this.verifyCodeInput.focus();
  }

  restart() {
    this.setState({
      active: false,
      counter: 60,
      firstTry: false,
    }, () => {
      this.interval = setInterval(this.tick, 1000);
    });
  }

  tick() {
    const counter = this.state.counter;
    this.setState({
      counter: counter - 1,
      active: counter <= 1,
    }, () => {
      if (this.state.counter <= 0) {
        clearInterval(this.interval);
      }
    });
  }

  handleChange(e, name) {
    this.setState({ [name]: e.target.value });
  }

  handleClick() {
    if (!this.state.active) {
      return;
    }
    this.props.handleClick && this.props.handleClick();
  }


  render() {
    let sendText = this.state.firstTry ? '发送' : '重新发送';
    if (!this.state.active) {
      sendText = `重发(${this.state.counter}s)`;
    }
    return (
      <FormCell warn={this.state.verifyCodeWarn}>
        <CellHeader>
          <Label>{this.props.verifyCodeLabel}</Label>
        </CellHeader>
        <CellBody>
          <Input
            ref={(verifyCodeInput) => { this.verifyCodeInput = verifyCodeInput; }}
            type="input"
            placeholder={this.props.verifyCodePlaceholder}
            value={this.state.verifyCode}
            onChange={(e) => { this.handleChange(e, 'verifyCode'); }}
          />
        </CellBody>
        <CellFooter>
          {this.state.verifyCodeWarn &&
          <Icon value="warn" style={{ marginRight: '5px' }} />
          }
          <Button
            size="small"
            type={this.state.active ? 'primary' : 'default'}
            style={{ padding: '0px 10px' }}
            onClick={this.handleClick}
            disabled={!this.state.active}
          >
            {sendText}
          </Button>
        </CellFooter>
      </FormCell>
    );
  }
}

VerifyCodeInput.propTypes = {
  verifyCodeLabel: PropTypes.string.isRequired,
  verifyCodePlaceholder: PropTypes.string.isRequired,
  handleClick: PropTypes.func.isRequired,
};
