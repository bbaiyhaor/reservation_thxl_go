/**
 * Created by shudi on 2016/10/22.
 */
import React from 'react';
import {Panel, PanelHeader} from 'react-weui';
import 'weui';

import {UserSignForm} from '#coms/user-form';
import PageBottom from '#coms/page-bottom';

let StudentLoginPage = React.createClass({
    contextTypes: {
        router: React.PropTypes.object
    },

    onSubmit() {

    },

    onCancel() {
        this.context.router.push("register");
    },

    render() {
        return (
            <div>
                <Panel access={true}>
                    <PanelHeader style={{fontSize: "18px"}}>学生登录</PanelHeader>
                    <UserSignForm titleTip="请输入学号、密码登录（与info账号不同）"
                               names={["学号", "密码"]}
                               types={["tel", "number"]}
                               placeholders={["请输入学号", "请输入密码"]}
                               submitText="登录"
                               cancelText="没有账户"
                               onSubmit={this.onSubmit}
                               onCancel={this.onCancel}/>
                    <div style={{color: "#999999", padding: "10px 20px", textAlign: "center", fontSize: "13px"}}>
                        账号密码遇到任何问题，请在工作时间（周一至周五8:00-12:00；13:00-17:00）致电010-62782007。
                    </div>

                </Panel>
                <PageBottom style={{color: "#999999", textAlign: "center", backgroundColor: "white", fontSize: "14px"}}
                            contents={["清华大学学生心理发展指导中心", "联系方式：010-62782007"]}/>
            </div>
        );
    },
});

export default StudentLoginPage;