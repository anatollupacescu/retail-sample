import * as React from 'react'
import Tabs from "react-bootstrap/Tabs";
import Tab from "react-bootstrap/Tab";
import {MySpinner} from "~ui/component/Spinner";
import {Wip} from "~ui/component/Wip";

export const MyTabs = () => {
  return (
      <Tabs defaultActiveKey="profile" id="uncontrolled-tab-example">
        <Tab eventKey="home" title="Home">
          <MySpinner/>
        </Tab>
        <Tab eventKey="profile" title="Profile">
          <Wip/>
        </Tab>
        <Tab eventKey="contact" title="Contact" disabled>
          Contact
        </Tab>
      </Tabs>
  )
};