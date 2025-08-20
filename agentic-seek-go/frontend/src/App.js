import React, { useState } from 'react';
import './App.css';

function App() {
  const [messages, setMessages] = useState([]);
  const [input, setInput] = useState('');
  const [file, setFile] = useState(null);
  const [loading, setLoading] = useState(false);

  const handleSend = async () => {
    if (input.trim() && !loading) {
      const newMessages = [...messages, { text: input, sender: 'user' }];
      setMessages(newMessages);
      setInput('');
      setLoading(true);

      try {
        const response = await fetch('/api/chat', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({ message: input }),
        });
        const data = await response.json();
        setMessages([...newMessages, { text: data.message, sender: 'bot' }]);
      } catch (error) {
        console.error('Error sending message:', error);
        setMessages([...newMessages, { text: 'Error sending message', sender: 'bot' }]);
      } finally {
        setLoading(false);
      }
    }
  };

  const handleFileChange = (e) => {
    setFile(e.target.files[0]);
  };

  const handleUpload = async () => {
    if (file && !loading) {
      setLoading(true);
      const formData = new FormData();
      formData.append('file', file);

      try {
        const response = await fetch('/api/upload', {
          method: 'POST',
          body: formData,
        });
        const data = await response.json();
        // For now, just log the extracted text
        console.log('Extracted text:', data.text);
        alert('File uploaded and text extracted. Check the console.');
      } catch (error) {
        console.error('Error uploading file:', error);
        alert('Error uploading file.');
      } finally {
        setLoading(false);
      }
    }
  };

  return (
    <div className="App">
      <div className="chat-window">
        {messages.map((msg, index) => (
          <div key={index} className={`message ${msg.sender}`}>
            {msg.text}
          </div>
        ))}
        {loading && <div className="message bot">Loading...</div>}
      </div>
      <div className="input-area">
        <input
          type="text"
          value={input}
          onChange={(e) => setInput(e.target.value)}
          onKeyPress={(e) => e.key === 'Enter' && handleSend()}
          disabled={loading}
        />
        <button onClick={handleSend} disabled={loading}>Send</button>
      </div>
      <div className="upload-area">
        <input type="file" onChange={handleFileChange} disabled={loading} />
        <button onClick={handleUpload} disabled={loading}>Upload</button>
      </div>
    </div>
  );
}

export default App;
