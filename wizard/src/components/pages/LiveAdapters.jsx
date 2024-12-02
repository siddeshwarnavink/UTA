import { faComputer, faServer, faTable, faUpload, faDownload } from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { Container, Table } from 'react-bootstrap';

import { Link } from 'react-router-dom';

const LiveAdapters = ({ routingTable, transmission }) => {
  return (
    <Container>
      <h3 className="my-4">
        <FontAwesomeIcon icon={faTable} />{" "}Live Adapters
      </h3>
      <Table>
        <thead>
          <tr>
            <th>UDP IP</th>
            <th></th>
            <th>Role</th>
            <th>From IP</th>
            <th>To IP</th>
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
                      <FontAwesomeIcon icon={faComputer} />{" "}
                      Client
                    </span>
                  ) : peer.role === "adapter-server" ? (
                    <span>
                      <FontAwesomeIcon icon={faServer} />{" "}
                      Server
                    </span>
                  ) : peer.role}</td>
                  <td>{peer.from_ip}</td>
                  <td>{peer.to_ip}</td>
                  <td>{new Date(peer.last_seen).toLocaleString()}</td>
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
