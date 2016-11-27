/**
 * Created by shudi on 2016/11/4.
 */
import React, { PropTypes } from 'react';
import { hashHistory } from 'react-router';
import { Panel, PanelHeader } from '#react-weui';
import 'weui';

import TeacherFeedbackForm from '#forms/TeacherFeedbackForm';
import PageBottom from '#coms/PageBottom';
import { AlertDialog, LoadingHud } from '#coms/Huds';
import { User, Application } from '#models/Models';

const propTypes = {
  location: PropTypes.object,
};

export default class TeacherFeedbackPage extends React.Component {
  static handleCancel() {
    hashHistory.goBack();
  }

  constructor(props) {
    super(props);
    this.state = {
      reservation: null,
      feedback: null,
      crisisLevel: 0,
      firstCategories: null,
      secondCategories: null,
    };
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  componentDidMount() {
    const reservationId = this.props.location.query.reservation_id;
    if (reservationId === '' || !User.teacher || !Application.reservations) {
      hashHistory.push('reservation');
      return;
    }
    let i = 0;
    let reservation = null;
    for (; i < Application.reservations.length; i += 1) {
      if (Application.reservations[i].id === reservationId) {
        reservation = Application.reservations[i];
        break;
      }
    }
    if (i === Application.reservations.length) {
      hashHistory.push('reservation');
      return;
    }
    this.loading.show('正在加载中');
    Application.getFeedbackCategories((data) => {
      this.setState({
        firstCategories: data.first_category,
        secondCategories: data.second_category,
      });
    }, (error) => {
      this.loading.hide();
      this.alert.show('', error, '好的', () => {
        hashHistory.push('reservation');
      });
    });
    setTimeout(() => {
      Application.getFeedbackByTeacher(reservation.id, reservation.source_id, (data) => {
        this.loading.hide();
        this.setState({
          reservation,
          feedback: data.feedback,
        });
      }, (error) => {
        this.loading.hide();
        this.alert.show('', error, '好的', () => {
          hashHistory.push('reservation');
        });
      });
    }, 500);
  }

  handleSubmit(payload) {
    this.loading.show('正在加载中');
    setTimeout(() => {
      Application.submitFeedbackByTeacher(payload, () => {
        this.loading.hide();
        this.alert.show('提交成功', '您已成功提交反馈', '好的', () => {
          hashHistory.push('reservation');
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
          <PanelHeader style={{ fontSize: '18px' }}>咨询师反馈表</PanelHeader>
          <TeacherFeedbackForm
            reservation={this.state.reservation}
            feedback={this.state.feedback}
            firstCategories={this.state.firstCategories}
            secondCategories={this.state.secondCategories}
            handleSubmit={this.handleSubmit}
            handleCancel={TeacherFeedbackPage.handleCancel}
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

TeacherFeedbackPage.propTypes = propTypes;
