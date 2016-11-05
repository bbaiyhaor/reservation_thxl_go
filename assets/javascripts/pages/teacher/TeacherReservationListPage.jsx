/* eslint max-len: ["off"] */
/**
 * Created by shudi on 2016/10/23.
 */
import React, { PropTypes } from 'react';
import { Link, hashHistory } from 'react-router';
import { Panel, PanelHeader, PanelBody, CellsTitle, MediaBox, MediaBoxTitle, MediaBoxBody, MediaBoxDescription, Button, SearchBar } from '#react-weui';
import 'weui';

import LogoutButton from '#coms/LogoutButton';
import PageBottom from '#coms/PageBottom';
import { AlertDialog, LoadingHud } from '#coms/Huds';
import { User, Application } from '#models/Models';

export default class TeacherReservationListPage extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      teacher: null,
      reservations: null,
    };
  }

  componentDidMount() {
    this.loading.show('正在加载中');
    User.updateSession(() => {
      Application.viewReservationsByTeacher(() => {
        setTimeout(() => {
          this.setState({
            teacher: User.teacher,
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
    }, () => {
      hashHistory.push('login');
    });
  }

  render() {
    return (
      <div>
        <Panel access>
          <PanelHeader style={{ fontSize: '18px' }}>
            {User.fullname !== '' ? `${User.fullname}，` : ''}欢迎使用咨询预约系统
            <div style={{ height: '20px' }}>
              <LogoutButton
                size="small"
                style={{ float: 'right' }}
              >
                退出登录
              </LogoutButton>
            </div>
            <CellsTitle>点击预约学生姓名可查看学生信息，红色咨询为危机个案</CellsTitle>
          </PanelHeader>
          <PanelBody>
            <TeacherReservationList reservations={this.state.reservations} />
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

class TeacherReservationList extends React.Component {
  static feedback(reservation) {
    hashHistory.push(`reservation/feedback?reservation_id=${reservation.id}`);
  }

  static renderStatusButton(reservation) {
    const style = {
      float: 'right',
      marginRight: '10px',
      zIndex: 10,
    };
    if (reservation.status === 1) {
      return (
        <Button
          style={{ ...style }}
          size="small"
          type="default"
          disabled
        >未预约</Button>
      );
    } else if (reservation.status === 2) {
      return (
        <Button
          style={{ ...style }}
          size="small"
          type="default"
          disabled
        >已预约</Button>
      );
    } else if (reservation.status === 3) {
      const type = reservation.student_crisis_level > 0 ? 'warn' : 'primary';
      return (
        <Button
          style={{ ...style }}
          size="small"
          type={type}
          onClick={(e) => {
            e.stopPropagation();
            TeacherReservationList.feedback(reservation);
          }}
        >反馈</Button>
      );
    }
    return null;
  }

  static renderStudentSpan(reservation) {
    const style = {
      marginLeft: '20px',
    };
    if (reservation.status === 2 || reservation.status === 3) {
      if (reservation.student_crisis_level > 0) {
        return <Link to={`student?student_id=${reservation.student_id}`} style={{ color: '#EF4F4F', ...style }}>学生：{reservation.student_fullname}</Link>;
      }
      return <Link to={`student?student_id=${reservation.student_id}`} style={{ color: '#999999', ...style }}>学生：{reservation.student_fullname}</Link>;
    }
    return null;
  }

  constructor(props) {
    super(props);
    this.state = {
      reservations: this.props.reservations,
      reservationsBak: this.props.reservations,
    };
    this.handleChange = this.handleChange.bind(this);
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
      } else if (reservation.teacher_mobile.indexOf(keyword) !== -1) {
        return true;
      } else if (reservation.teacher_username.indexOf(keyword) !== -1) {
        return true;
      } else if (reservation.start_time.indexOf(keyword) !== -1) {
        return true;
      } else if (reservation.end_time.indexOf(keyword) !== -1) {
        return true;
      } else if (reservation.student_fullname && reservation.student_fullname.indexOf(keyword) !== -1) {
        return true;
      } else if (reservation.student_username && reservation.student_username.indexOf(keyword) !== -1) {
        return true;
      }
      return false;
    });
    this.setState({ reservations: result });
  }

  render() {
    return (
      <div>
        <SearchBar
          onChange={this.handleChange}
        />
        {this.state.reservations && this.state.reservations.map(reservation =>
          <MediaBox
            key={`reservation-box-${reservation.id}`}
            type="appmsg"
            style={{ padding: '10px 15px' }}
          >
            <MediaBoxBody>
              <MediaBoxTitle style={{ marginBottom: '5px' }}>
                {reservation.start_time} - {reservation.end_time.slice(-5)}
                {TeacherReservationList.renderStatusButton(reservation)}
              </MediaBoxTitle>
              {reservation.teacher_fullname &&
              <MediaBoxDescription>
                {reservation.teacher_fullname}&nbsp;&nbsp;&nbsp;{reservation.teacher_mobile}
                {TeacherReservationList.renderStudentSpan(reservation)}
              </MediaBoxDescription>
              }
            </MediaBoxBody>
          </MediaBox>)
        }
      </div>
    );
  }
}

TeacherReservationList.propTypes = propTypes;
