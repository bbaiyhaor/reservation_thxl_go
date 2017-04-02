import 'weui';
import { AlertDialog, LoadingHud } from '#coms/Huds';
import { Application, User } from '#models/Models';
import { Panel, PanelHeader } from 'react-weui';
import React, { PropTypes } from 'react';
import PageBottom from '#coms/PageBottom';
import StudentFeedbackForm from '#forms/StudentFeedbackForm';

export default class StudentFeedbackPage extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      reservation: null,
      scores: [],
    };
    this.handleCancel = this.handleCancel.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  componentWillMount() {
    const reservationId = this.props.history.location.state.reservation_id || '';
    if (reservationId === '' || !User.student || !Application.reservations) {
      this.props.history.push('/reservation');
      return;
    }
    let i = 0;
    for (; i < Application.reservations.length; i += 1) {
      if (Application.reservations[i].id === reservationId) {
        this.setState({
          reservation: Application.reservations[i],
        });
        break;
      }
    }
    if (i === Application.reservations.length) {
      this.props.history.push('/reservation');
    }
  }

  componentDidMount() {
    this.loading.show('正在加载中');
    setTimeout(() => {
      Application.getFeedbackByStudent(this.state.reservation.id, this.state.reservation.source_id, (data) => {
        this.loading.hide();
        this.setState({
          scores: data.feedback.scores,
        });
      }, (error) => {
        this.loading.hide();
        this.alert.show('', error, '好的', () => {
          this.props.history.push('/reservation');
        });
      });
    }, 500);
  }

  handleCancel() {
    this.props.history.goBack();
  }

  handleSubmit(feedback1, feedback2, feedback3, feedback4, feedback5) {
    const scores = [feedback1, feedback2, feedback3, feedback4, feedback5];
    this.loading.show('正在加载中');
    setTimeout(() => {
      Application.submitFeedbackByStudent(this.state.reservation.id, this.state.reservation.source_id, scores, () => {
        this.loading.hide();
        this.alert.show('提交成功', '你已成功提交反馈', '好的', () => {
          this.props.history.push('/reservation');
        });
      }, (error) => {
        this.loading.hide();
        this.alert.show('提交失败', error, '好的');
      });
    }, 500);
  }

  render() {
    return (
      <div>
        <Panel>
          <PanelHeader style={{ fontSize: '18px' }}>学生咨询反馈表</PanelHeader>
          <StudentFeedbackForm
            reservation={this.state.reservation}
            scores={this.state.scores}
            handleSubmit={this.handleSubmit}
            handleCancel={this.handleCancel}
          />
          <LoadingHud ref={(loading) => { this.loading = loading; }} />
          <AlertDialog ref={(alert) => { this.alert = alert; }} />
          <PageBottom
            styles={{ color: '#999999', textAlign: 'center', backgroundColor: 'white', fontSize: '14px' }}
            contents={['清华大学学生心理发展指导中心', '联系方式：010-62782007']}
            height="55px"
          />
        </Panel>
      </div>
    );
  }
}

StudentFeedbackPage.propTypes = {
  history: PropTypes.object.isRequired,
};
