/**
 * Created by shudi on 2016/10/21.
 */
import React from "react";

export default class PageBottom extends React.Component{
    render() {
        let bottomStyle = {
            position: "fixed",
            bottom: "0px",
            borderTop: "1px solid #E5E5E5",
            padding: "5px 0 5px 0",
            width: "100%",
            zIndex: 10,
        };

        return (
            <div>
                <div style={{height: this.props.height}}></div>
                <div style={{...bottomStyle, ...this.props.style}}>
                    {
                        this.props.contents.map((content, index) => {
                            return (
                                <p key={`content-${index}`}>{content}</p>
                            );
                        })
                    }
                </div>
            </div>
        );
    }
}