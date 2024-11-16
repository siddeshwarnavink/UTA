import { useState, useEffect } from 'react';

const PeerTable = () => {
  const [routingTable, setRoutingTable] = useState({});

  useEffect(() => {
    const socket = new WebSocket("ws://localhost:3300/ws/peer-table");

    socket.onmessage = (event) => {
      const data = JSON.parse(event.data);
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
    <div>
      <h1>Peer Table</h1>
      <table>
        <thead>
          <tr>
            <th>IP</th>
            <th>Address</th>
            <th>Last Seen</th>
          </tr>
        </thead>
        <tbody>
          {Object.keys(routingTable).length > 0 ? (
            Object.entries(routingTable).map(([key, peer]) => (
              <tr key={key}>
                <td>{peer.IP}</td>
                <td>{peer.Address}</td>
                <td>{new Date(peer.LastSeen).toLocaleString()}</td>
              </tr>
            ))
          ) : (
            <tr>
              <td colSpan="3">No peers available</td>
            </tr>
          )}
        </tbody>
      </table>
    </div>
  );
}

export default PeerTable;
