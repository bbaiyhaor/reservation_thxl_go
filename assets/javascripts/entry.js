/**
 * Created by shudi on 2016/10/20.
 */
import React from "react";
import ReactDOM from "react-dom";
import {Link, Router} from "react-router";

import WeUI from 'react-weui';
const {Panel, PanelHeader, PanelBody, MediaBox, MediaBoxHeader, MediaBoxTitle, MediaBoxDescription} = WeUI;
import 'weui';

import StudentEntryIcon from '#imgs/mobile/student.png';
import TeacherEntryIcon from '#imgs/mobile/teacher.png';

let EntryPage = React.createClass({
    render() {
        return (
            <Panel access="true">
                <PanelHeader style={{fontSize: "4mm"}}>学生心理发展指导中心预约系统</PanelHeader>
                <PanelBody>
                    <MediaBox type="appmsg" href="">
                        <MediaBoxHeader>
                            <img src={StudentEntryIcon}/>
                        </MediaBoxHeader>
                        <PanelBody>
                            <MediaBoxTitle>我是学生</MediaBoxTitle>
                            <MediaBoxDescription>点击进入</MediaBoxDescription>
                        </PanelBody>
                    </MediaBox>
                    <MediaBox type="appmsg" href="">
                        <MediaBoxHeader>
                            <img src={TeacherEntryIcon}/>
                        </MediaBoxHeader>
                        <PanelBody>
                            <MediaBoxTitle>我是咨询师</MediaBoxTitle>
                            <MediaBoxDescription>点击进入</MediaBoxDescription>
                        </PanelBody>
                    </MediaBox>
                </PanelBody>
            </Panel>
        );
    },
});

ReactDOM.render(
    React.createElement(EntryPage, null),
    document.getElementById('content')
);