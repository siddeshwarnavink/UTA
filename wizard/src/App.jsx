import { useState, useEffect } from 'react';
import { BrowserRouter, Outlet, Route, Routes } from 'react-router-dom';
import { Alert, Spinner } from 'react-bootstrap';
import 'bootstrap/dist/css/bootstrap.min.css';

import './App.css'
import LiveAdapters from './components/pages/LiveAdapters';
import TheNavbar from './components/layout/TheNavbar';
import NetworkGraph from './components/pages/NetworkGraph'

const App = () => {
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [routingTable, setRoutingTable] = useState({});
  const [transmission, setTransmission] = useState(null);

  useEffect(() => {
    let ws;
    try {
      ws = new WebSocket('ws://localhost:3300/ws');

      ws.onmessage = (event) => {
        try {
          const parsedData = JSON.parse(event.data);
          if ("peerTable" in parsedData && parsedData.peerTable !== null) {
            console.log("peerTable", parsedData.peerTable)
            setError(null);
            setLoading(false);
            setRoutingTable(parsedData.peerTable);
            setTransmission(null);
          } else if ("transmission" in parsedData && parsedData.transmission !== null) {
            console.log("transmission", parsedData.transmission)
            setTransmission(prevState => {
              if (prevState != null) {
                return { ...prevState, [parsedData.transmission.ip]: parsedData.transmission.sent }
              } else {
                return { [parsedData.transmission.ip]: parsedData.transmission.sent }
              }
            });
          }
        } catch (err) {
          console.error("Error parsing message:", err);
        }
      };

      ws.onerror = (event) => {
        setError("WebSocket encountered an error.");
        console.error(event);
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

  if (loading) {
    return (
      <Spinner animation="border" role="status">
        <span className="visually-hidden">Loading...</span>
      </Spinner>
    );
  } else if (error) {
    return <Alert variant="danger">{error}</Alert>;
  }


  return (
    <BrowserRouter>
      <Routes>
        <Route
          element={
            <>
              <TheNavbar />
              <Outlet />
            </>
          }
        >
          <Route path="/" element={<LiveAdapters routingTable={routingTable} transmission={transmission} />} />
          <Route path="/network" element={<NetworkGraph routingTable={routingTable} transmission={transmission} />} />
        </Route>
      </Routes>
    </BrowserRouter>
  );
};

export default App;
