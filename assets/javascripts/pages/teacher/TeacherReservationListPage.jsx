/* eslint max-len: ["off"] */
import 'weui';
import { AlertDialog, LoadingHud } from '#coms/Huds';
import { Application, User } from '#models/Models';
import { Button, CellBody, CellFooter, CellsTitle, FormCell, Icon, Input, MediaBox, MediaBoxBody, MediaBoxDescription, MediaBoxTitle, Panel, PanelBody, PanelHeader, SearchBar } from '#react-weui';
import React, { PropTypes } from 'react';
import { Link } from 'react-router-dom';
import LogoutButton from '#coms/LogoutButton';
import PageBottom from '#coms/PageBottom';

export default class TeacherReservationListPage extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      teacher: null,
      reservations: null,
      studentUsername: '',
      studentUsernameWarn: false,
    };
    this.toChangePassword = this.toChangePassword.bind(this);
    this.handleChange = this.handleChange.bind(this);
    this.toStudentInfo = this.toStudentInfo.bind(this);
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
          this.props.history.push('/login');
        });
      });
    }, () => {
      this.props.history.push('/login');
    });
  }

  toChangePassword() {
    this.props.history.push('/password/change');
  }

  handleChange(e, name) {
    const value = e.target.value;
    this.setState({ [name]: value });
  }

  toStudentInfo() {
    this.setState({ studentUsernameWarn: false });
    if (this.state.studentUsername === '') {
      this.setState({ studentUsernameWarn: true });
      this.studentUsernameInput.focus();
      return;
    }
    Application.queryStudentInfoByTeacher(this.state.studentUsername, () => {
      this.props.history.push('/student', { student_username: `${this.state.studentUsername}` });
    }, (error) => {
      this.alert.show('查询失败', error, '好的', () => {
        this.alert.hide();
        this.studentUsernameInput.focus();
      });
    });
  }

  render() {
    return (
      <div>
        <Panel>
          <PanelHeader style={{ fontSize: '18px' }}>
            {User.fullname !== '' ? `${User.fullname}，` : ''}欢迎使用咨询预约系统
            <div style={{ height: '20px' }}>
              <LogoutButton
                size="small"
                style={{ float: 'right' }}
              >
                退出登录
              </LogoutButton>
              <Button
                size="small"
                style={{ float: 'right', marginRight: '15px' }}
                onClick={this.toChangePassword}
              >
                更改密码
              </Button>
            </div>
            <CellsTitle>点击预约学生姓名可查看学生信息，红色咨询为危机个案</CellsTitle>
            <CellsTitle>输入学生学号查询学生信息</CellsTitle>
            <FormCell warn style={{ padding: '10px 0px 0px 15px' }}>
              <CellBody>
                <Input
                  ref={(studentUsernameInput) => { this.studentUsernameInput = studentUsernameInput; }}
                  type="input"
                  placeholder="请输入学生学号"
                  value={this.state.studentUsername}
                  onChange={(e) => { this.handleChange(e, 'studentUsername'); }}
                />
              </CellBody>
              <CellFooter>
                {this.state.studentUsernameWarn &&
                  <Icon value="warn" style={{ marginRight: '5px' }} />
                }
                <Button size="small" onClick={this.toStudentInfo}>查询</Button>
              </CellFooter>
            </FormCell>
          </PanelHeader>
          <PanelBody>
            <TeacherReservationList
              reservations={this.state.reservations}
              history={this.props.history}
            />
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

TeacherReservationListPage.propTypes = {
  history: PropTypes.object.isRequired,
};

class TeacherReservationList extends React.Component {
  static renderStudentSpan(reservation) {
    const style = {
      marginLeft: '20px',
    };
    if (reservation.status === 2 || reservation.status === 3) {
      if (reservation.student_crisis_level > 0) {
        return (
          <Link
            to={{
              pathname: '/student',
              state: { student_id: `${reservation.student_id}` },
            }}
            style={{ color: '#EF4F4F', ...style }}
          >
            学生：{reservation.student_fullname}
          </Link>
        );
      }
      return (
        <Link
          to={{
            pathname: '/student',
            state: { student_id: `${reservation.student_id}` },
          }}
          style={{ color: '#999999', ...style }}
        >
          学生：{reservation.student_fullname}
        </Link>
      );
    }
    return null;
  }

  constructor(props) {
    super(props);
    this.state = {
      reservations: this.props.reservations,
      reservationsBak: this.props.reservations,
    };
    this.toFeedback = this.toFeedback.bind(this);
    this.handleChange = this.handleChange.bind(this);
  }

  componentWillReceiveProps(nextProps) {
    nextProps.reservations && this.setState({
      reservations: nextProps.reservations,
      reservationsBak: nextProps.reservations,
    });
  }

  toFeedback(reservation) {
    this.props.history.push('/reservation/feedback', { reservation_id: `${reservation.id}` });
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

  renderStatusButton(reservation) {
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
      const type = reservation.has_teacher_feedback ? 'default' : 'primary';
      return (
        <Button
          style={{ ...style }}
          size="small"
          type={type}
          onClick={(e) => {
            e.stopPropagation();
            this.toFeedback(reservation);
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
        {this.state.reservations && this.state.reservations.map(reservation =>
          <MediaBox
            key={`reservation-box-${reservation.id}`}
            type="appmsg"
            style={{ padding: '10px 15px' }}
          >
            <MediaBoxBody>
              <MediaBoxTitle style={{ marginBottom: '5px' }}>
                {reservation.start_time} - {reservation.end_time.slice(-5)}
                {this.renderStatusButton(reservation)}
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

TeacherReservationList.propTypes = {
  history: PropTypes.object.isRequired,
  reservations: PropTypes.arrayOf(React.PropTypes.object),
};

TeacherReservationList.defaultProps = {
  reservations: [],
};
