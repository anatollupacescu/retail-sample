import * as React from 'react'
import 'bootstrap/dist/css/bootstrap.min.css';
import {MyTabs} from "~ui/component/Tabs";

export const App = () => (
    <div>
      <title>Dashboard</title>
      <main role="main" className="container">
        <MyTabs/>
      </main>
    </div>
);