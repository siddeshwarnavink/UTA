import { faSave } from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { Editor } from '@monaco-editor/react';
import { useEffect } from 'react';
import { Button, Container, Tab, Tabs } from 'react-bootstrap';
import { useParams } from 'react-router-dom';

const AdapterConfig = ({ requestConfig }) => {
  let { ip } = useParams();
  useEffect(() => {
    requestConfig(ip);
  }, []);

  return (
    <Container>
      <h3 className="my-4">
        {ip}
      </h3>
      <Tabs
        defaultActiveKey="config"
        className="mb-3"
        fill
      >
        <Tab eventKey="config" title="Config">
          <Button className="float-end m-3 mb-4">
            <FontAwesomeIcon icon={faSave} />
            {" "}Save
          </Button>
          <Editor height="90vh" defaultLanguage="lua" defaultValue="-- some comment" />;
        </Tab>
        <Tab eventKey="activity" title="Activity">
          Activity log...
        </Tab>
      </Tabs>
    </Container>
  );
};

export default AdapterConfig;
