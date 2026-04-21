/* @refresh reload */
import { render } from "solid-js/web";
import App from "./App";
import "./index.css"
import { Route, Router } from "@solidjs/router";

const Root = (props) => {
  return (
    <>{props.children}</>
  )
}



render(() => (
  <Router root={Root}>
    <Route path="/" component={App} />
  </Router>
), document.getElementById("root"));
