/**
 * Created by shudi on 2016/10/24.
 */
import React from 'react';
import {hashHistory} from 'react-router';
import {Panel, PanelHeader} from 'react-weui';
import 'weui';

import MakeReservationForm from '#forms/MakeReservationForm';
import PageBottom from '#coms/PageBottom';
import {AlertDialog, LoadingHud} from '#coms/Huds';
import {User, Application} from '#models/Models';

export default class StudentMakeReservationPage extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            student: null,
            reservation: null,
        };
        this.handleCancel = this.handleCancel.bind(this);
    }

    componentWillMount() {
        let reservationId = this.props.location.query['reservation_id'];
        if (reservationId === '') {
            hashHistory.push('reservation');
        }
        Application.viewReservationsByStudent(() => {
            if (!User.student || !Application.reservations || Application.reservations.length === 0) {
                hashHistory.push('reservation');
            }
            let i = 0;
            for (; i < Application.reservations.length; i++) {
                if (Application.reservations[i]['id'] === reservationId) {
                    this.setState({
                        student: User.student,
                        reservation: Application.reservations[i],
                    });
                    break;
                }
            }
            if (i === Application.reservations.length) {
                hashHistory.push('reservation');
            }
        }, () => {
            hashHistory.push('reservation');
        });
    }

    handleCancel() {
        hashHistory.goBack();
    }

    render() {
        return (
            <div>
                <Panel access>
                    <PanelHeader style={{fontSize: "18px"}}>学生信息登记表</PanelHeader>
                    <MakeReservationForm student={this.state.student}
                                         reservation={this.state.reservation}
                                         handleCancel={this.handleCancel}/>
                </Panel>
                <LoadingHud ref="loading"/>
                <AlertDialog ref="alert"/>
                <PageBottom styles={{color: "#999999", textAlign: "center", backgroundColor: "white", fontSize: "14px"}}
                            contents={["清华大学学生心理发展指导中心", "联系方式：010-62782007"]}
                            height="55px"/>
            </div>
        );
    }
}