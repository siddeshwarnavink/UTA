import { BrowserRouter, Outlet, Route, Routes } from 'react-router-dom';
import 'bootstrap/dist/css/bootstrap.min.css';

import LiveAdapters from './components/pages/LiveAdapters';
import TheNavbar from './components/layout/TheNavbar';
import NetworkGraph from './components/pages/NetworkGraph'

const App = () => {
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
          <Route path="/" element={<LiveAdapters />} />
          <Route path="/network" element={<NetworkGraph />} />
        </Route>
      </Routes>
    </BrowserRouter>
  );
};

export default App;
