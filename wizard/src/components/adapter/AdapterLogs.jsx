import { useEffect, useState } from 'react';
import { Alert, Button, Spinner } from 'react-bootstrap';
import { Editor } from '@monaco-editor/react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faDownload, faRefresh } from '@fortawesome/free-solid-svg-icons';

const AdapterLogs = ({ ip, requestLogs }) => {
  const [page, setPage] = useState(1);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [logs, setLogs] = useState("-- Loading ...")

  const fetchLogs = async (p = undefined, append = false) => {
    try {
      const response = await requestLogs(ip, p || page);
      setLoading(false);
      if (append) {
        setLogs(curr => `${curr}${response}`);
      } else {
        setLogs(response);
      }
    } catch (err) {
      setLoading(false);
      setError(JSON.stringify(err));
    }
  };

  useEffect(() => {
    fetchLogs();
  }, []);

  const fetchNextPage = () => {
    fetchLogs(page + 1, true);
    setPage(curr => curr + 1);
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
        <Button className="float-end m-3 mb-4" onClick={() => fetchLogs()} variant="light">
          <FontAwesomeIcon icon={faRefresh} />
        </Button>
        <Editor height="30vh" defaultLanguage="text" value={logs} />
        <div className="d-flex justify-content-center mt-4">
          <Button variant="light" onClick={fetchNextPage}>
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
