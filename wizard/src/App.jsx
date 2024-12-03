import { useState, useEffect } from 'react';
import { BrowserRouter, Outlet, Route, Routes } from 'react-router-dom';
import { Alert, Spinner } from 'react-bootstrap';
import { v4 as uuidv4 } from 'uuid';
import 'bootstrap/dist/css/bootstrap.min.css';

import './App.css'
import LiveAdapters from './components/pages/LiveAdapters';
import TheNavbar from './components/layout/TheNavbar';
import NetworkGraph from './components/pages/NetworkGraph'
import AdapterConfig from './components/pages/AdapterConfig'

const App = () => {
  const [ws, setWs] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [routingTable, setRoutingTable] = useState({});
  const [transmission, setTransmission] = useState(null);

  useEffect(() => {
    try {
      const myws = new WebSocket('ws://localhost:3300/ws');
      setWs(myws);

      myws.onmessage = (event) => {
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
          } else {
            console.log("unknown message", parsedData)
          }
        } catch (err) {
          console.error("Error parsing message:", err);
        }
      };

      myws.onerror = (event) => {
        setError("WebSocket encountered an error.");
        console.error(event);
      };
    } catch (err) {
      setError("WebSocket initialization failed.");
    }

    return () => {
      if (ws !== null) {
        ws.close();
      }
    };
  }, []);

  const requestConfig = ip => {
    return new Promise((res, rej) => {
      if (!ws || ws.readyState !== WebSocket.OPEN) {
        console.error('socket not open or ready');
        rej(null);
        return;
      }

      // send request to ws
      const id = uuidv4();
      try {
        const messageToSend = JSON.stringify({ reqId: id, request: 0, payload: ip });
        ws.send(messageToSend);
      } catch (err) {
        console.error('Error sending message:', err);
        rej(err);
      }

      // listen for response in ws
      ws.onmessage = (event) => {
        try {
          const parsedData = JSON.parse(event.data);
          if ("response" in parsedData && parsedData.response !== null) {
            console.log("response", parsedData.response)
            if (parsedData.response.reqId === id) {
              res(parsedData.response.data);
            }
          }
        } catch (err) {
          console.error("Error parsing message:", err);
          rej(err);
        }
      };

    });
  }

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
          <Route exact path="/" element={
            <LiveAdapters
              routingTable={routingTable}
              transmission={transmission} />
          } />
          <Route exact path="/network" element={
            <NetworkGraph
              routingTable={routingTable}
              transmission={transmission} />} />
          <Route path="/c/:ip" element={
            <AdapterConfig
              routingTable={routingTable}
              requestConfig={requestConfig}
            />
          } />
        </Route>
      </Routes>
    </BrowserRouter>
  );
};

export default App;
