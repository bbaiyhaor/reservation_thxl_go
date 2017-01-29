/**
 * Created by shudi on 2016/11/4.
 */
import React, { PropTypes } from 'react';
import { Form, FormCell, CellsTitle, CellHeader, CellBody, CellFooter, Label, Select, Checkbox, TextArea, Icon, ButtonArea, Button } from '#react-weui';
import 'weui';

const propTypes = {
  handleSubmit: PropTypes.func.isRequired,
  handleCancel: PropTypes.func.isRequired,
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
      severity_0: false,
      severity_1: false,
      severity_2: false,
      severity_3: false,
      severity_4: false,
      medicalDiagnosis_0: false,
      medicalDiagnosis_1: false,
      medicalDiagnosis_2: false,
      medicalDiagnosis_3: false,
      medicalDiagnosis_4: false,
      medicalDiagnosis_5: false,
      medicalDiagnosis_6: false,
      medicalDiagnosis_7: false,
      medicalDiagnosis_8: false,
      medicalDiagnosis_9: false,
      medicalDiagnosis_10: false,
      crisis_0: false,
      crisis_1: false,
      crisis_2: false,
      crisis_3: false,
      record: '',
      crisisLevel: '0',
      firstCategoryWarn: false,
      secondCategoryWarn: false,
      recordWarn: false,
    };
    this.handleChange = this.handleChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  componentWillReceiveProps(nextProps) {
    this.setState({
      reservation: nextProps.reservation,
      feedback: nextProps.feedback,
      firstCategories: nextProps.feedback.first_category,
      secondCategories: nextProps.feedback.second_category,
    });
    if (nextProps.feedback) {
      this.setState({
        firstCategory: nextProps.feedback.category ? nextProps.feedback.category.substring(0, 1) : '',
        secondCategory: nextProps.feedback.category ? nextProps.feedback.category : '',
        severity_0: nextProps.feedback.severity && nextProps.feedback.severity[0] === 1,
        severity_1: nextProps.feedback.severity && nextProps.feedback.severity[1] === 1,
        severity_2: nextProps.feedback.severity && nextProps.feedback.severity[2] === 1,
        severity_3: nextProps.feedback.severity && nextProps.feedback.severity[3] === 1,
        severity_4: nextProps.feedback.severity && nextProps.feedback.severity[4] === 1,
        medicalDiagnosis_0: nextProps.feedback.medical_diagnosis && nextProps.feedback.medical_diagnosis[0] === 1,
        medicalDiagnosis_1: nextProps.feedback.medical_diagnosis && nextProps.feedback.medical_diagnosis[1] === 1,
        medicalDiagnosis_2: nextProps.feedback.medical_diagnosis && nextProps.feedback.medical_diagnosis[2] === 1,
        medicalDiagnosis_3: nextProps.feedback.medical_diagnosis && nextProps.feedback.medical_diagnosis[3] === 1,
        medicalDiagnosis_4: nextProps.feedback.medical_diagnosis && nextProps.feedback.medical_diagnosis[4] === 1,
        medicalDiagnosis_5: nextProps.feedback.medical_diagnosis && nextProps.feedback.medical_diagnosis[5] === 1,
        medicalDiagnosis_6: nextProps.feedback.medical_diagnosis && nextProps.feedback.medical_diagnosis[6] === 1,
        medicalDiagnosis_7: nextProps.feedback.medical_diagnosis && nextProps.feedback.medical_diagnosis[7] === 1,
        medicalDiagnosis_8: nextProps.feedback.medical_diagnosis && nextProps.feedback.medical_diagnosis[8] === 1,
        medicalDiagnosis_9: nextProps.feedback.medical_diagnosis && nextProps.feedback.medical_diagnosis[9] === 1,
        medicalDiagnosis_10: nextProps.feedback.medical_diagnosis && nextProps.feedback.medical_diagnosis[10] === 1,
        crisis_0: nextProps.feedback.crisis && nextProps.feedback.crisis[0] === 1,
        crisis_1: nextProps.feedback.crisis && nextProps.feedback.crisis[1] === 1,
        crisis_2: nextProps.feedback.crisis && nextProps.feedback.crisis[2] === 1,
        crisis_3: nextProps.feedback.crisis && nextProps.feedback.crisis[3] === 1,
        record: nextProps.feedback.record ? nextProps.feedback.record : '',
        crisisLevel: nextProps.feedback.crisis_level ? nextProps.feedback.crisis_level.toString() : '0',
      });
    }
  }

  handleChange(e, name) {
    const value = e.target.value;
    if (name && name !== '') {
      this.setState({ [name]: value });
      if (name === 'firstCategory') {
        this.setState({ secondCategory: '' });
      }
    } else {
      this.setState(prevState => ({
        [value]: !prevState[value],
      }));
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
    const severity = [];
    for (let i = 0; i < 5; i += 1) {
      severity.push(this.state[`severity_${i}`] ? 1 : 0);
    }
    const medicalDiagnosis = [];
    for (let i = 0; i < 11; i += 1) {
      medicalDiagnosis.push(this.state[`medicalDiagnosis_${i}`] ? 1 : 0);
    }
    const crisis = [];
    for (let i = 0; i < 4; i += 1) {
      crisis.push(this.state[`crisis_${i}`] ? 1 : 0);
    }
    const payload = {
      reservation_id: this.state.reservation.id,
      source_id: this.state.reservation.source_id,
      category: this.state.secondCategory,
      severity,
      medical_diagnosis: medicalDiagnosis,
      crisis,
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
        ref={(secondCategorySelect) => { this.secondCategorySelect = secondCategorySelect; }}
        value={this.state.secondCategory}
        onChange={(e) => { this.handleChange(e, 'secondCategory'); }}
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
        <FormCell checkbox>
          <CellHeader>
            <Checkbox
              value="severity_0"
              checked={this.state.severity_0}
              onChange={this.handleChange}
            />
          </CellHeader>
          <CellBody>
            缓考
          </CellBody>
        </FormCell>
        <FormCell checkbox>
          <CellHeader>
            <Checkbox
              value="severity_1"
              checked={this.state.severity_1}
              onChange={this.handleChange}
            />
          </CellHeader>
          <CellBody>
            休学复学
          </CellBody>
        </FormCell>
        <FormCell checkbox>
          <CellHeader>
            <Checkbox
              value="severity_2"
              checked={this.state.severity_2}
              onChange={this.handleChange}
            />
          </CellHeader>
          <CellBody>
            家属陪读
          </CellBody>
        </FormCell>
        <FormCell checkbox>
          <CellHeader>
            <Checkbox
              value="severity_3"
              checked={this.state.severity_3}
              onChange={this.handleChange}
            />
          </CellHeader>
          <CellBody>
            家属不知情
          </CellBody>
        </FormCell>
        <FormCell checkbox>
          <CellHeader>
            <Checkbox
              value="severity_4"
              checked={this.state.severity_4}
              onChange={this.handleChange}
            />
          </CellHeader>
          <CellBody>
            任何其他需要知会院系关注的原因
          </CellBody>
        </FormCell>
        <CellsTitle>疑似或明确的医疗诊断</CellsTitle>
        <FormCell checkbox>
          <CellHeader>
            <Checkbox
              value="medicalDiagnosis_0"
              checked={this.state.medicalDiagnosis_0}
              onChange={this.handleChange}
            />
          </CellHeader>
          <CellBody>
            服药
          </CellBody>
        </FormCell>
        <FormCell checkbox>
          <CellHeader>
            <Checkbox
              value="medicalDiagnosis_1"
              checked={this.state.medicalDiagnosis_1}
              onChange={this.handleChange}
            />
          </CellHeader>
          <CellBody>
            精神分裂
          </CellBody>
        </FormCell>
        <FormCell checkbox>
          <CellHeader>
            <Checkbox
              value="medicalDiagnosis_2"
              checked={this.state.medicalDiagnosis_2}
              onChange={this.handleChange}
            />
          </CellHeader>
          <CellBody>
            双相情感障碍
          </CellBody>
        </FormCell>
        <FormCell checkbox>
          <CellHeader>
            <Checkbox
              value="medicalDiagnosis_3"
              checked={this.state.medicalDiagnosis_3}
              onChange={this.handleChange}
            />
          </CellHeader>
          <CellBody>
            焦虑症（状态）
          </CellBody>
        </FormCell>
        <FormCell checkbox>
          <CellHeader>
            <Checkbox
              value="medicalDiagnosis_4"
              checked={this.state.medicalDiagnosis_4}
              onChange={this.handleChange}
            />
          </CellHeader>
          <CellBody>
            抑郁症（状态）
          </CellBody>
        </FormCell>
        <FormCell checkbox>
          <CellHeader>
            <Checkbox
              value="medicalDiagnosis_5"
              checked={this.state.medicalDiagnosis_5}
              onChange={this.handleChange}
            />
          </CellHeader>
          <CellBody>
            强迫症
          </CellBody>
        </FormCell>
        <FormCell checkbox>
          <CellHeader>
            <Checkbox
              value="medicalDiagnosis_6"
              checked={this.state.medicalDiagnosis_6}
              onChange={this.handleChange}
            />
          </CellHeader>
          <CellBody>
            进食障碍
          </CellBody>
        </FormCell>
        <FormCell checkbox>
          <CellHeader>
            <Checkbox
              value="medicalDiagnosis_7"
              checked={this.state.medicalDiagnosis_7}
              onChange={this.handleChange}
            />
          </CellHeader>
          <CellBody>
            失眠
          </CellBody>
        </FormCell>
        <FormCell checkbox>
          <CellHeader>
            <Checkbox
              value="medicalDiagnosis_8"
              checked={this.state.medicalDiagnosis_8}
              onChange={this.handleChange}
            />
          </CellHeader>
          <CellBody>
            其他精神症状
          </CellBody>
        </FormCell>
        <FormCell checkbox>
          <CellHeader>
            <Checkbox
              value="medicalDiagnosis_9"
              checked={this.state.medicalDiagnosis_9}
              onChange={this.handleChange}
            />
          </CellHeader>
          <CellBody>
            躯体疾病
          </CellBody>
        </FormCell>
        <FormCell checkbox>
          <CellHeader>
            <Checkbox
              value="medicalDiagnosis_10"
              checked={this.state.medicalDiagnosis_10}
              onChange={this.handleChange}
            />
          </CellHeader>
          <CellBody>
            不遵医嘱
          </CellBody>
        </FormCell>
        <CellsTitle>危急情况</CellsTitle>
        <FormCell checkbox>
          <CellHeader>
            <Checkbox
              value="crisis_0"
              checked={this.state.crisis_0}
              onChange={this.handleChange}
            />
          </CellHeader>
          <CellBody>
            自伤
          </CellBody>
        </FormCell>
        <FormCell checkbox>
          <CellHeader>
            <Checkbox
              value="crisis_1"
              checked={this.state.crisis_1}
              onChange={this.handleChange}
            />
          </CellHeader>
          <CellBody>
            伤害他人
          </CellBody>
        </FormCell>
        <FormCell checkbox>
          <CellHeader>
            <Checkbox
              value="crisis_2"
              checked={this.state.crisis_2}
              onChange={this.handleChange}
            />
          </CellHeader>
          <CellBody>
            自杀念头
          </CellBody>
        </FormCell>
        <FormCell checkbox>
          <CellHeader>
            <Checkbox
              value="crisis_3"
              checked={this.state.crisis_3}
              onChange={this.handleChange}
            />
          </CellHeader>
          <CellBody>
            自杀未遂
          </CellBody>
        </FormCell>
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
                ref={(firstCategorySelect) => { this.firstCategorySelect = firstCategorySelect; }}
                value={this.state.firstCategory}
                onChange={(e) => { this.handleChange(e, 'firstCategory'); }}
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
          {this.renderEmphasis()}
          <CellsTitle>
            咨询记录<span style={{ color: 'red' }}>*</span>
          </CellsTitle>
          <FormCell warn={this.state.recordWarn}>
            <CellBody>
              <TextArea
                ref={(recordInput) => { this.recordInput = recordInput; }}
                placeholder="请输入咨询记录"
                rows="4"
                value={this.state.record}
                onChange={(e) => { this.handleChange(e, 'record'); }}
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
                value={this.state.crisisLevel}
                onChange={(e) => { this.handleChange(e, 'crisisLevel'); }}
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
