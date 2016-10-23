/**
 * Created by shudi on 2016/10/20.
 */
import React from 'react';
import ReactDOM from 'react-dom';
import {Link, Router} from 'react-router';
import {Panel, PanelHeader, PanelBody, MediaBox, MediaBoxHeader, MediaBoxBody, MediaBoxTitle, MediaBoxDescription} from 'react-weui';
import 'weui';

import PageBottom from '#coms/page-bottom';
import StudentEntryIcon from '#imgs/mobile/student.png';
import TeacherEntryIcon from '#imgs/mobile/teacher.png';

class EntryPage extends React.Component{
    render() {
        let studentEntryIcon = <img src={StudentEntryIcon}/>;
        let teacherEntryIcon = <img src={TeacherEntryIcon}/>;

        return (
            <Panel access>
                <PanelHeader style={{fontSize: "18px"}}>
                    学生心理发展指导中心预约系统
                </PanelHeader>
                <PanelBody>
                    <MediaBox type="appmsg" href="/m/student">
                        <MediaBoxHeader>{studentEntryIcon}</MediaBoxHeader>
                        <MediaBoxBody>
                            <MediaBoxTitle>我是学生</MediaBoxTitle>
                            <MediaBoxDescription>点击进入</MediaBoxDescription>
                        </MediaBoxBody>
                    </MediaBox>
                    <MediaBox type="appmsg" href="/m/teacher">
                        <MediaBoxHeader>{teacherEntryIcon}</MediaBoxHeader>
                        <MediaBoxBody>
                            <MediaBoxTitle>我是咨询师</MediaBoxTitle>
                            <MediaBoxDescription>点击进入</MediaBoxDescription>
                        </MediaBoxBody>
                    </MediaBox>
                </PanelBody>
                <PageBottom style={{color: "#999999", textAlign: "center"}}
                            contents={["清华大学学生心理发展指导中心", "联系方式：010-62782007"]}/>
            </Panel>
        );
    }
}

ReactDOM.render(
    React.createElement(EntryPage, null),
    document.getElementById('content')
);