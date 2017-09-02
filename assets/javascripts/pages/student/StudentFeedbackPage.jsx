import 'weui';
import { AlertDialog, LoadingHud } from '#coms/Huds';
import { Application, User } from '#models/Models';
import { Panel, PanelHeader } from 'react-weui';
import PageBottom from '#coms/PageBottom';
import PropTypes from 'prop-types';
import React from 'react';
import StudentFeedbackForm from '#forms/StudentFeedbackForm';
import queryString from 'query-string';

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

  componentDidMount() {
    this.loading.show('正在加载中');
    const parsedQuery = queryString.parse(this.props.history.location.search);
    const reservationId = parsedQuery.reservation_id;
    User.updateSession(() => {
      Application.getFeedbackByStudent(reservationId, (data) => {
        setTimeout(() => {
          this.setState({
            scores: data.feedback.scores,
            reservation: data.reservation,
          }, () => {
            this.loading.hide();
          });
        }, 500);
      }, (error) => {
        this.loading.hide();
        this.alert.show('', error, '好的', () => {
          this.props.history.push('/reservation');
        });
      });
    }, () => {
      this.props.history.push('/login');
    });
  }

  handleCancel() {
    this.props.history.goBack();
  }

  handleSubmit(feedback1, feedback2, feedback3, feedback4, feedback5) {
    const scores = [feedback1, feedback2, feedback3, feedback4, feedback5];
    this.loading.show('正在加载中');
    setTimeout(() => {
      Application.submitFeedbackByStudent(this.state.reservation.id, scores, () => {
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
