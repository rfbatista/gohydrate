import * as React from "react";

const NavigationBar = (props) => {
  return (
    <h1>
      Hello from React!<div>{props.title}</div>
    </h1>
  );
};

export default NavigationBar;
