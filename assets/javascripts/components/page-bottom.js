/**
 * Created by shudi on 2016/10/21.
 */
import React from "react";

let PageBottom = React.createClass({
    render() {
        let bottomStyle = {
            position: "fixed",
            bottom: "0px",
            borderTop: "1px solid #E5E5E5",
            padding: "10px 0 10px 0",
            width: "100%",
        };

        return (
            <div style={{...bottomStyle, ...this.props.style}}>
                {
                    this.props.contents.map(content => {
                        return (
                            <p>{content}</p>
                        );
                    })
                }
            </div>
        );
    },
});

export default PageBottom;