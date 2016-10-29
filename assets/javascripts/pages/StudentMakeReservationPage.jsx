/**
 * Created by shudi on 2016/10/24.
 */
import React, { PropTypes } from 'react';
import { hashHistory } from 'react-router';
import { Panel, PanelHeader } from '#react-weui';
import 'weui';

import MakeReservationForm from '#forms/MakeReservationForm';
import PageBottom from '#coms/PageBottom';
import { AlertDialog, LoadingHud } from '#coms/Huds';
import { User, Application } from '#models/Models';

const propTypes = {
  location: PropTypes.object,
};

export default class StudentMakeReservationPage extends React.Component {
  static handleCancel() {
    hashHistory.goBack();
  }

  constructor(props) {
    super(props);
    this.state = {
      student: null,
      reservation: null,
    };
  }

  componentWillMount() {
    const reservationId = this.props.location.query.reservation_id;
    if (reservationId === '' || !User.student || !Application.reservations || Application.reservations.length === 0) {
      hashHistory.push('reservation');
      return;
    }
    let i = 0;
    for (; i < Application.reservations.length; i += 1) {
      if (Application.reservations[i].id === reservationId) {
        this.setState({
          student: User.student,
          reservation: Application.reservations[i],
        });
        break;
      }
    }
    if (i === Application.reservations.length) {
      hashHistory.push('reservation');
      return;
    }
  }

  render() {
    return (
      <div>
        <Panel access>
          <PanelHeader style={{ fontSize: '18px' }}>学生信息登记表</PanelHeader>
          <MakeReservationForm
            student={this.state.student}
            reservation={this.state.reservation}
            handleCancel={StudentMakeReservationPage.handleCancel}
          />
        </Panel>
        <LoadingHud ref={(loading) => { this.loading = loading; }} />
        <AlertDialog ref={(alert) => { this.alert = alert; }} />
        <PageBottom
          styles={{ color: '#999999', textAlign: 'center', backgroundColor: 'white', fontSize: '14px' }}
          contents={['清华大学学生心理发展指导中心', '联系方式：010-62782007']}
          height="55px"
        />
      </div>
    );
  }
}

StudentMakeReservationPage.propTypes = propTypes;
