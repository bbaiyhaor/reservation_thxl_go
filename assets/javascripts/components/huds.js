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
            alertShow: false,
            alert: {
                title: '',
                msg: '',
                buttons: [
                    {
                        label: '',
                        onClick: this.hide.bind(this),
                    }
                ],
            },
        };
        this.show = this.show.bind(this);
        this.hide = this.hide.bind(this);
    }

    show(title, msg, label, click) {
        this.setState({
            alertShow: true,
            alert: {
                title: title,
                msg: msg,
                buttons: [
                    {
                        label: label,
                        onClick: click || this.hide.bind(this),
                    }
                ],
            },
        });
    }

    hide() {
        this.setState({
            alertShow: false,
            alert: {
                title: '',
                msg: '',
            },
        });
    }

    render() {
        return (
            <Alert title={this.state.alert.title} buttons={this.state.alert.buttons} show={this.state.alertShow}>
                {this.state.alert.msg}
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
        this.show = this.show.bind(this);
        this.hide = this.hide.bind(this);
    }

    show(msg, duration) {
        this.setState({
            showLoading: true,
            loadingMsg: msg,
        });
        if (duration && duration > 0) {
            setTimeout(() => {
                this.hide();
            }, duration);
        }
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