import { useEffect, useState } from 'react';
import { Alert, Button, Spinner } from 'react-bootstrap';
import { Editor } from '@monaco-editor/react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faSave } from '@fortawesome/free-solid-svg-icons';

const AdapterConfig = ({ ip, requestConfig }) => {
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
      <>
        <Button className="float-end m-3 mb-4">
          <FontAwesomeIcon icon={faSave} />
          {" "}Save
        </Button>
        <Editor height="40vh" defaultLanguage="lua" defaultValue={config} />
      </>
    );
  }

  return content;
};


export default AdapterConfig;
