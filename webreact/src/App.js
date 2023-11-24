// src/App.js
import React, { useState, useEffect } from 'react';
import './App.css';
import { w3cwebsocket as W3CWebSocket } from 'websocket';

const client = new W3CWebSocket('ws://localhost:8080/ws');

function App() {
  const [messages, setMessages] = useState([]);
  const [input, setInput] = useState('');
  const [username, setUsername] = useState('');

  useEffect(() => {
    client.onopen = () => {
      console.log('WebSocket Client Connected');
    };
    client.onmessage = (message) => {
      const msg = JSON.parse(message.data);
      setMessages((prevMessages) => [...prevMessages, msg]);
    };
  }, []);

  const sendMessage = () => {
    if (input && username) {
      const message = { username, content: input };
      client.send(JSON.stringify(message));
      setInput('');
    }
  };

  return (
    <div className="App">
      <div>
        <h2>WebSocket Chat</h2>
        <div className="chat-container">
          {messages.map((msg, index) => (
            <div key={index} className="message">
              <strong>{msg.username}:</strong> {msg.content}
            </div>
          ))}
        </div>
        <div className="input-container">
          <input
            type="text"
            placeholder="Your username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
          />
          <input
            type="text"
            placeholder="Type a message..."
            value={input}
            onChange={(e) => setInput(e.target.value)}
          />
          <button onClick={sendMessage}>Send</button>
        </div>
      </div>
    </div>
  );
}

export default App;
