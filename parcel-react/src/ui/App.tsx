import * as React from 'react'
import 'bootstrap/dist/css/bootstrap.min.css';
import Button from 'react-bootstrap/Button';
import {Tabs} from "~ui/component/Tabs";
import { useState } from 'react';
import Spinner from "react-bootstrap/Spinner";

export const App = () => {

  let [enabled, setEnabled ] = useState(true);

  let updateButtonState = () => {
    if (enabled) {
      setEnabled(false)
    }
    setTimeout(() => {
      setEnabled(true)
    }, 3000);
  };

  return (
      <div>
        <title>Dashboard</title>
        <main role="main" className="container">
          <Tabs/>
          <Button size={"lg"} onClick={updateButtonState} disabled={!enabled}>
            {(!enabled && (
                <Spinner animation="border" role="status" />
            )) || (
                "jora"
            ) }
          </Button>
        </main>
      </div>
  )
};