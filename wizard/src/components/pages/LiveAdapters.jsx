import { faComputer, faServer, faTable } from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { useState, useEffect } from 'react';
import { Alert, Container, Spinner, Table } from 'react-bootstrap';

const LiveAdapters = () => {
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [routingTable, setRoutingTable] = useState({});

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await fetch("/api/peer-table");
        if (!response.ok) {
          throw new Error(`Error: ${response.status}`);
        }
        const result = await response.json();
        setRoutingTable(result);
      } catch (err) {
        console.error(err);
        setError(err.message);
      } finally {
        setLoading(false);
      }
    }

    fetchData();

    const intervalId = setInterval(fetchData, 5000);

    return () => clearInterval(intervalId);
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
                .filter(([_, item]) => item.Role !== "wizard")
                .map(([key, peer]) => (
                  <tr key={key}>
                    <td>{peer.IP}</td>
                    <td>{peer.Role === "adapter-client" ? (
                      <span>
                        <FontAwesomeIcon icon={faComputer} />{" "}
                        Client
                      </span>
                    ) : peer.Role === "adapter-server" ? (
                      <span>
                        <FontAwesomeIcon icon={faServer} />{" "}
                        Server
                      </span>
                    ) : peer.Role}</td>
                    <td>{peer.FromIP}</td>
                    <td>{peer.ToIP}</td>
                    <td>{new Date(peer.LastSeen).toLocaleString()}</td>
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
