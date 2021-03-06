import 'weui';
import { Dialog, Toast } from 'react-weui';
import React from 'react';

export class AlertDialog extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      alertShow: false,
      alert: {
        title: '',
        msg: '',
        buttons: [
          {
            label: '',
            onClick: this.hide.bind(this),
          },
        ],
      },
    };
  }

  show(title, msg, label, click) {
    this.setState({
      alertShow: true,
      alert: {
        title,
        msg,
        buttons: [
          {
            label,
            onClick: click || this.hide.bind(this),
          },
        ],
      },
    });
  }

  hide() {
    this.setState({
      alertShow: false,
      alert: {
        title: '',
        msg: '',
      },
    });
  }

  render() {
    return (
      <Dialog
        title={this.state.alert.title}
        buttons={this.state.alert.buttons}
        show={this.state.alertShow}
      >
        {this.state.alert.msg}
      </Dialog>
    );
  }
}

export class ConfirmDialog extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      confirmShow: false,
      confirm: {
        title: '',
        msg: '',
        buttons: [
          {
            type: 'default',
            label: '',
          },
          {
            type: 'primary',
            label: '',
          },
        ],
      },
    };
    this.show = this.show.bind(this);
    this.hide = this.hide.bind(this);
  }

  show(title, msg, label1, label2, click1, click2) {
    this.setState({
      confirmShow: true,
      confirm: {
        title,
        msg,
        buttons: [
          {
            type: 'default',
            label: label1,
            onClick: click1 || this.hide.bind(this),
          },
          {
            type: 'primary',
            label: label2,
            onClick: click2,
          },
        ],
      },
    });
  }

  hide() {
    this.setState({
      confirmShow: false,
      confirm: {
        title: '',
        msg: '',
        buttons: [
          {
            type: 'default',
            label: '',
          },
          {
            type: 'primary',
            label: '',
          },
        ],
      },
    });
  }

  render() {
    return (
      <Dialog
        title={this.state.confirm.title}
        buttons={this.state.confirm.buttons}
        show={this.state.confirmShow}
      >
        {this.state.confirm.msg}
      </Dialog>
    );
  }
}

export class LoadingHud extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      showLoading: false,
      loadingMsg: '',
    };
    this.show = this.show.bind(this);
    this.hide = this.hide.bind(this);
  }

  show(msg, duration) {
    this.setState({
      showLoading: true,
      loadingMsg: msg,
    });
    if (duration && duration > 0) {
      setTimeout(() => {
        this.hide();
      }, duration);
    }
  }

  hide() {
    this.setState({
      showLoading: false,
      loadingMsg: '',
    });
  }

  render() {
    return (
      <Toast icon="loading" show={this.state.showLoading}>
        {this.state.loadingMsg}
      </Toast>
    );
  }
}
