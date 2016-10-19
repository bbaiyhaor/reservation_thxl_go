/**
 * Created by shudi on 2016/10/20.
 */
import React from "react";
import ReactDOM from 'react-dom';
import {Route, Router, hashHistory} from "react-router";

import EntryPage from "./user_mobile/page/EntryPage";

let resourceCount = 0;
let loaded = false;
let resourceReady = (timeout) => {
    resourceCount -= 1;
    if (resourceCount <= 0) {
        setTimeout(() => {
            if (!loaded) {
                ReactDOM.render((
                    <Router history={hashHistory}>
                        <Route name="home" path="/" component={EntryPage}/>
                        <Route name="entry" path="entry" component={EntryPage}/>
                        <Route path="*" component={EntryPage}/>
                    </Router>
                ), document.getElementById('user_mobile'));
                loaded = true;
            }
        }, timeout);
    }
};

let domReady = (timeout) => {
    let images = [];
    for (let prop in window.assets) {
        if (prop.indexOf("img") == 0) {
            images.push(window.assets[prop]);
        }
    }
    resourceCount = images.length;
    if (resourceCount <= 0) {
        resourceReady(timeout);
    } else {
        images.forEach(src => {
            let image = new Image();
            image.onload = () => resourceReady(timeout);
            image.onerror = () => resourceReady(timeout);
            image.src = src;
        });
        setTimeout(() => {
            resourceCount = 0;
            resourceReady(timeout);
        }, 2000);
    }
};

if (typeof document.onreadystatechange === "undefined") {
    window.onload = () => domReady(0);
} else {
    document.onreadystatechange = () => {
        if (document.readyState !== "complete") {
            return;
        }
        domReady(100);
    }
}