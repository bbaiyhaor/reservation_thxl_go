/* eslint max-len: ["off"] */
import 'weui';
import { Button, ButtonArea, CellBody, CellFooter, CellHeader, CellsTitle, Form, FormCell, Icon, Input, Label, Select, TextArea } from 'react-weui';
import PropTypes from 'prop-types';
import React from 'react';

export default class MakeReservationForm extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      startTime: '',
      endTime: '',
      fullnameWarn: false,
      genderWarn: false,
      birthdayWarn: false,
      schoolWarn: false,
      gradeWarn: false,
      currentAddressWarn: false,
      familyAddressWarn: false,
      mobileWarn: false,
      emailWarn: false,
      problemWarn: false,
    };
    this.handleChange = this.handleChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  componentWillReceiveProps(nextProps) {
    nextProps.student && this.setState({
      fullname: nextProps.student.fullname,
      gender: nextProps.student.gender,
      birthday: nextProps.student.birthday,
      school: nextProps.student.school,
      grade: nextProps.student.grade,
      currentAddress: nextProps.student.current_address,
      familyAddress: nextProps.student.family_address,
      mobile: nextProps.student.mobile,
      email: nextProps.student.email,
      experienceTime: nextProps.student.experience_time,
      experienceLocation: nextProps.student.experience_location,
      experienceTeacher: nextProps.student.experience_teacher,
      fatherAge: nextProps.student.father_age,
      fatherJob: nextProps.student.father_job,
      fatherEdu: nextProps.student.father_edu,
      motherAge: nextProps.student.mother_age,
      motherJob: nextProps.student.mother_job,
      motherEdu: nextProps.student.mother_edu,
      parentMarriage: nextProps.student.parent_marriage,
      significant: nextProps.student.significant,
      problem: nextProps.student.problem,
    });
    nextProps.reservation && this.setState({
      reservation: nextProps.reservation,
      startTime: nextProps.reservation.start_time,
      endTime: nextProps.reservation.end_time,
      teacherFullname: nextProps.reservation.teachaer_fullname,
    });
  }

  handleChange(e) {
    const target = e.target;
    const name = target.name;
    const value = target.value;
    this.setState({
      [name]: value,
    });
  }

  handleSubmit() {
    this.setState({
      fullnameWarn: false,
      genderWarn: false,
      birthdayWarn: false,
      schoolWarn: false,
      gradeWarn: false,
      currentAddressWarn: false,
      familyAddressWarn: false,
      mobileWarn: false,
      emailWarn: false,
      problemWarn: false,
    });
    if (this.state.fullname === '') {
      this.setState({ fullnameWarn: true });
      this.fullnameInput.focus();
      return;
    }
    if (this.state.gender === '') {
      this.setState({ genderWarn: true });
      this.gradeInput.focus();
      return;
    }
    if (this.state.birthday === '') {
      this.setState({ birthdayWarn: true });
      this.birthdayInput.focus();
      return;
    }
    if (this.state.school === '') {
      this.setState({ schoolWarn: true });
      this.schoolInput.focus();
      return;
    }
    if (this.state.grade === '') {
      this.setState({ gradeWarn: true });
      this.gradeInput.focus();
      return;
    }
    if (this.state.currentAddress === '') {
      this.setState({ currentAddressWarn: true });
      this.currentAddressInput.focus();
      return;
    }
    if (this.state.familyAddress === '') {
      this.setState({ familyAddressWarn: true });
      this.familyAddressInput.focus();
      return;
    }
    if (this.state.mobile === '') {
      this.setState({ mobileWarn: true });
      this.mobileInput.focus();
      return;
    }
    if (this.state.email === '') {
      this.setState({ emailWarn: true });
      this.emailInput.focus();
      return;
    }
    if (this.state.problem === '') {
      this.setState({ problemWarn: true });
      this.problemInput.focus();
      return;
    }
    this.props.handleSubmit(this.state.reservation, this.state.fullname, this.state.gender, this.state.birthday, this.state.school, this.state.grade, this.state.currentAddress, this.state.familyAddress, this.state.mobile, this.state.email, this.state.experienceTime, this.state.experienceLocation, this.state.experienceTeacher, this.state.fatherAge, this.state.fatherJob, this.state.fatherEdu, this.state.motherAge, this.state.motherJob, this.state.motherEdu, this.state.parentMarriage, this.state.significant, this.state.problem);
  }

  render() {
    return (
      <div>
        {this.props.reservation &&
          <CellsTitle>
            正在预约：{this.props.reservation.start_time}-{this.props.reservation.end_time.slice(-5)} {this.props.reservation.teacher_fullname}
          </CellsTitle>
        }
        <Form>
          <FormCell warn={this.state.fullnameWarn}>
            <CellHeader>
              <Label>姓名<span style={{ color: 'red' }}>*</span></Label>
            </CellHeader>
            <CellBody>
              <Input
                ref={(fullnameInput) => { this.fullnameInput = fullnameInput; }}
                name="fullname"
                type="text"
                placeholder="请输入姓名"
                value={this.state.fullname ? this.state.fullname : ''}
                onChange={this.handleChange}
              />
            </CellBody>
            {this.state.fullnameWarn &&
              <CellFooter>
                <Icon value="warn" />
              </CellFooter>
            }
          </FormCell>
          <FormCell warn={this.state.genderWarn} select selectPos="after">
            <CellHeader>
              <Label>性别<span style={{ color: 'red' }}>*</span></Label>
            </CellHeader>
            <CellBody>
              <Select
                name="gender"
                value={this.state.gender ? this.state.gender : ''}
                onChange={this.handleChange}
              >
                <option value="">请选择</option>
                <option value="男">男</option>
                <option value="女">女</option>
              </Select>
            </CellBody>
            {this.state.genderWarn &&
              <CellFooter style={{ marginRight: '25px' }}>
                <Icon value="warn" />
              </CellFooter>
            }
          </FormCell>
          <FormCell warn={this.state.birthdayWarn}>
            <CellHeader>
              <Label>出生日期<span style={{ color: 'red' }}>*</span></Label>
            </CellHeader>
            <CellBody>
              <Input
                ref={(birthdayInput) => { this.birthdayInput = birthdayInput; }}
                name="birthday"
                type="date"
                placeholder="请输入出生日期"
                value={this.state.birthday ? this.state.birthday : ''}
                onChange={this.handleChange}
              />
            </CellBody>
            {this.state.birthdayWarn &&
              <CellFooter>
                <Icon value="warn" />
              </CellFooter>
            }
          </FormCell>
          <FormCell warn={this.state.schoolWarn}>
            <CellHeader>
              <Label>院系<span style={{ color: 'red' }}>*</span></Label>
            </CellHeader>
            <CellBody>
              <Input
                ref={(schoolInput) => { this.schoolInput = schoolInput; }}
                name="school"
                type="text"
                placeholder="请输入院系"
                value={this.state.school ? this.state.school : ''}
                onChange={this.handleChange}
              />
            </CellBody>
            {this.state.schoolWarn &&
              <CellFooter>
                <Icon value="warn" />
              </CellFooter>
            }
          </FormCell>
          <FormCell warn={this.state.gradeWarn}>
            <CellHeader>
              <Label>年级<span style={{ color: 'red' }}>*</span></Label>
            </CellHeader>
            <CellBody>
              <Input
                ref={(gradeInput) => { this.gradeInput = gradeInput; }}
                name="garde"
                type="text"
                placeholder="请输入年级"
                value={this.state.grade ? this.state.grade : ''}
                onChange={this.handleChange}
              />
            </CellBody>
            {this.state.gradeWarn &&
              <CellFooter>
                <Icon value="warn" />
              </CellFooter>
            }
          </FormCell>
          <FormCell warn={this.state.currentAddressWarn}>
            <CellHeader>
              <Label>现在住址<span style={{ color: 'red' }}>*</span></Label>
            </CellHeader>
            <CellBody>
              <Input
                ref={(currentAddressInput) => { this.currentAddressInput = currentAddressInput; }}
                name="currentAddress"
                type="text"
                placeholder="请输入现在住址"
                value={this.state.currentAddress ? this.state.currentAddress : ''}
                onChange={this.handleChange}
              />
            </CellBody>
            {this.state.currentAddressWarn &&
              <CellFooter>
                <Icon value="warn" />
              </CellFooter>
            }
          </FormCell>
          <FormCell warn={this.state.familyAddressWarn}>
            <CellHeader>
              <Label>家庭住址<span style={{ color: 'red' }}>*</span></Label>
            </CellHeader>
            <CellBody>
              <Input
                ref={(familyAddressInput) => { this.familyAddressInput = familyAddressInput; }}
                name="familyAddress"
                type="text"
                placeholder="请输入家庭住址"
                value={this.state.familyAddress ? this.state.familyAddress : ''}
                onChange={this.handleChange}
              />
            </CellBody>
            {this.state.familyAddressWarn &&
              <CellFooter>
                <Icon value="warn" />
              </CellFooter>
            }
          </FormCell>
          <FormCell warn={this.state.mobileWarn}>
            <CellHeader>
              <Label>联系电话<span style={{ color: 'red' }}>*</span></Label>
            </CellHeader>
            <CellBody>
              <Input
                ref={(mobileInput) => { this.mobileInput = mobileInput; }}
                name="mobile"
                type="tel"
                placeholder="请输入联系电话"
                value={this.state.mobile ? this.state.mobile : ''}
                onChange={this.handleChange}
              />
            </CellBody>
            {this.state.mobileWarn &&
              <CellFooter>
                <Icon value="warn" />
              </CellFooter>
            }
          </FormCell>
          <FormCell warn={this.state.emailWarn}>
            <CellHeader>
              <Label>邮箱<span style={{ color: 'red' }}>*</span></Label>
            </CellHeader>
            <CellBody>
              <Input
                ref={(emailInput) => { this.emailInput = emailInput; }}
                name="email"
                type="text"
                placeholder="请输入邮箱"
                value={this.state.email ? this.state.email : ''}
                onChange={this.handleChange}
              />
            </CellBody>
            {this.state.emailWarn &&
              <CellFooter>
                <Icon value="warn" />
              </CellFooter>
            }
          </FormCell>
          <CellsTitle>过往咨询经历</CellsTitle>
          <FormCell>
            <CellHeader>
              <Label>时间</Label>
            </CellHeader>
            <CellBody>
              <Input
                name="experienceTime"
                type="text"
                placeholder="请输入咨询时间"
                value={this.state.experienceTime ? this.state.experienceTime : ''}
                onChange={this.handleChange}
              />
            </CellBody>
          </FormCell>
          <FormCell>
            <CellHeader>
              <Label>地点</Label>
            </CellHeader>
            <CellBody>
              <Input
                name="experienceLocation"
                type="text"
                placeholder="请输入咨询地点"
                value={this.state.experienceLocation ? this.state.experienceLocation : ''}
                onChange={this.handleChange}
              />
            </CellBody>
          </FormCell>
          <FormCell>
            <CellHeader>
              <Label>咨询师姓名</Label>
            </CellHeader>
            <CellBody>
              <Input
                name="experienceTeacher"
                type="text"
                placeholder="请输入咨询师姓名"
                value={this.state.experienceTeacher ? this.state.experienceTeacher : ''}
                onChange={this.handleChange}
              />
            </CellBody>
          </FormCell>
          <CellsTitle>家庭情况</CellsTitle>
          <FormCell>
            <CellHeader>
              <Label>父亲年龄</Label>
            </CellHeader>
            <CellBody>
              <Input
                name="fatherAge"
                type="text"
                placeholder="请输入父亲年龄"
                value={this.state.fatherAge ? this.state.fatherAge : ''}
                onChange={this.handleChange}
              />
            </CellBody>
          </FormCell>
          <FormCell>
            <CellHeader>
              <Label>父亲职业</Label>
            </CellHeader>
            <CellBody>
              <Input
                name="fatherJob"
                type="text"
                placeholder="请输入父亲职业"
                value={this.state.fatherJob ? this.state.fatherJob : ''}
                onChange={this.handleChange}
              />
            </CellBody>
          </FormCell>
          <FormCell>
            <CellHeader>
              <Label>父亲学历</Label>
            </CellHeader>
            <CellBody>
              <Input
                name="fatherEdu"
                type="text"
                placeholder="请输入父亲学历"
                value={this.state.fatherEdu ? this.state.fatherEdu : ''}
                onChange={this.handleChange}
              />
            </CellBody>
          </FormCell>
          <FormCell>
            <CellHeader>
              <Label>母亲年龄</Label>
            </CellHeader>
            <CellBody>
              <Input
                name="motherAge"
                type="text"
                placeholder="请输入母亲年龄"
                value={this.state.motherAge ? this.state.motherAge : ''}
                onChange={this.handleChange}
              />
            </CellBody>
          </FormCell>
          <FormCell>
            <CellHeader>
              <Label>母亲职业</Label>
            </CellHeader>
            <CellBody>
              <Input
                name="motherJob"
                type="text"
                placeholder="请输入母亲职业"
                value={this.state.motherJob ? this.state.motherJob : ''}
                onChange={this.handleChange}
              />
            </CellBody>
          </FormCell>
          <FormCell>
            <CellHeader>
              <Label>母亲学历</Label>
            </CellHeader>
            <CellBody>
              <Input
                name="motherEdu"
                type="text"
                placeholder="请输入母亲学历"
                value={this.state.motherEdu ? this.state.motherEdu : ''}
                onChange={this.handleChange}
              />
            </CellBody>
          </FormCell>
          <FormCell select selectPos="after">
            <CellHeader>
              <Label>父母婚姻状况</Label>
            </CellHeader>
            <CellBody>
              <Select
                name="parentMarriage"
                value={this.state.parentMarriage ? this.state.parentMarriage : ''}
                onChange={this.handleChange}
              >
                <option value="">请选择</option>
                <option value="良好">良好</option>
                <option value="一般">一般</option>
                <option value="离婚">离婚</option>
                <option value="再婚">再婚</option>
              </Select>
            </CellBody>
          </FormCell>
          <CellsTitle>
            在近三个月里，是否发生了对你有重大意义的事（如亲友的死亡、法律诉讼、失恋等）？
          </CellsTitle>
          <FormCell>
            <CellBody>
              <TextArea
                name="significant"
                placeholder="请输入"
                rows="3"
                maxLength={300}
                value={this.state.significant ? this.state.significant : ''}
                onChange={this.handleChange}
              />
            </CellBody>
          </FormCell>
          <CellsTitle>
            你现在需要接受帮助的主要问题是什么？<span style={{ color: 'red' }}>*</span>
          </CellsTitle>
          <FormCell warn={this.state.problemWarn}>
            <CellBody>
              <TextArea
                ref={(problemInput) => { this.problemInput = problemInput; }}
                name="problem"
                placeholder="请输入"
                rows="3"
                maxLength={300}
                value={this.state.problem ? this.state.problem : ''}
                onChange={this.handleChange}
              />
            </CellBody>
            {this.state.problemWarn &&
              <CellFooter>
                <Icon value="warn" />
              </CellFooter>
            }
          </FormCell>
        </Form>
        <ButtonArea direction="horizontal">
          <Button onClick={this.handleSubmit}>确定预约</Button>
          <Button type="default" onClick={this.props.handleCancel}>取消预约</Button>
        </ButtonArea>
        <div style={{ height: '10px' }} />
      </div>
    );
  }
}

MakeReservationForm.propTypes = {
  student: PropTypes.object,
  reservation: PropTypes.object,
  handleSubmit: PropTypes.func.isRequired,
  handleCancel: PropTypes.func.isRequired,
};

MakeReservationForm.defaultProps = {
  student: null,
  reservation: null,
};
