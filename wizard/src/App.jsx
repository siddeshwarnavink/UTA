import React, { useState, useEffect } from 'react';
import { BrowserRouter, Outlet, Route, Routes } from 'react-router-dom';
import { Alert, Spinner } from 'react-bootstrap';
import { v4 as uuidv4 } from 'uuid';
import './bootstrap.min.css'; // Source: https://bootswatch.com/lumen/
import './App.css';
import LiveAdapters from './components/pages/LiveAdapters';
import TheNavbar from './components/layout/TheNavbar';
import NetworkGraph from './components/pages/NetworkGraph';
import ViewAdapter from './components/pages/ViewAdapter';
import Login from './components/pages/Login';

const App = () => {
  const [ws, setWs] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [routingTable, setRoutingTable] = useState({});
  const [transmission, setTransmission] = useState(null);
  const [isAuthenticated, setIsAuthenticated] = useState(false);

  useEffect(() => {
    try {
      const myws = new WebSocket('ws://localhost:3300/ws');
      setWs(myws);

      myws.onmessage = (event) => {
        try {
          const parsedData = JSON.parse(event.data);
          if ('peerTable' in parsedData && parsedData.peerTable !== null) {
            setError(null);
            setLoading(false);
            setRoutingTable(parsedData.peerTable);
            setTransmission(null);
          } else if (
            'transmission' in parsedData &&
            parsedData.transmission !== null
          ) {
            setTransmission((prevState) => {
              if (prevState != null) {
                return { ...prevState, [parsedData.transmission.ip]: parsedData.transmission.sent };
              } else {
                return { [parsedData.transmission.ip]: parsedData.transmission.sent };
              }
            });
          } else {
            console.log('unknown message', parsedData);
          }
        } catch (err) {
          console.error('Error parsing message:', err);
        }
      };

      myws.onerror = (event) => {
        setError('WebSocket encountered an error.');
        console.error(event);
      };
    } catch (err) {
      setError('WebSocket initialization failed.');
    }

    return () => {
      if (ws !== null) {
        ws.close();
      }
    };
  }, []);

  const wsRequest = (ip, reqType) => {
    return new Promise((res, rej) => {
      if (!ws || ws.readyState !== WebSocket.OPEN) {
        rej('Websocket is closed. Try refreshing the page');
        return;
      }

      const id = uuidv4();
      try {
        const messageToSend = JSON.stringify({ reqId: id, request: reqType, payload: ip });
        ws.send(messageToSend);
      } catch (err) {
        rej('Failed to send config request');
      }

      ws.onmessage = (event) => {
        try {
          const parsedData = JSON.parse(event.data);
          if ('response' in parsedData && parsedData.response !== null) {
            if (parsedData.response.reqId === id) {
              res(parsedData.response.data);
            }
          }
        } catch (err) {
          rej('Failed to listen for config response');
        }
      };
    });
  };

  const requestConfig = (ip) => wsRequest(ip, 0);
  const requestLogs = (ip, page) => wsRequest(`${ip},${page}`, 1);
  const requestSaveConfig = (ip, config) => wsRequest(`${ip},${config}`, 2);

  if (!isAuthenticated) {
    return <Login setIsAuthenticated={setIsAuthenticated} />;
  }

  if (loading) {
    return (
      <Spinner animation="border" role="status">
        <span className="visually-hidden">Loading...</span>
      </Spinner>
    );
  }

  if (error) {
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
          <Route
            exact
            path="/"
            element={
              <LiveAdapters routingTable={routingTable} transmission={transmission} />
            }
          />
          <Route
            exact
            path="/network"
            element={
              <NetworkGraph routingTable={routingTable} transmission={transmission} />
            }
          />
          <Route
            path="/c/:ip"
            element={
              <ViewAdapter
                routingTable={routingTable}
                requestConfig={requestConfig}
                requestLogs={requestLogs}
                requestSaveConfig={requestSaveConfig}
              />
            }
          />
        </Route>
      </Routes>
    </BrowserRouter>
  );
};

export default App;
