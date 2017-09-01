/* eslint max-len: ["off"] */
import 'weui';
// import 'react-weui/lib/react-weui.min.css';
import { AlertDialog, ConfirmDialog, LoadingHud } from '#coms/Huds';
import { Application, User } from '#models/Models';
import { Button, Cell, CellBody, Cells, CellsTitle, MediaBox, Panel, PanelBody, PanelHeader, SearchBar, Toptips } from 'react-weui';
import LogoutButton from '#coms/LogoutButton';
import MobileDetect from '#utils/MobileDetect';
import PageBottom from '#coms/PageBottom';
import PropTypes from 'prop-types';
import React from 'react';

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
    User.updateSession(() => {
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
          this.props.history.push('/login');
        });
      });
    }, () => {
      this.props.history.push('/login');
    });
  }

  render() {
    return (
      <div>
        {
          MobileDetect.isWechat() ?
            <div style={{ marginBottom: '25px' }}>
              <Toptips type="info" show>
                如遇预约错误，请使用系统自带浏览器。
              </Toptips>
            </div> : null
        }
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
            </div>
            <CellsTitle>可预约时间为当前时间加一周减去1.5小时，请根据您的需要选择相应咨询师和时间段进行预约。</CellsTitle>
          </PanelHeader>
          <PanelBody>
            <StudentReservationList
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

StudentReservationListPage.propTypes = {
  history: PropTypes.object.isRequired,
};

class StudentReservationList extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      reservations: this.props.reservations,
      reservationsBak: this.props.reservations,
    };
    this.toFeedback = this.toFeedback.bind(this);
    this.handleChange = this.handleChange.bind(this);
    this.makeReservation = this.makeReservation.bind(this);
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
    const { history } = this.props;
    this.confirm.show('',
      '确定预约后请准确填写个人信息，方便心理咨询中心老师与你取得联系。',
      '暂不预约', '立即预约', null,
      () => history.push('/reservation/make', { reservation_id: `${reservation.id}` }),
    );
  }

  renderButton(reservation) {
    if (reservation.status === 1) {
      return (
        <Button
          size="small"
          onClick={(e) => {
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
            e.stopPropagation();
            this.toFeedback(reservation);
          }}
        >反馈</Button>
      );
    }
    return null;
  }

  render() {
    const normalStyle = {
      fontSize: '14px',
    };
    const hightlightStyle = {
      color: '#EF4F4F',
      ...normalStyle,
    };
    return (
      <div>
        <SearchBar
          onChange={this.handleChange}
        />
        <MediaBox type="small_appmsg">
          <Cells>
            {this.state.reservations && this.state.reservations.map(reservation =>
              (<Cell key={`reservation-cell-${reservation.id}`}>
                <CellBody>
                  {reservation.student_id && reservation.student_id === User.userId ?
                    <p style={{ ...hightlightStyle }}>
                      {reservation.start_time}-{reservation.end_time.slice(-5)}({reservation.start_weekday}) {reservation.teacher_fullname}
                    </p> : <p style={{ ...normalStyle }}>
                      {reservation.start_time}-{reservation.end_time.slice(-5)}({reservation.start_weekday}) {reservation.teacher_fullname}
                    </p>
                  }
                </CellBody>
                {this.renderButton(reservation)}
              </Cell>))
            }
          </Cells>
        </MediaBox>
        <ConfirmDialog ref={(confirm) => { this.confirm = confirm; }} />
      </div>
    );
  }
}

StudentReservationList.propTypes = {
  history: PropTypes.object.isRequired,
  reservations: PropTypes.arrayOf(PropTypes.object),
};

StudentReservationList.defaultProps = {
  reservations: [],
};
