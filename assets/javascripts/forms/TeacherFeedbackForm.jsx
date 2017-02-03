/**
 * Created by shudi on 2016/11/4.
 */
import React, { PropTypes } from 'react';
import { Form, FormCell, CellsTitle, CellHeader, CellBody, CellFooter, Label, Select, Checkbox, TextArea, Icon, ButtonArea, Button } from '#react-weui';
import 'weui';

const propTypes = {
  handleSubmit: PropTypes.func.isRequired,
  handleCancel: PropTypes.func.isRequired,
  showAlert: PropTypes.func.isRequired,
};

export default class TeacherFeedbackForm extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      reservation: null,
      feedback: null,
      firstCategories: null,
      secondCategories: null,
      firstCategory: '',
      secondCategory: '',
      severity: [],
      varSeverity: [],
      medicalDiagnosis: [],
      varMedicalDiagnosis: [],
      crisis: [],
      varCrisis: [],
      record: '',
      crisisLevel: '0',
      firstCategoryWarn: false,
      secondCategoryWarn: false,
      recordWarn: false,
      categoryShowTips: '',
      categoryShowCheckTips: ['A1', 'A2', 'E1', 'E2', 'E3', 'F1', 'F2', 'F3', 'F4', 'F5', 'F6', 'H1', 'H2', 'H3', 'H4', 'H5', 'H6'],
      categoryShowNeedTips: ['G1', 'G2', 'G3', 'G4', 'J1', 'J2'],
    };
    this.handleChange = this.handleChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  componentWillReceiveProps(nextProps) {
    this.setState({
      reservation: nextProps.reservation,
      feedback: nextProps.feedback,
    });
    if (nextProps.feedback) {
      let severity = nextProps.feedback.severity;
      const varSeverity = nextProps.feedback.var_severity;
      while (severity.length < varSeverity.length) {
        severity.push(0);
      }
      severity = severity.slice(0, varSeverity.length);

      let medicalDiagnosis = nextProps.feedback.medical_diagnosis;
      const varMedicalDiagnosis = nextProps.feedback.var_medical_diagnosis;
      while (medicalDiagnosis.length < varMedicalDiagnosis.length) {
        medicalDiagnosis.push(0);
      }
      medicalDiagnosis = medicalDiagnosis.slice(0, varMedicalDiagnosis.length);

      let crisis = nextProps.feedback.crisis;
      const varCrisis = nextProps.feedback.var_crisis;
      while (crisis.length < varCrisis.length) {
        crisis.push(0);
      }
      crisis = crisis.slice(0, varCrisis.length);

      this.setState({
        firstCategories: nextProps.feedback.var_first_category,
        secondCategories: nextProps.feedback.var_second_category,
        firstCategory: nextProps.feedback.category ? nextProps.feedback.category.substring(0, 1) : '',
        secondCategory: nextProps.feedback.category ? nextProps.feedback.category : '',
        severity,
        varSeverity,
        medicalDiagnosis,
        varMedicalDiagnosis,
        crisis,
        varCrisis,
        record: nextProps.feedback.record ? nextProps.feedback.record : '',
        crisisLevel: nextProps.feedback.crisis_level ? nextProps.feedback.crisis_level.toString() : '0',
      });
    }
  }

  handleChange(e) {
    const name = e.target.name;
    const value = e.target.value;
    const checked = e.target.type === 'checkbox' ? e.target.checked : value;
    switch (name) {
      case 'severity':
        this.setState((prevState) => {
          const { severity } = prevState;
          severity[Number(value)] = checked ? 1 : 0;
          return { severity };
        });
        break;
      case 'medical_diagnosis':
        this.setState((prevState) => {
          const { medicalDiagnosis } = prevState;
          medicalDiagnosis[Number(value)] = checked ? 1 : 0;
          return { medicalDiagnosis };
        });
        break;
      case 'crisis':
        this.setState((prevState) => {
          const { crisis } = prevState;
          crisis[Number(value)] = checked ? 1 : 0;
          return { crisis };
        });
        break;
      default:
        this.setState({ [name]: value });
        if (name === 'firstCategory') {
          this.setState({ secondCategory: '' });
        } else if (name === 'secondCategory') {
          if (value === 'A3') {
            this.setState((prevState) => {
              const { severity } = prevState;
              severity[0] = 1;
              return {
                severity,
                categoryShowTips: '',
              };
            });
          } else if (value === 'A4') {
            this.setState((prevState) => {
              const { severity } = prevState;
              severity[1] = 1;
              return {
                severity,
                categoryShowTips: '',
              };
            });
          } else {
            if (this.state.categoryShowCheckTips.includes(value)) {
              this.setState({ categoryShowTips: '请核查是否需要重点标记' });
            }
            if (this.state.categoryShowNeedTips.includes(value)) {
              this.setState({ categoryShowTips: '请选择合适的重点标记，否则不能够成功提交反馈表' });
            }
          }
        }
    }
  }

  handleSubmit() {
    this.setState({
      firstCategoryWarn: false,
      secondCategoryWarn: false,
      recordWarn: false,
    });
    if (this.state.firstCategory === '') {
      this.setState({ firstCategoryWarn: true });
      this.firstCategorySelect.focus();
      return;
    }
    if (this.state.secondCategory === '') {
      this.setState({ secondCategoryWarn: true });
      this.secondCategorySelect.focus();
      return;
    }
    if (this.state.record === '') {
      this.setState({ recordWarn: true });
      this.recordInput.focus();
      return;
    }
    if (!this.state.severity.includes(1) && !this.state.medicalDiagnosis.includes(1) && !this.state.crisis.includes(1)) {
      this.setState({ categoryShowTips: '请选择合适的重点标记，否则不能够成功提交反馈表' });
      this.props.showAlert('提交失败', '请选择合适的重点标记，否则不能够成功提交反馈表', '好的');
      return;
    }
    const payload = {
      reservation_id: this.state.reservation.id,
      source_id: this.state.reservation.source_id,
      category: this.state.secondCategory,
      severity: this.state.severity,
      medical_diagnosis: this.state.medicalDiagnosis,
      crisis: this.state.crisis,
      record: this.state.record,
      crisis_level: Number(this.state.crisisLevel),
    };
    this.props.handleSubmit(payload);
  }

  renderSecondCategories() {
    if (this.state.firstCategory === '') {
      return (
        <Select />
      );
    }
    return (
      <Select
        name="secondCategory"
        ref={(secondCategorySelect) => { this.secondCategorySelect = secondCategorySelect; }}
        value={this.state.secondCategory}
        onChange={this.handleChange}
      >
        {this.state.secondCategories && this.state.secondCategories[this.state.firstCategory] &&
          Object.keys(this.state.secondCategories[this.state.firstCategory]).map(name =>
            <option
              key={`second_category_option_${name}`}
              value={name}
            >
              {this.state.secondCategories[this.state.firstCategory][name]}
            </option>,
          )
        }
      </Select>
    );
  }

  renderEmphasis() {
    return (
      <div>
        <CellsTitle>严重程度</CellsTitle>
        {this.state.varSeverity && this.state.varSeverity.map((vs, i) =>
          <FormCell checkbox key={`checkbox-severity-${i}`}>
            <CellHeader>
              <Checkbox
                name="severity"
                value={i}
                checked={this.state.severity && this.state.severity.length > i && this.state.severity[i] === 1}
                onChange={this.handleChange}
              />
            </CellHeader>
            <CellBody>
              {vs}
            </CellBody>
          </FormCell>)
        }
        <CellsTitle>疑似或明确的医疗诊断</CellsTitle>
        {this.state.varMedicalDiagnosis && this.state.varMedicalDiagnosis.map((vmd, i) =>
          <FormCell checkbox key={`checkbox-medical-diagnosis-${i}`}>
            <CellHeader>
              <Checkbox
                name="medical_diagnosis"
                value={i}
                checked={this.state.medicalDiagnosis && this.state.medicalDiagnosis.length > i && this.state.medicalDiagnosis[i] === 1}
                onChange={this.handleChange}
              />
            </CellHeader>
            <CellBody>
              {vmd}
            </CellBody>
          </FormCell>)
        }
        <CellsTitle>危急情况</CellsTitle>
        {this.state.varCrisis && this.state.varCrisis.map((vc, i) =>
          <FormCell checkbox key={`checkbox-crisis-${i}`}>
            <CellHeader>
              <Checkbox
                name="crisis"
                value={i}
                checked={this.state.crisis && this.state.crisis.length > i && this.state.crisis[i] === 1}
                onChange={this.handleChange}
              />
            </CellHeader>
            <CellBody>
              {vc}
            </CellBody>
          </FormCell>)
        }
      </div>
    );
  }

  render() {
    return (
      <div>
        {this.state.reservation &&
          <CellsTitle>
            正在反馈：{this.state.reservation.start_time}-{this.state.reservation.end_time.slice(-5)} {this.state.reservation.teacher_fullname}
          </CellsTitle>
        }
        <Form checkbox>
          <FormCell warn={this.state.firstCategoryWarn} select selectPos="after">
            <CellHeader>
              <Label>评估分类<span style={{ color: 'red' }}>*</span></Label>
            </CellHeader>
            <CellBody>
              <Select
                name="firstCategory"
                ref={(firstCategorySelect) => { this.firstCategorySelect = firstCategorySelect; }}
                value={this.state.firstCategory}
                onChange={this.handleChange}
              >
                {this.state.firstCategories && Object.keys(this.state.firstCategories).map(name =>
                  <option
                    key={`first_category_option_${name}`}
                    value={name}
                  >
                    {this.state.firstCategories[name]}
                  </option>,
                )}
              </Select>
            </CellBody>
            {this.state.firstCategoryWarn &&
              <CellFooter style={{ marginRight: '25px' }}>
                <Icon value="warn" />
              </CellFooter>
            }
          </FormCell>
          <FormCell warn={this.state.secondCategoryWarn} select selectPos="after">
            <CellHeader>
              <Label>二级分类<span style={{ color: 'red' }}>*</span></Label>
            </CellHeader>
            <CellBody>
              {this.renderSecondCategories()}
            </CellBody>
            {this.state.secondCategoryWarn &&
              <CellFooter style={{ marginRight: '25px' }}>
                <Icon value="warn" />
              </CellFooter>
            }
          </FormCell>
          {this.state.categoryShowTips !== '' &&
            <CellsTitle>
              <span style={{ color: 'red' }}>{this.state.categoryShowTips}</span>
            </CellsTitle>
          }
          {this.renderEmphasis()}
          <CellsTitle>
            咨询记录<span style={{ color: 'red' }}>*</span>
          </CellsTitle>
          <FormCell warn={this.state.recordWarn}>
            <CellBody>
              <TextArea
                name="record"
                ref={(recordInput) => { this.recordInput = recordInput; }}
                placeholder="请输入咨询记录"
                rows="4"
                value={this.state.record}
                onChange={this.handleChange}
              />
            </CellBody>
            {this.state.recordWarn &&
              <CellFooter>
                <Icon value="warn" />
              </CellFooter>
            }
          </FormCell>
          <FormCell select selectPos="after">
            <CellHeader>
              <Label>是否危机个案</Label>
            </CellHeader>
            <CellBody>
              <Select
                name="crisisLevel"
                value={this.state.crisisLevel}
                onChange={this.handleChange}
              >
                <option value="0">否</option>
                <option value="3">三星</option>
                <option value="4">四星</option>
                <option value="5">五星</option>
              </Select>
            </CellBody>
          </FormCell>
        </Form>
        <ButtonArea direction="horizontal">
          <Button onClick={this.handleSubmit}>提交反馈</Button>
          <Button type="default" onClick={this.props.handleCancel}>返回</Button>
        </ButtonArea>
        <div style={{ height: '10px' }} />
      </div>
    );
  }
}

TeacherFeedbackForm.propTypes = propTypes;
