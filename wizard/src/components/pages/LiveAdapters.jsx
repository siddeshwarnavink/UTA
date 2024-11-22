import { faComputer, faServer, faTable } from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { useState, useEffect } from 'react';
import { Alert, Container, Spinner, Table } from 'react-bootstrap';

const LiveAdapters = () => {
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [routingTable, setRoutingTable] = useState({});

  useEffect(() => {
    let ws;
    try {
      ws = new WebSocket('ws://localhost:3300/ws');

      ws.onmessage = (event) => {
        try {
          const parsedData = JSON.parse(event.data);
          setError(null);
          setRoutingTable(parsedData);
          setLoading(false);
        } catch (err) {
          console.error("Error parsing message:", err);
        }
      };

      ws.onerror = (event) => {
        setError("WebSocket encountered an error.");
      };
    } catch (err) {
      setError("WebSocket initialization failed.");
    }

    return () => {
      if (ws) {
        ws.close();
      }
    };
  }, []);

  return (
    <Container>
      <h3 className="my-4">
        <FontAwesomeIcon icon={faTable} />{" "}Live Adapters
      </h3>
      {loading ? (
        <Spinner animation="border" role="status">
          <span className="visually-hidden">Loading...</span>
        </Spinner>
      ) : error ? (
        <Alert variant="danger">{error}</Alert>
      ) : (
        <Table>
          <thead>
            <tr>
              <th>UDP IP</th>
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
                    <td>{peer.ip}</td>
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

      )}
    </Container>
  );
}

export default LiveAdapters;
