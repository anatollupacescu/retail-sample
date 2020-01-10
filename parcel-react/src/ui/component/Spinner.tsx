import Button from "react-bootstrap/Button";
import Spinner from "react-bootstrap/Spinner";
import * as React from "react";
import {useState} from "react";
import Col from "react-bootstrap/Col";
import Container from "react-bootstrap/Container";
import Row from "react-bootstrap/Row";

export const MySpinner = () => {
  let [enabled, setEnabled] = useState(true);

  let updateButtonState = () => {
    if (enabled) {
      setEnabled(false)
    }
    setTimeout(() => {
      setEnabled(true)
    }, 3000);
  };

  return (
  <Container>
    <Row className="py-3">
      <Col>
        <Button size={"lg"} onClick={updateButtonState} disabled={!enabled}>
          {(!enabled && (
              <Spinner animation="border" role="status"/>
          )) || (
              "jora"
          )}
        </Button>
      </Col>
    </Row>
  </Container>
  )
};