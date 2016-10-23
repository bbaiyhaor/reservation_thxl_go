/**
 * Created by shudi on 2016/10/23.
 */
import React from 'react';
import {Dialog, Toast} from 'react-weui';
const {Alert} = Dialog;
import 'weui';

export class AlertDialog extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            alertTitle: '',
            alertShow: false,
            alertMsg: '',
            alertButtons: [
                {
                    label: '',
                    onClick: this.hide.bind(this),
                }
            ],
        };
    }

    show(title, msg, label) {
        this.setState({
            alertTitle: title,
            alertShow: true,
            alertMsg: msg,
            alertButtons: [
                {
                    label: label,
                    onClick: this.hide.bind(this),
                }
            ],
        });
    }

    hide() {
        this.setState({
            alertTitle: '',
            alertShow: false,
            alertMsg: '',
        });
    }

    render() {
        return (
            <Alert title={this.state.alertTitle} buttons={this.state.alertButtons} show={this.state.alertShow}>
                {this.state.alertMsg}
            </Alert>
        );
    }
}

export class LoadingHud extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            showLoading: false,
            loadingMsg: '',
        };
    }

    show(msg) {
        this.setState({
            showLoading: true,
            loadingMsg: msg,
        });
    }

    hide() {
        this.setState({
            showLoading: false,
            loadingMsg: '',
        });
    }

    render() {
        return (
            <Toast icon="loading" show={this.state.showLoading}>
                {this.state.loadingMsg}
            </Toast>
        );
    }
}