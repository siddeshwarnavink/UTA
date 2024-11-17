import { faComputer, faServer, faTable } from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { useState, useEffect } from 'react';
import { Container, Table } from 'react-bootstrap';

const LiveAdapters = () => {
  const [routingTable, setRoutingTable] = useState({});

  useEffect(() => {
    const socket = new WebSocket("ws://localhost:3300/ws/peer-table");

    socket.onmessage = (event) => {
      const data = JSON.parse(event.data);
      console.log("peer table", data);
      setRoutingTable(data);
    };

    socket.onerror = (error) => {
      console.error("WebSocket error:", error);
    };

    return () => {
      socket.close();
    };
  }, []);

  return (
    <Container>
      <h3 className="my-4">
        <FontAwesomeIcon icon={faTable} />{" "}Live Adapters
      </h3>
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
    </Container>
  );
}

export default LiveAdapters;
