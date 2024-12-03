import { faSave } from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { Editor } from '@monaco-editor/react';
import { useEffect, useState } from 'react';
import { Alert, Button, Container, Spinner, Tab, Tabs } from 'react-bootstrap';
import { useParams } from 'react-router-dom';

const AdapterConfig = ({ requestConfig }) => {
  const { ip } = useParams();
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [config, setConfig] = useState("-- Loading ...")

  useEffect(() => {
    (async () => {
      try {
        const response = await requestConfig(ip);
        setLoading(false);
        setConfig(response);
      } catch (err) {
        setLoading(false);
        setError(JSON.stringify(err));
      }
    })();
  }, []);

  let content = <Spinner />;

  if (!loading && error) {
    content = (
      <Alert variant="danger">
        {error}
      </Alert>
    );
  } else if (!loading && !error) {
    content = (
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
          <Editor height="90vh" defaultLanguage="lua" defaultValue={config} />;
        </Tab>
        <Tab eventKey="activity" title="Activity">
          Activity log...
        </Tab>
      </Tabs>
    );
  }

  return (
    <Container>
      <h3 className="my-4">
        {ip}
      </h3>
      {content}
    </Container>
  );
};

export default AdapterConfig;
