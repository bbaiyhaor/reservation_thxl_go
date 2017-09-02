import 'weui';
import { Msg } from 'react-weui';
import PropTypes from 'prop-types';
import React from 'react';

export default class StudentMakeReservationSuccessPage extends React.Component {
  constructor(props) {
    super(props);
    const { history } = this.props;
    this.state = {
      buttons: [{
        type: 'primary',
        label: '返回首页',
        onClick: () => {
          history.push('/reservation');
        },
      }],
    };
  }

  render() {
    return (
      <Msg
        type="success"
        title="预约成功"
        description="你已预约成功，请关注短信提醒。"
        buttons={this.state.buttons}
      />
    );
  }
}

StudentMakeReservationSuccessPage.propTypes = {
  history: PropTypes.object.isRequired,
};
