import { useMemo } from 'react';
import { faComputer, faServer, faUpload, faDownload, faHatWizard } from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { Col, Container, Row, Table, Card } from 'react-bootstrap';
import { Link } from 'react-router-dom';
import moment from 'moment';

import classes from './LiveAdapters.module.css';
import liveIcon from '../../assets/live.gif';

const LiveAdapters = ({ routingTable, transmission }) => {

  const noClients = useMemo(() => {
    if (Object.keys(routingTable).length > 0) {
      return Object.entries(routingTable)
        .filter(([_, item]) => item.role === "adapter-client")
        .length;
    }
    return 0;
  }, [routingTable]);

  const noServer = useMemo(() => {
    if (Object.keys(routingTable).length > 0) {
      return Object.entries(routingTable)
        .filter(([_, item]) => item.role === "adapter-server")
        .length;
    }
    return 0;
  }, [routingTable]);

  const noWizard = useMemo(() => {
    if (Object.keys(routingTable).length > 0) {
      return Object.entries(routingTable)
        .filter(([_, item]) => item.role === "wizard")
        .length;
    }
    return 0;
  }, [routingTable]);

  return (
    <Container>
      <Row className={classes.useless_thing}>
        <Col>
          <Card className={[classes.card, classes.cyan].join(" ")}>
            <div>
              <div className={classes.count}>{noClients}</div>
              <div classname={classes.label}>
                Client
              </div>
            </div>
            <div className={classes.icon}>
              <FontAwesomeIcon icon={faComputer} />
            </div>
          </Card>
        </Col>
        <Col>
          <Card className={[classes.card, classes.yellow].join(" ")}>
            <div>
              <div className={classes.count}>{noServer}</div>
              <div classname={classes.label}>Server</div>
            </div>
            <div className={classes.icon}>
              <FontAwesomeIcon icon={faServer} />
            </div>
          </Card>
        </Col>
        <Col>
          <Card className={[classes.card, classes.red].join(" ")}>
            <div>
              <div className={classes.count}>{noWizard}</div>
              <div classname={classes.label}>
                Wizard
              </div>
            </div>
            <div className={classes.icon}>
              <FontAwesomeIcon icon={faHatWizard} />
            </div>
          </Card>
        </Col>
      </Row>
      <h3 className="my-4">
        <img src={liveIcon} alt="Live now..." className={classes.liveIcon} />
        Live Adapters
      </h3>
      <Table>
        <thead>
          <tr>
            <th>ID</th>
            <th></th>
            <th>Role</th>
            <th>TCP IP</th>
            <th>Last Seen</th>
          </tr>
        </thead>
        <tbody>
          {Object.keys(routingTable).length > 0 ? (
            Object.entries(routingTable)
              .filter(([_, item]) => item.role !== "wizard")
              .map(([key, peer]) => (
                <tr key={key}>
                  <td>
                    <Link to={'/c/' + peer.ip}>
                      {peer.ip}
                    </Link>
                  </td>
                  <td>
                    {(transmission !== null && peer.ip in transmission) ?
                      transmission[peer.ip] ? (
                        <FontAwesomeIcon icon={faUpload} />
                      ) : (
                        <FontAwesomeIcon icon={faDownload} />
                      )
                      : null}
                  </td>
                  <td>{peer.role === "adapter-client" ? (
                    <span>
                      <FontAwesomeIcon icon={faComputer} style={{ color: "#1c702f" }} />{" "}
                      Client
                    </span>
                  ) : peer.role === "adapter-server" ? (
                    <span>
                      <FontAwesomeIcon icon={faServer} style={{ color: "#cb9906" }} />{" "}
                      Server
                    </span>
                  ) : peer.role}</td>
                  <td>
                    {peer.role === "adapter-client" ? peer.from_ip : peer.to_ip}
                  </td>
                  <td>
                    {moment(peer.last_seen).fromNow()}
                  </td>
                </tr>
              ))
          ) : (
            <tr>
              <td colSpan="3">No peers available</td>
            </tr>
          )}
        </tbody>
      </Table>
    </Container>
  );
}

export default LiveAdapters;
