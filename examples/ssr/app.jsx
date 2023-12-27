import * as React from "react";
import { TestComponent } from "./component";

const NavigationBar = (props) => {
  return (
    <h1>
      Hello from React!
      <TestComponent title={props.title} />
    </h1>
  );
};

export default NavigationBar;
