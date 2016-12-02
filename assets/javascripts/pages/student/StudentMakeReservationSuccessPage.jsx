/**
 * Created by shudi on 2016/11/28.
 */
import React from 'react';
import { hashHistory } from 'react-router';
import { Msg } from '#react-weui';
import 'weui';

export default class StudentMakeReservationSuccessPage extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      buttons: [{
        type: 'primary',
        label: '返回首页',
        onClick: () => {
          hashHistory.push('reservation');
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
