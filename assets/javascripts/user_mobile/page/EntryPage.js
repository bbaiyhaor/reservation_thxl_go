/**
 * Created by shudi on 2016/10/20.
 */
import React from "react";
import {Link, Router} from "react-router";

import WeUI from 'react-weui';
const {Panel, PanelHeader, PanelBody, MediaBox, MediaBoxHeader, MediaBoxTitle, MediaBoxDescription} = WeUI;
import 'weui';

import StudentImg from '../../../images/user_mobile/student.png';
import TeacherImg from '../../../images/user_mobile/teacher.png';

let EntryPage = React.createClass({
    getInitialState() {
        return {
            studentEntry: "entry",
            teacherEntry: "entry",
        };
    },

    render() {
        return (
            <Panel access="true">
                <PanelHeader style={{fontSize: "4mm"}}>学生心理发展指导中心预约系统</PanelHeader>
                <PanelBody>
                    <Link to={this.state.studentEntry}>
                        <MediaBox type="appmsg">
                            <MediaBoxHeader>
                                <img src={StudentImg}/>
                            </MediaBoxHeader>
                            <PanelBody>
                                <MediaBoxTitle>我是学生</MediaBoxTitle>
                                <MediaBoxDescription>点击进入</MediaBoxDescription>
                            </PanelBody>
                        </MediaBox>
                    </Link>
                    <Link to={this.state.teacherEntry}>
                        <MediaBox type="appmsg">
                            <MediaBoxHeader>
                                <img src={TeacherImg}/>
                            </MediaBoxHeader>
                            <PanelBody>
                                <MediaBoxTitle>我是咨询师</MediaBoxTitle>
                                <MediaBoxDescription>点击进入</MediaBoxDescription>
                            </PanelBody>
                        </MediaBox>
                    </Link>
                </PanelBody>
            </Panel>
        );
    },
});

export default EntryPage;