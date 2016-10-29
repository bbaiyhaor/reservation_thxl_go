/* eslint max-len: ["off"] */
/**
 * Created by shudi on 2016/10/23.
 */
import React from 'react';
import { hashHistory } from 'react-router';
import { Panel, PanelHeader, PanelBody, CellsTitle, MediaBox, Button, Cells, Cell, CellBody } from 'react-weui';
import 'weui';

import UserLogoutButton from '#coms/LogoutButton';
import PageBottom from '#coms/PageBottom';
import { AlertDialog, ConfirmDialog, LoadingHud } from '#coms/Huds';
import { User, Application } from '#models/Models';

class StudentReservationListPage extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      student: null,
      reservations: null,
    };
    this.showAlert = this.showAlert.bind(this);
  }

  componentDidMount() {
    this.loading.show('正在加载中');
    Application.viewReservationsByStudent(() => {
      setTimeout(() => {
        this.setState({
          student: User.student,
          reservations: Application.reservations,
        }, () => {
          this.loading.hide();
        });
      }, 500);
    }, (status) => {
      this.loading.hide();
      this.alert.show('', status, '好的', () => {
        hashHistory.push('login');
      });
    });
  }

  showAlert(title, msg, label) {
    this.alert.show(title, msg, label);
  }

  render() {
    return (
      <div>
        <Panel access>
          <PanelHeader style={{ fontSize: '18px' }}>
            {User.student && User.student.fullname ? `${User.student.fullname}，` : ''}欢迎使用咨询预约系统
            <div style={{ height: '20px' }}>
              <UserLogoutButton
                size="small"
                style={{ float: 'right' }}
                alert={this.showAlert}
              >
                退出登录
              </UserLogoutButton>
            </div>
            <CellsTitle>请根据您的需要选择相应咨询师和时间段进行预约</CellsTitle>
          </PanelHeader>
          <PanelBody>
            <StudentReservationList reservations={this.state.reservations} />
          </PanelBody>
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

const propTypes = {
  reservations: React.PropTypes.arrayOf(React.PropTypes.object),
};

class StudentReservationList extends React.Component {
  constructor(props) {
    super(props);
    this.makeReservation = this.makeReservation.bind(this);
    this.feedback = this.feedback.bind(this);
  }

  makeReservation(reservation) {
    this.confirm.show('',
      '确定预约后请准确填写个人信息，方便心理咨询中心老师与你取得联系。',
      '暂不预约', '立即预约', null,
      () => hashHistory.push(`/reservation/make?reservation_id=${reservation.id}`)
    );
  }

  feedback(reservation) {
    console.log(reservation);
  }

  renderButton(reservation) {
    if (reservation.status === 1) {
      return (
        <Button
          size="small"
          onClick={(e) => {
            e.preventDefault();
            e.stopPropagation();
            this.makeReservation(reservation);
          }}
        >预约</Button>
      );
    } else if (reservation.status === 2) {
      return <Button size="small" type="default" disabled>已预约</Button>;
    } else if (reservation.status === 3) {
      return (
        <Button
          size="small"
          type="warn"
          onClick={(e) => {
            e.preventDefault();
            e.stopPropagation();
            this.feedback(reservation);
          }}
        >反馈</Button>
      );
    }
    return null;
  }

  render() {
    return (
      <div>
        <MediaBox type="small_appmsg">
          <Cells access>
            {this.props.reservations && this.props.reservations.map(reservation =>
              <Cell key={`reservation-cell-${reservation.id}`}>
                <CellBody>
                  <p style={{ fontSize: '14px' }}>
                    {reservation.start_time} - {reservation.end_time.slice(-5)}　{reservation.teacher_fullname}
                  </p>
                </CellBody>
                {this.renderButton(reservation)}
              </Cell>)
            }
          </Cells>
        </MediaBox>
        <ConfirmDialog ref={(confirm) => { this.confirm = confirm; }} />
      </div>
    );
  }
}

StudentReservationList.propTypes = propTypes;

export default StudentReservationListPage;
