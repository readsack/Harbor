/* @refresh reload */
import { render } from "solid-js/web";
import App from "./App";
import Login from "./Login"
import Addr from "./Addr"
import "./index.css"
import { Route, Router } from "@solidjs/router";
import Org from "./Org";

const Root = (props) => {
    return (
        <>{props.children}</>
    )
}

render(() => (
    <Router root={Root}>
        <Route path="/" component={App}/>
        <Route path="/login" component={Login}/>
        <Route path="/addr" component={Addr}/>
        <Route path="/org" component={Org}/>
    </Router>
), document.getElementById("root"));
