/**
 * Created by shudi on 2016/10/24.
 */
import React from 'react';
import ReactDOM from 'react-dom';
import {Link} from 'react-router';
import {Form, FormCell, CellHeader, CellBody, CellFooter, CellsTitle, Label, Input, Icon, Select, TextArea, ButtonArea, Button} from 'react-weui';
import 'weui';

export default class ReservationMakeForm extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            reservation: null,
            usernameWarn: false,
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
        this.setInputValue = this.setInputValue.bind(this);
    }

    componentWillReceiveProps(nextProps) {
        nextProps['student'] && this.setInputValue(nextProps['student']);
        nextProps['reservation'] && this.setState({reservation: nextProps['reservation']});
    }

    componentDidMount() {
        this.setInputValue(this.props.student);
    }
    
    setInputValue(student) {
        if (student) {
            ReactDOM.findDOMNode(this.refs['usernameInput']).value = student['fullname'] ? student['fullname'] : '';
            ReactDOM.findDOMNode(this.refs['genderInput']).value = student['gender'] ? student['gender'] : '';
            ReactDOM.findDOMNode(this.refs['birthdayInput']).value = student['birthday'] ? student['birthday'] : '';
            ReactDOM.findDOMNode(this.refs['schoolInput']).value = student['school'] ? student['school'] : '';
            ReactDOM.findDOMNode(this.refs['gradeInput']).value = student['grade'] ? student['grade'] : '';
            ReactDOM.findDOMNode(this.refs['currentAddressInput']).value = student['current_address'] ? student['current_address'] : '';
            ReactDOM.findDOMNode(this.refs['familyAddressInput']).value = student['family_address'] ? student['family_address'] : '';
            ReactDOM.findDOMNode(this.refs['mobileInput']).value = student['mobile'] ? student['mobile'] : '';
            ReactDOM.findDOMNode(this.refs['emailInput']).value = student['email'] ? student['email'] : '';
            ReactDOM.findDOMNode(this.refs['experienceTimeInput']).value = student['experience_time'] ? student['experience_time'] : '';
            ReactDOM.findDOMNode(this.refs['experienceLocationInput']).value = student['experience_location'] ? student['experience_location'] : '';
            ReactDOM.findDOMNode(this.refs['experienceTeacherInput']).value = student['experience_teacher'] ? student['experience_teacher'] : '';
            ReactDOM.findDOMNode(this.refs['fatherAgeInput']).value = student['father_age'] ? student['father_age'] : '';
            ReactDOM.findDOMNode(this.refs['fatherJobInput']).value = student['father_job'] ? student['father_job'] : '';
            ReactDOM.findDOMNode(this.refs['fatherEduInput']).value = student['father_edu'] ? student['father_edu'] : '';
            ReactDOM.findDOMNode(this.refs['motherAgeInput']).value = student['mother_age'] ? student['mother_age'] : '';
            ReactDOM.findDOMNode(this.refs['motherJobInput']).value = student['mother_job'] ? student['mother_job'] : '';
            ReactDOM.findDOMNode(this.refs['motherEduInput']).value = student['mother_edu'] ? student['mother_edu'] : '';
            ReactDOM.findDOMNode(this.refs['parentMarriageInput']).value = student['parent_marriage'] ? student['parent_marriage'] : '';
            ReactDOM.findDOMNode(this.refs['significantInput']).children[0].value = student['significant'] ? student['significant'] : '';
            ReactDOM.findDOMNode(this.refs['problemInput']).children[0].value = student['problem'] ? student['problem'] : '';
        }
    }

    render() {
        return (
            <div>
                {
                    this.state.reservation ?
                        <CellsTitle>正在预约：{this.state.reservation['start_time']}-{this.state.reservation['end_time'].slice(-5)} {this.state.reservation['teacher_fullname']}</CellsTitle> : null
                }
                <Form>
                    <FormCell warn={this.state.usernameWarn}>
                        <CellHeader>
                            <Label>姓名<span style={{color: "red"}}>*</span></Label>
                        </CellHeader>
                        <CellBody>
                            <Input ref="usernameInput"
                                   type="input"
                                   placeholder="请输入姓名"/>
                        </CellBody>
                        {
                            this.state.usernameWarn ?
                                <CellFooter>
                                    <Icon value="warn"/>
                                </CellFooter> : null
                        }
                    </FormCell>
                    <FormCell warn={this.state.genderWarn} select selectPos="after">
                        <CellHeader>
                            <Label>性别<span style={{color: "red"}}>*</span></Label>
                        </CellHeader>
                        <CellBody>
                            <Select ref="genderInput">
                                <option value="">请选择</option>
                                <option value="男">男</option>
                                <option value="女">女</option>
                            </Select>
                        </CellBody>
                        {
                            this.state.genderWarn ?
                                <CellFooter style={{marginRight: "25px"}}>
                                    <Icon value="warn"/>
                                </CellFooter> : null
                        }
                    </FormCell>
                    <FormCell warn={this.state.birthdayWarn}>
                        <CellHeader>
                            <Label>出生日期<span style={{color: "red"}}>*</span></Label>
                        </CellHeader>
                        <CellBody>
                            <Input ref="birthdayInput"
                                   type="input"
                                   placeholder="请输入出生日期"/>
                        </CellBody>
                        {
                            this.state.birthdayWarn ?
                                <CellFooter>
                                    <Icon value="warn"/>
                                </CellFooter> : null
                        }
                    </FormCell>
                    <FormCell warn={this.state.schoolWarn}>
                        <CellHeader>
                            <Label>院系<span style={{color: "red"}}>*</span></Label>
                        </CellHeader>
                        <CellBody>
                            <Input ref="schoolInput"
                                   type="input"
                                   placeholder="请输入院系"/>
                        </CellBody>
                        {
                            this.state.schoolWarn ?
                                <CellFooter>
                                    <Icon value="warn"/>
                                </CellFooter> : null
                        }
                    </FormCell>
                    <FormCell warn={this.state.gradeWarn}>
                        <CellHeader>
                            <Label>年级<span style={{color: "red"}}>*</span></Label>
                        </CellHeader>
                        <CellBody>
                            <Input ref="gradeInput"
                                   type="input"
                                   placeholder="请输入年级"/>
                        </CellBody>
                        {
                            this.state.gradeWarn ?
                                <CellFooter>
                                    <Icon value="warn"/>
                                </CellFooter> : null
                        }
                    </FormCell>
                    <FormCell warn={this.state.currentAddressWarn}>
                        <CellHeader>
                            <Label>现在住址<span style={{color: "red"}}>*</span></Label>
                        </CellHeader>
                        <CellBody>
                            <Input ref="currentAddressInput"
                                   type="input"
                                   placeholder="请输入现在住址"/>
                        </CellBody>
                        {
                            this.state.currentAddressWarn ?
                                <CellFooter>
                                    <Icon value="warn"/>
                                </CellFooter> : null
                        }
                    </FormCell>
                    <FormCell warn={this.state.familyAddressWarn}>
                        <CellHeader>
                            <Label>家庭住址<span style={{color: "red"}}>*</span></Label>
                        </CellHeader>
                        <CellBody>
                            <Input ref="familyAddressInput"
                                   type="input"
                                   placeholder="请输入家庭住址"/>
                        </CellBody>
                        {
                            this.state.familyAddressWarn ?
                                <CellFooter>
                                    <Icon value="warn"/>
                                </CellFooter> : null
                        }
                    </FormCell>
                    <FormCell warn={this.state.mobileWarn}>
                        <CellHeader>
                            <Label>联系电话<span style={{color: "red"}}>*</span></Label>
                        </CellHeader>
                        <CellBody>
                            <Input ref="mobileInput"
                                   type="tel"
                                   placeholder="请输入联系电话"/>
                        </CellBody>
                        {
                            this.state.mobileWarn ?
                                <CellFooter>
                                    <Icon value="warn"/>
                                </CellFooter> : null
                        }
                    </FormCell>
                    <FormCell warn={this.state.emailWarn}>
                        <CellHeader>
                            <Label>邮箱<span style={{color: "red"}}>*</span></Label>
                        </CellHeader>
                        <CellBody>
                            <Input ref="emailInput"
                                   type="input"
                                   placeholder="请输入邮箱"/>
                        </CellBody>
                        {
                            this.state.emailWarn ?
                                <CellFooter>
                                    <Icon value="warn"/>
                                </CellFooter> : null
                        }
                    </FormCell>
                    <CellsTitle>过往咨询经历</CellsTitle>
                    <FormCell>
                        <CellHeader>
                            <Label>时间</Label>
                        </CellHeader>
                        <CellBody>
                            <Input ref="experienceTimeInput"
                                   type="input"
                                   placeholder="请输入咨询时间"/>
                        </CellBody>
                    </FormCell>
                    <FormCell>
                        <CellHeader>
                            <Label>地点</Label>
                        </CellHeader>
                        <CellBody>
                            <Input ref="experienceLocationInput"
                                   type="input"
                                   placeholder="请输入咨询地点"/>
                        </CellBody>
                    </FormCell>
                    <FormCell>
                        <CellHeader>
                            <Label>咨询师姓名</Label>
                        </CellHeader>
                        <CellBody>
                            <Input ref="experienceTeacherInput"
                                   type="input"
                                   placeholder="请输入咨询师姓名"/>
                        </CellBody>
                    </FormCell>
                    <CellsTitle>家庭情况</CellsTitle>
                    <FormCell>
                        <CellHeader>
                            <Label>父亲年龄</Label>
                        </CellHeader>
                        <CellBody>
                            <Input ref="fatherAgeInput"
                                   type="input"
                                   placeholder="请输入父亲年龄"/>
                        </CellBody>
                    </FormCell>
                    <FormCell>
                        <CellHeader>
                            <Label>父亲职业</Label>
                        </CellHeader>
                        <CellBody>
                            <Input ref="fatherJobInput"
                                   type="input"
                                   placeholder="请输入父亲职业"/>
                        </CellBody>
                    </FormCell>
                    <FormCell>
                        <CellHeader>
                            <Label>父亲学历</Label>
                        </CellHeader>
                        <CellBody>
                            <Input ref="fatherEduInput"
                                   type="input"
                                   placeholder="请输入父亲学历"/>
                        </CellBody>
                    </FormCell>
                    <FormCell>
                        <CellHeader>
                            <Label>母亲年龄</Label>
                        </CellHeader>
                        <CellBody>
                            <Input ref="motherAgeInput"
                                   type="input"
                                   placeholder="请输入母亲年龄"/>
                        </CellBody>
                    </FormCell>
                    <FormCell>
                        <CellHeader>
                            <Label>母亲职业</Label>
                        </CellHeader>
                        <CellBody>
                            <Input ref="motherJobInput"
                                   type="input"
                                   placeholder="请输入母亲职业"/>
                        </CellBody>
                    </FormCell>
                    <FormCell>
                        <CellHeader>
                            <Label>母亲学历</Label>
                        </CellHeader>
                        <CellBody>
                            <Input ref="motherEduInput"
                                   type="input"
                                   placeholder="请输入母亲学历"/>
                        </CellBody>
                    </FormCell>
                    <FormCell select selectPos="after">
                        <CellHeader>
                            <Label>父母婚姻状况</Label>
                        </CellHeader>
                        <CellBody>
                            <Select ref="parentMarriageInput">
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
                            <TextArea ref="significantInput"
                                      placeholder="请输入"
                                      rows="3"
                                      maxlength="300"/>
                        </CellBody>
                    </FormCell>
                    <CellsTitle>
                        你现在需要接受帮助的主要问题是什么？
                    </CellsTitle>
                    <FormCell warn={this.state.problemWarn}>
                        <CellBody>
                            <TextArea ref="problemInput"
                                      placeholder="请输入"
                                      rows="3"
                                      maxlength="300"/>
                        </CellBody>
                        {
                            this.state.problemWarn ?
                                <CellFooter>
                                    <Icon value="warn"/>
                                </CellFooter> : null
                        }
                    </FormCell>
                </Form>
                <ButtonArea direction="horizontal">
                    <Button onClick={this.onSubmit}>确定预约</Button>
                    <Button type="default" onClick={this.props.onCancel}>取消预约</Button>
                </ButtonArea>
                <div style={{height: "10px"}}></div>
            </div>
        );
    }
}