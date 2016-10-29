/**
 * Created by shudi on 2016/10/22.
 */
import React from 'react';
import { Button } from '#react-weui';
import 'weui';

import { User } from '#models/Models';

const propTypes = {
  children: React.PropTypes.node,
  alert: React.PropTypes.func,
};

class UserLogoutButton extends React.Component {
  constructor(props) {
    super(props);
    this.logout = this.logout.bind(this);
  }

  logout() {
    User.logout((payload) => {
      if (payload.redirect_url) {
        window.location.href = payload.redirect_url;
      }
    }, (status) => {
      this.props.alert('登出失败', status, '好的');
    });
  }

  render() {
    const { children, ...others } = this.props;

    return (
      <Button {...others} onClick={this.logout}>
        {children}
      </Button>
    );
  }
}

UserLogoutButton.propTypes = propTypes;

export default UserLogoutButton;
