/**
 * Created by shudi on 2016/10/23.
 */
import React from 'react';
import {hashHistory} from 'react-router';
import {Panel, PanelHeader, PanelBody, CellsTitle, MediaBox, MediaBoxDescription, MediaBoxInfo, MediaBoxInfoMeta, Button, Cells, Cell, CellBody, CellFooter} from 'react-weui';
import 'weui';

import UserLogoutButton from '#coms/LogoutButton';
import PageBottom from '#coms/PageBottom';
import {AlertDialog, ConfirmDialog, LoadingHud} from '#coms/Huds';
import {User, Application} from '#models/Models';

export default class StudentReservationListPage extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            student: null,
            reservations: null,
        };
        this.alert = this.alert.bind(this);
    }

    componentDidMount() {
        this.refs['loading'].show('正在加载中');
        Application.ViewReservationsByStudent(() => {
            setTimeout(() => {
                this.setState({
                    student: User.student,
                    reservations: Application.reservations,
                }, () => {
                    this.refs['loading'].hide();
                });
            }, 500);
        }, (status) => {
            this.refs['loading'].hide();
            this.refs['alert'].show('', status, '好的', () => {
                hashHistory.push('login');
            });
        });
    }

    alert(title, msg, label) {
        this.refs['alert'].show(title, msg, label);
    }

    render() {
        return (
            <div>
                <Panel access>
                    <PanelHeader style={{fontSize: "18px"}}>
                        {User.student && User.student['fullname'] ? User.student['fullname'] + '，' : ''}欢迎使用咨询预约系统
                        <div style={{height: "20px"}}>
                            <UserLogoutButton size="small" style={{float: "right"}} alert={this.alert}>退出登录</UserLogoutButton>
                        </div>
                        <CellsTitle>请根据您的需要选择相应咨询师和时间段进行预约</CellsTitle>
                    </PanelHeader>
                    <PanelBody>
                        <StudentReservationList reservations={this.state.reservations}/>
                    </PanelBody>
                </Panel>
                <LoadingHud ref="loading"/>
                <AlertDialog ref="alert"/>
                <PageBottom style={{color: "#999999", textAlign: "center", backgroundColor: "white", fontSize: "14px"}}
                            contents={["清华大学学生心理发展指导中心", "联系方式：010-62782007"]}
                            height="55px"/>
            </div>
        )
    }
}

class StudentReservationList extends React.Component {
    constructor(props) {
        super(props);
        this.makeReservation = this.makeReservation.bind(this);
        this.feedback = this.feedback.bind(this);
    }

    makeReservation(reservation) {
        this.refs['confirm'].show('',
            '确定预约后请准确填写个人信息，方便心理咨询中心老师与你取得联系。',
            '暂不预约', '立即预约', null, () => {
                hashHistory.push(`/reservation/make?reservation_id=${reservation['id']}`);
            });
    }

    feedback(reservation) {
        console.log(reservation);
    }

    render() {
        return (
            <div>
                <MediaBox type="small_appmsg">
                    <Cells access>
                        {
                            this.props.reservations && this.props.reservations.map((reservation, index) => {
                                return (
                                    <Cell key={`reservation-cell-${reservation.id}`}>
                                        <CellBody>
                                            <p style={{fontSize: "14px"}}>{reservation['start_time']} - {reservation['end_time'].slice(-5)}　{reservation['teacher_fullname']}</p>
                                        </CellBody>
                                        {
                                            reservation['status'] === 1 ?
                                                <Button size="small" onClick={(e) => {
                                                    this.makeReservation(reservation);
                                                    e.stopPropagation();
                                                    e.preventDefault();
                                                }}>预约</Button> :
                                                (
                                                    reservation['status'] === 2 ?
                                                        <Button size="small" type="default" disabled>已预约</Button> :
                                                        (
                                                            reservation['status'] === 3 ?
                                                                <Button size="small" type="warn" onClick={(e) => {
                                                                    this.feedback(reservation);
                                                                    e.stopPropagation();
                                                                    e.preventDefault();
                                                                }}>反馈</Button> : null
                                                        )
                                                )
                                        }
                                    </Cell>
                                );
                            })
                        }
                    </Cells>
                </MediaBox>
                <ConfirmDialog ref="confirm"/>
            </div>
        );
    }
}