/* eslint max-len: ["off"] */
/**
 * Created by shudi on 2016/10/23.
 */
import React, { PropTypes } from 'react';
import { hashHistory } from 'react-router';
import { Panel, PanelHeader, PanelBody, CellsTitle, MediaBox, Button, Cells, Cell, CellBody, SearchBar } from '#react-weui';
import 'weui';

import LogoutButton from '#coms/LogoutButton';
import PageBottom from '#coms/PageBottom';
import { AlertDialog, ConfirmDialog, LoadingHud } from '#coms/Huds';
import { User, Application } from '#models/Models';

export default class StudentReservationListPage extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      student: null,
      reservations: null,
    };
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
    }, (error) => {
      this.loading.hide();
      this.alert.show('', error, '好的', () => {
        hashHistory.push('login');
      });
    });
  }

  render() {
    return (
      <div>
        <Panel access>
          <PanelHeader style={{ fontSize: '18px' }}>
            {User.student && User.student.fullname ? `${User.student.fullname}，` : ''}欢迎使用咨询预约系统
            <div style={{ height: '20px' }}>
              <LogoutButton
                size="small"
                style={{ float: 'right' }}
              >
                退出登录
              </LogoutButton>
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
  reservations: PropTypes.arrayOf(React.PropTypes.object),
};

class StudentReservationList extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      reservations: this.props.reservations,
      reservationsBak: this.props.reservations,
    };
    this.handleChange = this.handleChange.bind(this);
    this.makeReservation = this.makeReservation.bind(this);
    this.feedback = this.feedback.bind(this);
  }

  componentWillReceiveProps(nextProps) {
    nextProps.reservations && this.setState({
      reservations: nextProps.reservations,
      reservationsBak: nextProps.reservations,
    });
  }

  handleChange(text) {
    const keyword = [text];
    if (keyword === '') {
      this.setState(prevState => ({
        reservations: prevState.reservationsBak,
      }));
    }
    const result = this.state.reservationsBak.filter((reservation) => {
      if (reservation.teacher_fullname.indexOf(keyword) !== -1) {
        return true;
      } else if (reservation.start_time.indexOf(keyword) !== -1) {
        return true;
      } else if (reservation.end_time.indexOf(keyword) !== -1) {
        return true;
      }
      return false;
    });
    this.setState({ reservations: result });
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
        <SearchBar
          onChange={this.handleChange}
        />
        <MediaBox type="small_appmsg">
          <Cells access>
            {this.state.reservations && this.state.reservations.map(reservation =>
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
