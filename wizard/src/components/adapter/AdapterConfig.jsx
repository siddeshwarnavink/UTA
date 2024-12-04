import { useEffect, useState } from 'react';
import { Alert, Button, Spinner } from 'react-bootstrap';
import { Editor } from '@monaco-editor/react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faRefresh, faSave } from '@fortawesome/free-solid-svg-icons';
import Swal from 'sweetalert2'

const AdapterConfig = ({ ip, requestConfig, requestSaveConfig }) => {
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

  const saveConfig = async () => {
    try {
      await requestSaveConfig(ip, config);
      Swal.fire({
        title: "Adapter configuration saved",
        icon: "success"
      });
    } catch (err) {
      Swal.fire({
        title: "Failed to save configuration",
        text: "Try refreshing the page and try again",
        icon: "error"
      });
    }
  }

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
        <div className="float-end">
          <Button className="m-3" variant="light" onClick={saveConfig}>
            <FontAwesomeIcon icon={faRefresh} />
          </Button>
          <Button onClick={saveConfig}>
            <FontAwesomeIcon icon={faSave} />
            {" "}Save
          </Button>
        </div>
        <Editor
          height="40vh"
          defaultLanguage="lua"
          value={config}
          onChange={val => setConfig(val)}
        />
      </>
    );
  }

  return content;
};


export default AdapterConfig;
