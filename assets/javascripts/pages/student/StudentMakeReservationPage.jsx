import 'weui';
import { AlertDialog, LoadingHud } from '#coms/Huds';
import { Application, User } from '#models/Models';
import { Panel, PanelHeader } from 'react-weui';
import React, { PropTypes } from 'react';
import MakeReservationForm from '#forms/MakeReservationForm';
import PageBottom from '#coms/PageBottom';

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

  componentWillMount() {
    const reservationId = this.props.history.location.state.reservation_id || '';
    if (reservationId === '' || !User.student || !Application.reservations || Application.reservations.length === 0) {
      this.props.history.push('/reservation');
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
      this.props.history.push('/reservation');
    }
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
