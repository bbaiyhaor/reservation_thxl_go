/* eslint max-len: ["off"] */
import 'weui';
import { Button, ButtonArea, CellBody, CellFooter, CellHeader, CellsTitle, Form, FormCell, Icon, Input, Label, Select, TextArea } from 'react-weui';
import React, { PropTypes } from 'react';

export default class MakeReservationForm extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      fullname: this.props.student.fullname,
      gender: this.props.student.gender,
      birthday: this.props.student.birthday,
      school: this.props.student.school,
      grade: this.props.student.grade,
      currentAddress: this.props.student.current_address,
      familyAddress: this.props.student.family_address,
      mobile: this.props.student.mobile,
      email: this.props.student.email,
      experienceTime: this.props.student.experience_time,
      experienceLocation: this.props.student.experience_location,
      experienceTeacher: this.props.student.experience_teacher,
      fatherAge: this.props.student.father_age,
      fatherJob: this.props.student.father_job,
      fatherEdu: this.props.student.father_edu,
      motherAge: this.props.student.mother_age,
      motherJob: this.props.student.mother_job,
      motherEdu: this.props.student.mother_edu,
      parentMarriage: this.props.student.parent_marriage,
      significant: this.props.student.significant,
      problem: this.props.student.problem,
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

  handleChange(e, name) {
    const value = e.target.value;
    if (name && name !== '') {
      this.setState({ [name]: value });
    } else {
      this.setState(prevState => ({
        [value]: !prevState[value],
      }));
    }
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
    this.props.handleSubmit(this.props.reservation, this.state.fullname, this.state.gender, this.state.birthday, this.state.school, this.state.grade, this.state.currentAddress, this.state.familyAddress, this.state.mobile, this.state.email, this.state.experienceTime, this.state.experienceLocation, this.state.experienceTeacher, this.state.fatherAge, this.state.fatherJob, this.state.fatherEdu, this.state.motherAge, this.state.motherJob, this.state.motherEdu, this.state.parentMarriage, this.state.significant, this.state.problem);
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
                type="input"
                placeholder="请输入姓名"
                value={this.state.fullname}
                onChange={(e) => { this.handleChange(e, 'fullname'); }}
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
                value={this.state.gender}
                onChange={(e) => { this.handleChange(e, 'gender'); }}
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
                type="input"
                placeholder="请输入出生日期"
                value={this.state.birthday}
                onChange={(e) => { this.handleChange(e, 'birthday'); }}
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
                type="input"
                placeholder="请输入院系"
                value={this.state.school}
                onChange={(e) => { this.handleChange(e, 'school'); }}
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
                type="input"
                placeholder="请输入年级"
                value={this.state.grade}
                onChange={(e) => { this.handleChange(e, 'grade'); }}
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
                type="input"
                placeholder="请输入现在住址"
                value={this.state.currentAddress}
                onChange={(e) => { this.handleChange(e, 'currentAddress'); }}
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
                type="input"
                placeholder="请输入家庭住址"
                value={this.state.familyAddress}
                onChange={(e) => { this.handleChange(e, 'familyAddress'); }}
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
                type="tel"
                placeholder="请输入联系电话"
                value={this.state.mobile}
                onChange={(e) => { this.handleChange(e, 'mobile'); }}
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
                type="input"
                placeholder="请输入邮箱"
                value={this.state.email}
                onChange={(e) => { this.handleChange(e, 'email'); }}
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
                type="input"
                placeholder="请输入咨询时间"
                value={this.state.experienceTime}
                onChange={(e) => { this.handleChange(e, 'experienceTime'); }}
              />
            </CellBody>
          </FormCell>
          <FormCell>
            <CellHeader>
              <Label>地点</Label>
            </CellHeader>
            <CellBody>
              <Input
                type="input"
                placeholder="请输入咨询地点"
                value={this.state.experienceLocation}
                onChange={(e) => { this.handleChange(e, 'experienceLocation'); }}
              />
            </CellBody>
          </FormCell>
          <FormCell>
            <CellHeader>
              <Label>咨询师姓名</Label>
            </CellHeader>
            <CellBody>
              <Input
                type="input"
                placeholder="请输入咨询师姓名"
                value={this.state.experienceTeacher}
                onChange={(e) => { this.handleChange(e, 'experienceTeacher'); }}
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
                type="input"
                placeholder="请输入父亲年龄"
                value={this.state.fatherAge}
                onChange={(e) => { this.handleChange(e, 'fatherAge'); }}
              />
            </CellBody>
          </FormCell>
          <FormCell>
            <CellHeader>
              <Label>父亲职业</Label>
            </CellHeader>
            <CellBody>
              <Input
                type="input"
                placeholder="请输入父亲职业"
                value={this.state.fatherJob}
                onChange={(e) => { this.handleChange(e, 'fatherJob'); }}
              />
            </CellBody>
          </FormCell>
          <FormCell>
            <CellHeader>
              <Label>父亲学历</Label>
            </CellHeader>
            <CellBody>
              <Input
                type="input"
                placeholder="请输入父亲学历"
                value={this.state.fatherEdu}
                onChange={(e) => { this.handleChange(e, 'fatherEdu'); }}
              />
            </CellBody>
          </FormCell>
          <FormCell>
            <CellHeader>
              <Label>母亲年龄</Label>
            </CellHeader>
            <CellBody>
              <Input
                type="input"
                placeholder="请输入母亲年龄"
                value={this.state.motherAge}
                onChange={(e) => { this.handleChange(e, 'motherAge'); }}
              />
            </CellBody>
          </FormCell>
          <FormCell>
            <CellHeader>
              <Label>母亲职业</Label>
            </CellHeader>
            <CellBody>
              <Input
                type="input"
                placeholder="请输入母亲职业"
                value={this.state.motherJob}
                onChange={(e) => { this.handleChange(e, 'motherJob'); }}
              />
            </CellBody>
          </FormCell>
          <FormCell>
            <CellHeader>
              <Label>母亲学历</Label>
            </CellHeader>
            <CellBody>
              <Input
                type="input"
                placeholder="请输入母亲学历"
                value={this.state.motherEdu}
                onChange={(e) => { this.handleChange(e, 'motherEdu'); }}
              />
            </CellBody>
          </FormCell>
          <FormCell select selectPos="after">
            <CellHeader>
              <Label>父母婚姻状况</Label>
            </CellHeader>
            <CellBody>
              <Select
                value={this.state.parentMarriage}
                onChange={(e) => { this.handleChange(e, 'parentMarriage'); }}
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
                placeholder="请输入"
                rows="3"
                maxLength={300}
                value={this.state.significant}
                onChange={(e) => { this.handleChange(e, 'significant'); }}
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
                placeholder="请输入"
                rows="3"
                maxLength={300}
                value={this.state.problem}
                onChange={(e) => { this.handleChange(e, 'problem'); }}
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
  student: PropTypes.object.isRequired,
  reservation: PropTypes.object.isRequired,
  handleSubmit: PropTypes.func.isRequired,
  handleCancel: PropTypes.func.isRequired,
};
