/**
 * Created by shudi on 2016/10/20.
 */
import React from "react";
import {Link, Router} from "react-router";

import WeUI from 'react-weui';
const {Panel, PanelHeader, PanelBody, MediaBox, MediaBoxHeader, MediaBoxTitle, MediaBoxDescription} = WeUI;

let EntryPage = React.createClass({
    getInitialState() {
        return {
            studentEntry: "",
            teacherEntry: "",
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
                                <img src={window.assets["imgStudent"]}/>
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
                                <img src={window.assets["imgTeacher"]}/>
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