/**
 * Created by shudi on 2016/10/22.
 */
import React from 'react';
import ReactDOM from 'react-dom';
import {Link} from 'react-router';
import {CellsTitle, Form, FormCell, CellHeader, Label, CellBody, Input, CellFooter, Icon, ButtonArea, Button, Checkbox} from 'react-weui';
import 'weui';

export class LoginForm extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            username: '',
            usernameWarn: false,
            passwordWarn: false,
        };
        this.onSubmit = this.onSubmit.bind(this);
    }

    onSubmit() {
        let username = ReactDOM.findDOMNode(this.refs['usernameInput']).value;
        let password = ReactDOM.findDOMNode(this.refs['passwordInput']).value;
        this.setState({
            usernameWarn: false,
            passwordWarn: false,
        });
        if (username === '') {
            this.setState({usernameWarn: true});
            return;
        }
        if (password === '') {
            this.setState({passwordWarn: true});
            return;
        }
        this.props.onSubmit && this.props.onSubmit(username, password);
    }

    render() {
        return (
            <div>
                {
                    (this.props.titleTip && this.props.titleTip !== '') ?
                        <CellsTitle>{this.props.titleTip}</CellsTitle> : null
                }
                <Form>
                    <FormCell warn={this.state.usernameWarn}>
                        <CellHeader>
                            <Label>{this.props.usernameLabel}</Label>
                        </CellHeader>
                        <CellBody>
                            <Input ref="usernameInput"
                                   type={this.props.usernameType}
                                   placeholder={this.props.usernamePlaceholder}/>
                        </CellBody>
                        {
                            this.state.usernameWarn ?
                                <CellFooter>
                                    <Icon value="warn"/>
                                </CellFooter> : null
                        }
                    </FormCell>
                    <FormCell warn={this.state.passwordWarn}>
                        <CellHeader>
                            <Label>{this.props.passwordLabel}</Label>
                        </CellHeader>
                        <CellBody>
                            <Input ref="passwordInput"
                                   type={this.props.passwordType}
                                   placeholder={this.props.passwordPlaceholder}/>
                        </CellBody>
                        {
                            this.state.passwordWarn ?
                                <CellFooter>
                                    <Icon value="warn"/>
                                </CellFooter> : null
                        }
                    </FormCell>
                </Form>
                <ButtonArea direction="horizontal">
                    <Button onClick={this.onSubmit}>{this.props.submitText}</Button>
                    <Button type="default" onClick={this.props.onCancel}>{this.props.cancelText}</Button>
                </ButtonArea>
            </div>
        );
    }
};

export class RegisterForm extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            username: '',
            usernameWarn: false,
            passwordWarn: false,
            confirmPasswordWarn: false,
            protocolWarn: false,
        };
        this.onSubmit = this.onSubmit.bind(this);
    }

    onSubmit() {
        let username = ReactDOM.findDOMNode(this.refs['usernameInput']).value;
        let password = ReactDOM.findDOMNode(this.refs['passwordInput']).value;
        let confirmPassword = ReactDOM.findDOMNode(this.refs['confirmPasswordInput']).value;
        let protocolChecked = ReactDOM.findDOMNode(this.refs['protocolCheckbox']).children[0].checked;
        this.setState({
            usernameWarn: false,
            passwordWarn: false,
            confirmPasswordWarn: false,
            protocolWarn: false,
        });
        if (username === '') {
            this.setState({usernameWarn: true});
            return;
        }
        if (password === '') {
            this.setState({passwordWarn: true});
            return;
        }
        if (confirmPassword === '') {
            this.setState({confirmPasswordWarn: true});
            return;
        }
        if (!protocolChecked) {
            this.setState({protocolWarn: true});
            return;
        }
        if (password !== confirmPassword) {
            this.setState({
                passwordWarn: true,
                confirmPasswordWarn: true,
            });
            this.props.showAlert && this.props.showAlert('注册失败', '两次密码不一致，请重新输入');
            ReactDOM.findDOMNode(this.refs['passwordInput']).value = '';
            ReactDOM.findDOMNode(this.refs['confirmPasswordInput']).value = '';
            return;
        }
        this.props.onSubmit && this.props.onSubmit(username, password);
    }

    render() {
        return (
            <div>
                {
                    (this.props.titleTip && this.props.titleTip !== '') ?
                        <CellsTitle>{this.props.titleTip}</CellsTitle> : null
                }
                <Form checkbox={this.props.protocol ? true : false}>
                    <FormCell warn={this.state.usernameWarn}>
                        <CellHeader>
                            <Label>{this.props.usernameLabel}</Label>
                        </CellHeader>
                        <CellBody>
                            <Input ref="usernameInput"
                                   type={this.props.usernameType}
                                   placeholder={this.props.usernamePlaceholder}/>
                        </CellBody>
                        {
                            this.state.usernameWarn ?
                                <CellFooter>
                                    <Icon value="warn"/>
                                </CellFooter> : null
                        }
                    </FormCell>
                    <FormCell warn={this.state.passwordWarn}>
                        <CellHeader>
                            <Label>{this.props.passwordLabel}</Label>
                        </CellHeader>
                        <CellBody>
                            <Input ref="passwordInput"
                                   type={this.props.passwordType}
                                   placeholder={this.props.passwordPlaceholder}/>
                        </CellBody>
                        {
                            this.state.passwordWarn ?
                                <CellFooter>
                                    <Icon value="warn"/>
                                </CellFooter> : null
                        }
                    </FormCell>
                    <FormCell warn={this.state.confirmPasswordWarn}>
                        <CellHeader>
                            <Label>{this.props.confirmPasswordLabel}</Label>
                        </CellHeader>
                        <CellBody>
                            <Input ref="confirmPasswordInput"
                                   type={this.props.confirmPasswordType}
                                   placeholder={this.props.confirmPasswordPlaceholder}/>
                        </CellBody>
                        {
                            this.state.confirmPasswordWarn ?
                                <CellFooter>
                                    <Icon value="warn"/>
                                </CellFooter> : null
                        }
                    </FormCell>
                    {
                        this.props.protocol ?
                            <FormCell checkbox warn={this.state.protocolWarn}>
                                <CellHeader>
                                    <Checkbox ref="protocolCheckbox" defaultChecked/>
                                </CellHeader>
                                <CellBody>
                                    {this.props.protocolPrefix}
                                    <Link to={this.props.protocolLink}>{this.props.protocol}</Link>
                                    {this.props.protocolSuffix}
                                </CellBody>
                                {
                                    this.state.protocolWarn ?
                                        <CellFooter>
                                            <Icon value="warn"/>
                                        </CellFooter> : null
                                }
                            </FormCell> : null
                    }
                </Form>
                <ButtonArea direction="horizontal">
                    <Button onClick={this.onSubmit}>{this.props.submitText}</Button>
                    <Button type="default" onClick={this.props.onCancel}>{this.props.cancelText}</Button>
                </ButtonArea>
            </div>
        );
    }
};