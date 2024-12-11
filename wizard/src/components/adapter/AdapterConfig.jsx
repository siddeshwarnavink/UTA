import { useEffect, useState } from "react";
import { Alert, Button, Spinner, Form } from "react-bootstrap";
import { Editor } from "@monaco-editor/react";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faRefresh, faSave } from "@fortawesome/free-solid-svg-icons";
import Swal from "sweetalert2";

const AdapterConfig = ({ ip, requestConfig, requestSaveConfig }) => {
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [config, setConfig] = useState("-- Loading ...");
  const [showEditor, setShowEditor] = useState(false);
  const [formData, setFormData] = useState({
    mode: "serverMode",
    encryptPort: "127.0.0.1:9999",
    decryptPort: "127.0.0.1:8888",
    crypto: "AES",
  });

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
  }, [ip, requestConfig]);

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
  };

  const saveConfig = async () => {
    try {
      const generatedConfig = `
        conf.serverMode(${formData.mode === "serverMode"}) 
        conf.encryptPort("${formData.encryptPort}") 
        conf.decryptPort("${formData.decryptPort}") 
        conf.crypto("${formData.crypto}") 
      `;
      await requestSaveConfig(ip, generatedConfig);
      Swal.fire({
        title: "Adapter configuration saved",
        icon: "success",
      });
    } catch (err) {
      Swal.fire({
        title: "Failed to save configuration",
        text: "Try refreshing the page and try again",
        icon: "error",
      });
    }
  };

  let content = <Spinner />;

  if (!loading && error) {
    content = <Alert variant="danger">{error}</Alert>;
  } else if (!loading && !error) {
    content = (
      <>
        <h2>Configuration Form</h2>
        <Form>
          {/* 1. Mode */}
          <Form.Group className="mb-3">
            <Form.Label>Mode</Form.Label>
            <Form.Check
              type="radio"
              label="Server Mode"
              name="mode"
              value="serverMode"
              checked={formData.mode === "serverMode"}
              onChange={handleInputChange}
            />
            <Form.Check
              type="radio"
              label="Radio Mode"
              name="mode"
              value="radioMode"
              checked={formData.mode === "radioMode"}
              onChange={handleInputChange}
            />
          </Form.Group>

          {/* 2. Encrypt Port */}
          <Form.Group className="mb-3">
            <Form.Label>Encrypt Port</Form.Label>
            <Form.Control
              type="text"
              placeholder="Enter encrypt port"
              name="encryptPort"
              value={formData.encryptPort}
              onChange={handleInputChange}
            />
          </Form.Group>

          {/* 3. Decrypt Port */}
          <Form.Group className="mb-3">
            <Form.Label>Decrypt Port</Form.Label>
            <Form.Control
              type="text"
              placeholder="Enter decrypt port"
              name="decryptPort"
              value={formData.decryptPort}
              onChange={handleInputChange}
            />
          </Form.Group>

          {/* 4. Crypto */}
          <Form.Group className="mb-3">
            <Form.Label>Crypto</Form.Label>
            <Form.Control
              type="text"
              placeholder="Enter crypto type"
              name="crypto"
              value={formData.crypto}
              onChange={handleInputChange}
            />
          </Form.Group>
        </Form>

        {/* Advanced Button */}
        <Button
          variant="secondary"
          className="my-3"
          onClick={() => setShowEditor(!showEditor)}
        >
          Advanced
        </Button>

        {/* Editor Section */}
        {showEditor && (
          <Editor
            height="40vh"
            defaultLanguage="lua"
            value={config}
            onChange={(val) => setConfig(val)}
          />
        )}
      </>
    );
  }

  return content;
};

export default AdapterConfig;
