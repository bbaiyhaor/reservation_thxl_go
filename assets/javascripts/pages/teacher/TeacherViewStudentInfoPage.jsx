import 'weui';
import { AlertDialog, LoadingHud } from '#coms/Huds';
import { Cell, CellBody, CellFooter, Cells, CellsTitle, Panel, PanelHeader, TextArea } from 'react-weui';
import { Application } from '#models/Models';
import PageBottom from '#coms/PageBottom';
import PropTypes from 'prop-types';
import React from 'react';
import queryString from 'query-string';

export default class TeacherViewStudentInfoPage extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      student: null,
      reservations: null,
    };
  }

  componentDidMount() {
    this.loading.show('正在加载中');
    const parsedQuery = queryString.parse(this.props.history.location.search);
    const studentId = parsedQuery.student_id ? parsedQuery.student_id : '';
    const studentUsername = parsedQuery.student_username ? parsedQuery.student_username : '';
    if (studentId !== '') {
      setTimeout(() => {
        Application.getStudentInfoByTeacher(studentId, (data) => {
          this.loading.hide();
          this.setState({
            student: data.student,
            reservations: data.reservations,
          });
        }, (error) => {
          this.loading.hide();
          this.alert.show('', error, '好的', () => {
            this.props.history.push('/reservation');
          });
        });
      }, 500);
    } else if (studentUsername && studentUsername !== '') {
      setTimeout(() => {
        Application.queryStudentInfoByTeacher(studentUsername, (data) => {
          this.loading.hide();
          this.setState({
            student: data.student,
            reservations: data.reservations,
          });
        }, (error) => {
          this.loading.hide();
          this.alert.show('', error, '好的', () => {
            this.props.history.push('/reservation');
          });
        });
      }, 500);
    } else {
      this.loading.hide();
      this.props.history.push('/reservation');
    }
  }

  render() {
    return (
      <div>
        <Panel>
          <PanelHeader style={{ fontSize: '18px' }}>学生信息表</PanelHeader>
          <StudentInfoPanel
            student={this.state.student}
            reservations={this.state.reservations}
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

TeacherViewStudentInfoPage.propTypes = {
  history: PropTypes.object.isRequired,
};

class StudentInfoPanel extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      student: null,
      reservations: null,
      reservationShows: {},
    };
    this.showStudentReservation = this.showStudentReservation.bind(this);
  }

  componentWillReceiveProps(nextProps) {
    const reservationShows = {};
    if (nextProps.reservations) {
      for (let i = 0; i < nextProps.reservations.length; i += 1) {
        reservationShows[nextProps.reservations[i].id] = false;
      }
    }
    this.setState({
      student: nextProps.student,
      reservations: nextProps.reservations,
      reservationShows,
    });
  }

  showStudentReservation(reservationId) {
    const reservationShows = {};
    for (let i = 0; i < this.state.reservations.length; i += 1) {
      reservationShows[this.state.reservations[i].id] = false;
    }
    reservationShows[reservationId] = !this.state.reservationShows[reservationId];
    this.setState({ reservationShows });
  }

  renderStudentReservations() {
    if (this.state.reservations && this.state.reservations.length > 0) {
      return (
        <Cells style={{ marginTop: '0px' }}>
          <CellsTitle>咨询经历</CellsTitle>
          {this.state.reservations.map((reservation) => {
            const cellBodyStyle = reservation.status !== 3 ? {
              color: '#888',
            } : null;
            const cellFooterStyle = this.state.reservationShows[reservation.id] ? {
              transform: 'rotate(90deg)',
            } : null;
            return (
              <div key={`student_reservation_div_${reservation.id}`}>
                <Cells style={{ marginTop: '0px' }}>
                  <Cell onClick={() => { this.showStudentReservation(reservation.id); }}>
                    <CellBody style={{ ...cellBodyStyle }}>
                      {reservation.start_time}-{reservation.end_time.slice(-5)} {reservation.teacher_fullname}
                    </CellBody>
                    <CellFooter style={{ ...cellFooterStyle }} />
                  </Cell>
                </Cells>
                {this.state.reservationShows[reservation.id] && <StudentReservationCell reservation={reservation} />}
              </div>
            );
          })}
        </Cells>
      );
    }
    return null;
  }

  render() {
    return (
      <div>
        <Cells>
          <Cell>
            <CellBody>学号</CellBody>
            <CellFooter>{this.state.student ? this.state.student.username : ''}</CellFooter>
          </Cell>
          <Cell>
            <CellBody>姓名</CellBody>
            <CellFooter>{this.state.student ? this.state.student.fullname : ''}</CellFooter>
          </Cell>
          <Cell>
            <CellBody>性别</CellBody>
            <CellFooter>{this.state.student ? this.state.student.gender : ''}</CellFooter>
          </Cell>
          <Cell>
            <CellBody>出生日期</CellBody>
            <CellFooter>{this.state.student ? this.state.student.birthday : ''}</CellFooter>
          </Cell>
          <Cell>
            <CellBody>院系</CellBody>
            <CellFooter>{this.state.student ? this.state.student.school : ''}</CellFooter>
          </Cell>
          <Cell>
            <CellBody>年级</CellBody>
            <CellFooter>{this.state.student ? this.state.student.grade : ''}</CellFooter>
          </Cell>
          <Cell>
            <CellBody>现在住址</CellBody>
            <CellFooter>{this.state.student ? this.state.student.current_address : ''}</CellFooter>
          </Cell>
          <Cell>
            <CellBody>家庭住址</CellBody>
            <CellFooter>{this.state.student ? this.state.student.family_address : ''}</CellFooter>
          </Cell>
          <Cell>
            <CellBody>联系电话</CellBody>
            <CellFooter>{this.state.student ? this.state.student.mobile : ''}</CellFooter>
          </Cell>
          <Cell>
            <CellBody>邮箱</CellBody>
            <CellFooter>{this.state.student ? this.state.student.email : ''}</CellFooter>
          </Cell>
          <CellsTitle>过往咨询经历</CellsTitle>
          <Cell>
            <CellBody>时间</CellBody>
            <CellFooter>{this.state.student ? this.state.student.experience_time : ''}</CellFooter>
          </Cell>
          <Cell>
            <CellBody>地点</CellBody>
            <CellFooter>{this.state.student ? this.state.student.experience_location : ''}</CellFooter>
          </Cell>
          <Cell>
            <CellBody>咨询师姓名</CellBody>
            <CellFooter>{this.state.student ? this.state.student.experience_teacher : ''}</CellFooter>
          </Cell>
          <CellsTitle>家庭情况</CellsTitle>
          <Cell>
            <CellBody>父亲年龄</CellBody>
            <CellFooter>{this.state.student ? this.state.student.father_age : ''}</CellFooter>
          </Cell>
          <Cell>
            <CellBody>父亲职业</CellBody>
            <CellFooter>{this.state.student ? this.state.student.father_job : ''}</CellFooter>
          </Cell>
          <Cell>
            <CellBody>父亲学历</CellBody>
            <CellFooter>{this.state.student ? this.state.student.father_edu : ''}</CellFooter>
          </Cell>
          <Cell>
            <CellBody>母亲年龄</CellBody>
            <CellFooter>{this.state.student ? this.state.student.mother_age : ''}</CellFooter>
          </Cell>
          <Cell>
            <CellBody>母亲职业</CellBody>
            <CellFooter>{this.state.student ? this.state.student.mother_job : ''}</CellFooter>
          </Cell>
          <Cell>
            <CellBody>母亲学历</CellBody>
            <CellFooter>{this.state.student ? this.state.student.mother_edu : ''}</CellFooter>
          </Cell>
          <Cell>
            <CellBody>父母婚姻状况</CellBody>
            <CellFooter>{this.state.student ? this.state.student.parent_marriage : ''}</CellFooter>
          </Cell>
          <CellsTitle>
            在近三个月里，是否发生了对你有重大意义的事（如亲友的死亡、法律诉讼、失恋等）？
          </CellsTitle>
          <Cell>
            <CellBody>
              <TextArea
                rows="3"
                value={this.state.student ? this.state.student.significant : ''}
                disabled
              />
            </CellBody>
          </Cell>
          <CellsTitle>
            你现在需要接受帮助的主要问题是什么？
          </CellsTitle>
          <Cell>
            <CellBody>
              <TextArea
                rows="3"
                value={this.state.student ? this.state.student.problem : ''}
                disabled
              />
            </CellBody>
          </Cell>
          <Cell>
            <CellBody>档案分类</CellBody>
            <CellFooter>{this.state.student ? this.state.student.archive_category : ''}</CellFooter>
          </Cell>
          <Cell>
            <CellBody>档案编号</CellBody>
            <CellFooter>{this.state.student ? this.state.student.archive_number : ''}</CellFooter>
          </Cell>
          <Cell>
            <CellBody>是否危机个案</CellBody>
            <CellFooter>{this.state.student && this.state.student.crisis_level > 0 ? '是' : '否'}</CellFooter>
          </Cell>
          <Cell>
            <CellBody>绑定咨询师</CellBody>
            <CellFooter>{this.state.student ? `${this.state.student.binded_teacher_username} ${this.state.student.binded_teacher_fullname}` : ''}</CellFooter>
          </Cell>
        </Cells>
        {this.renderStudentReservations()}
        <div style={{ height: '10px' }} />
      </div>
    );
  }
}

StudentInfoPanel.propTypes = {
  reservations: PropTypes.arrayOf(PropTypes.object),
  student: PropTypes.object,
};

StudentInfoPanel.defaultProps = {
  reservations: [],
  student: null,
};

function StudentReservationCell({ reservation }) {
  return (
    <Cells style={{ marginTop: '0px' }}>
      <Cell>
        <CellBody>学生反馈</CellBody>
        <CellFooter>{reservation.student_feedback.scores}</CellFooter>
      </Cell>
      <Cell>
        <CellBody>评估分类</CellBody>
        <CellFooter>{reservation.teacher_feedback.category}</CellFooter>
      </Cell>
      {reservation.teacher_feedback.severity !== '' &&
      <div>
        <CellsTitle>
          严重程度
        </CellsTitle>
        <Cell>
          <CellBody>
            <TextArea
              value={reservation.teacher_feedback.severity}
              disabled
            />
          </CellBody>
        </Cell>
      </div>
      }
      {reservation.teacher_feedback.medical_diagnosis !== '' &&
      <div>
        <CellsTitle>
          疑似或明确的医疗诊断
        </CellsTitle>
        <Cell>
          <CellBody>
            <TextArea
              value={reservation.teacher_feedback.medical_diagnosis}
              disabled
            />
          </CellBody>
        </Cell>
      </div>
      }
      {reservation.teacher_feedback.crisis !== '' &&
      <div>
        <CellsTitle>
          危急情况
        </CellsTitle>
        <Cell>
          <CellBody>
            <TextArea
              value={reservation.teacher_feedback.crisis}
              disabled
            />
          </CellBody>
        </Cell>
      </div>
      }
      <CellsTitle>咨询记录</CellsTitle>
      <Cell>
        <CellBody>
          <TextArea
            rows="3"
            value={reservation.teacher_feedback.record}
            disabled
          />
        </CellBody>
      </Cell>
    </Cells>
  );
}

StudentReservationCell.propTypes = {
  reservation: PropTypes.object.isRequired,
};
