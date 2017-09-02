import 'weui';
import { AlertDialog, LoadingHud } from '#coms/Huds';
import { Application, User } from '#models/Models';
import { Panel, PanelHeader } from 'react-weui';
import MakeReservationForm from '#forms/MakeReservationForm';
import PageBottom from '#coms/PageBottom';
import PropTypes from 'prop-types';
import React from 'react';
import queryString from 'query-string';

export default class StudentMakeReservationPage extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      student: null,
      reservation: null,
    };
    this.handleCancel = this.handleCancel.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  componentDidMount() {
    this.loading.show('正在加载中');
    const parsedQuery = queryString.parse(this.props.history.location.search);
    const reservationId = parsedQuery.reservation_id;
    const sourceId = parsedQuery.source_id;
    const startTime = parsedQuery.start_time;
    User.updateSession(() => {
      Application.validReservationByStudent(reservationId, sourceId, startTime, () => {
        setTimeout(() => {
          this.setState({
            student: User.student,
            reservation: Application.reservation,
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

  handleSubmit(reservation, fullname, gender, birthday, school, grade, currentAddress, familyAddress, mobile, email, experienceTime, experienceLocation, experienceTeacher, fatherAge, fatherJob, fatherEdu, motherAge, motherJob, motherEdu, parentMarriage, significant, problem) {
    const payload = {
      reservation_id: reservation.id,
      source_id: reservation.source_id,
      start_time: reservation.start_time,
      student_fullname: fullname,
      student_gender: gender,
      student_birthday: birthday,
      student_school: school,
      student_grade: grade,
      student_current_address: currentAddress,
      student_family_address: familyAddress,
      student_mobile: mobile,
      student_email: email,
      student_experience_time: experienceTime,
      student_experience_location: experienceLocation,
      student_experience_teacher: experienceTeacher,
      student_father_age: fatherAge,
      student_father_job: fatherJob,
      student_father_edu: fatherEdu,
      student_mother_age: motherAge,
      student_mother_job: motherJob,
      student_mother_edu: motherEdu,
      student_parent_marriage: parentMarriage,
      student_significant: significant,
      student_problem: problem,
    };
    this.loading.show('正在加载中');
    setTimeout(() => {
      Application.makeReservationByStudent(payload, () => {
        this.loading.hide();
        this.props.history.push('/reservation/make/success');
      }, (error) => {
        this.loading.hide();
        this.alert.show('预约失败', error, '好的');
      });
    }, 500);
  }

  handleCancel() {
    this.props.history.goBack();
  }

  render() {
    return (
      <div>
        <Panel>
          <PanelHeader style={{ fontSize: '18px' }}>学生信息登记表</PanelHeader>
          <MakeReservationForm
            student={this.state.student}
            reservation={this.state.reservation}
            handleSubmit={this.handleSubmit}
            handleCancel={this.handleCancel}
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

StudentMakeReservationPage.propTypes = {
  history: PropTypes.object.isRequired,
};
