import 'weui';
import React, { PropTypes } from 'react';
import { AlertDialog } from '#coms/Huds';
import { Button } from '#react-weui';
import { User } from '#models/Models';

export default class LogoutButton extends React.Component {
  constructor(props) {
    super(props);
    this.logout = this.logout.bind(this);
  }

  logout() {
    User.logout((data) => {
      if (data.redirect_url) {
        window.location.href = data.redirect_url;
      }
    }, (error) => {
      this.alert.show('登出失败', error, '好的');
    });
  }

  render() {
    const { children, ...others } = this.props;

    return (
      <div>
        <Button {...others} onClick={this.logout}>
          {children}
        </Button>
        <AlertDialog ref={(alert) => { this.alert = alert; }} />
      </div>
    );
  }
}

LogoutButton.propTypes = {
  children: PropTypes.node.isRequired,
};
