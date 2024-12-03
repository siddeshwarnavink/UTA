import { Navbar, Container, Nav, Form, Button } from 'react-bootstrap';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faBell, faCircleNodes, faInfinity, faQuestion, faTable } from '@fortawesome/free-solid-svg-icons';
import { LinkContainer } from 'react-router-bootstrap';

const TheNavbar = () => {
  return (
    <Navbar expand="lg" className="bg-primary" data-bs-theme="dark">
      <Container>
        <LinkContainer to="/">
          <Navbar.Brand>
            <FontAwesomeIcon icon={faInfinity} />{"   "}
            UTA Wizard
          </Navbar.Brand>
        </LinkContainer>
        <Navbar.Toggle aria-controls="navbar-nav" />
        <Navbar.Collapse id="navbar-nav">
          <>
            <Nav className="justify-content-end flex-grow-1 pe-3">
              <LinkContainer to="/">
                <Nav.Link>
                  <FontAwesomeIcon icon={faTable} />{" "}
                  Live Adapters
                </Nav.Link>
              </LinkContainer>
              <LinkContainer to="/network">
                <Nav.Link>
                  <FontAwesomeIcon icon={faCircleNodes} />{" "}
                  Network
                </Nav.Link>
              </LinkContainer>
            </Nav>
            <Form className="d-flex">
              <Form.Control
                type="search"
                placeholder="Search..."
                className="me-2"
                aria-label="Search"
              />
              <Button variant="secondary">Search</Button>
            </Form>
            <Nav className="mx-3">
              <Nav.Link>
                <FontAwesomeIcon icon={faBell} />
              </Nav.Link>
              <Nav.Link>
                <FontAwesomeIcon icon={faQuestion} />
              </Nav.Link>
            </Nav>
          </>
        </Navbar.Collapse>
      </Container>
    </Navbar>
  );
}

export default TheNavbar;
