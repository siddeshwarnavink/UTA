import { useEffect, useState } from 'react';
import { Alert, Button, Spinner } from 'react-bootstrap';
import { Editor } from '@monaco-editor/react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faDownload, faRefresh } from '@fortawesome/free-solid-svg-icons';

const AdapterLogs = ({ ip, requestLogs }) => {
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [logs, setLogs] = useState("-- Loading ...")

  const fetchLogs = async () => {
    try {
      const response = await requestLogs(ip);
      setLoading(false);
      setLogs(response);
    } catch (err) {
      setLoading(false);
      setError(JSON.stringify(err));
    }
  };

  useEffect(() => {
    fetchLogs();
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
        <Button className="float-end m-3 mb-4" onClick={() => fetchLogs()} variant="light">
          <FontAwesomeIcon icon={faRefresh} />
          {" "}Refresh
        </Button>
        <Editor height="30vh" defaultLanguage="text" defaultValue={logs} />
        <div class="d-flex justify-content-center mt-4">
          <Button variant="light">
            <FontAwesomeIcon icon={faDownload} />
            {" "}Load more
          </Button>
        </div>
      </>
    );
  }

  return content;
};


export default AdapterLogs;
