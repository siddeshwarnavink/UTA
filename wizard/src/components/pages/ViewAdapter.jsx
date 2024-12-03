import { Container, Tab, Tabs } from 'react-bootstrap';
import { useParams } from 'react-router-dom';

import AdapterConfig from '../adapter/AdapterConfig';
import AdapterLogs from '../adapter/AdapterLogs';

const ViewAdapter = ({ requestConfig, requestLogs, requestSaveConfig }) => {
  const { ip } = useParams();

  return (
    <Container>
      <h3 className="my-4">{ip}</h3>
      <Tabs
        defaultActiveKey="config"
        className="mb-3"
        fill
      >
        <Tab eventKey="config" title="Config" >
          <AdapterConfig ip={ip} requestConfig={requestConfig} requestSaveConfig={requestSaveConfig} />
        </Tab>
        <Tab eventKey="activity" title="Activity" mountOnEnter>
          <AdapterLogs ip={ip} requestLogs={requestLogs} />
        </Tab>
      </Tabs>
    </Container>
  );
}

export default ViewAdapter
