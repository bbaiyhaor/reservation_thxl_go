/**
 * Created by shudi on 2016/10/22.
 */
import React from 'react';
import {Link} from 'react-router';
import {CellsTitle, Form, FormCell, CellHeader, Label, CellBody, Input, ButtonArea, Button, Checkbox} from 'react-weui';
import 'weui';

export let UserSignForm = React.createClass({
    render() {
        return (
            <div>
                {
                    (this.props.titleTip && this.props.titleTip !== '') ?
                        <CellsTitle>{this.props.titleTip}</CellsTitle> : null
                }
                <Form radio={false} checkbox={true}>
                    {
                        this.props.names && this.props.names.length > 0 ?
                            this.props.names.map((name, index) => {
                                return (
                                    <FormCell key={`cell-${index}`} vcode={false} warn={false} radio={false}>
                                        <CellHeader>
                                            <Label>{name}</Label>
                                        </CellHeader>
                                        <CellBody>
                                            <Input type={this.props.types[index]}
                                                   placeholder={this.props.placeholders[index]}/>
                                        </CellBody>
                                    </FormCell>
                                );
                            }) : null
                    }
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
    },
});