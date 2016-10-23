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
            alertShow: false,
            alertTitle: '',
            alertMsg: '',
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
    render() {
        return (
            <div>
                {
                    (this.props.titleTip && this.props.titleTip !== '') ?
                        <CellsTitle>{this.props.titleTip}</CellsTitle> : null
                }
                <Form radio={false} checkbox={false}>
                    <FormCell key={`cell-${index}`} vcode={false} warn={false} radio={false}>
                        <CellHeader>
                            <Label>{name}</Label>
                        </CellHeader>
                        <CellBody>
                            <Input type={this.props.types[index]}
                                   placeholder={this.props.placeholders[index]}/>
                        </CellBody>
                    </FormCell>
                    {
                        this.props.protocolNames && this.props.protocolNames.length > 0 ?
                            this.props.protocolNames.map((protocolName, index) => {
                                return (
                                    <FormCell key={`protocol-${index}`} checkbox={true} vcode={false} warn={false}>
                                        <CellHeader>
                                            <Checkbox name="checkbox2" value="2" defaultChecked={true}/>
                                        </CellHeader>
                                        <CellBody>
                                            {this.props.protocolPrefix[index]}
                                            <Link to={this.props.protocolLinks[index]}>{protocolName}</Link>
                                            {this.props.protocolSuffix[index]}
                                        </CellBody>
                                    </FormCell>
                                )
                            }) : null
                    }
                </Form>
                <ButtonArea direction="horizontal">
                    <Button type="primary" size="normal" disabled={false} href="javascript:;" onClick={this.props.onSubmit}>{this.props.submitText}</Button>
                    <Button type="default" size="normal" disabled={false} href="javascript:;" onClick={this.props.onCancel}>{this.props.cancelText}</Button>
                </ButtonArea>
            </div>
        );
    }
}