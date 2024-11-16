import { Navbar, Container, Nav } from 'react-bootstrap';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faCircleNodes, faInfinity, faTable } from '@fortawesome/free-solid-svg-icons';
import { LinkContainer } from 'react-router-bootstrap';

const TheNavbar = () => {
  return (
    <Navbar expand="lg" className="bg-body-tertiary" data-bs-theme="dark">
      <Container>
        <LinkContainer to="/">
          <Navbar.Brand>
            <FontAwesomeIcon icon={faInfinity} />{"  "}
            UTA Wizard
          </Navbar.Brand>
        </LinkContainer>
        <Navbar.Toggle aria-controls="navbar-nav" />
        <Navbar.Collapse id="navbar-nav">
          <Nav className="me-auto">
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
        </Navbar.Collapse>
      </Container>
    </Navbar>
  );
}

export default TheNavbar;
