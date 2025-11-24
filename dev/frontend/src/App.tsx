import { useState } from "react";
import { useStore } from "./stores";
import "./App.css";

function App() {
  const { user, isLoggedIn, login, logout, messages, addMessage } = useStore();
  const [username, setUsername] = useState("");
  const [messageText, setMessageText] = useState("");

  const handleLogin = () => {
    if (username.trim()) {
      login(username);
      setUsername("");
    }
  };

  const handleSendMessage = () => {
    if (messageText.trim() && user) {
      addMessage(messageText, user.name);
      setMessageText("");
    }
  };

  return (
    <div className="app-container">
      <h1>Zustand Chat App Demo</h1>

      <div className="card">
        <h2>User Status</h2>
        {isLoggedIn ? (
          <div>
            <p>
              Logged in as: <strong>{user?.name}</strong>
            </p>
            <button onClick={logout}>Logout</button>
          </div>
        ) : (
          <div className="login-form">
            <input
              type="text"
              placeholder="Enter username"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
            />
            <button onClick={handleLogin}>Login</button>
          </div>
        )}
      </div>

      <div className="card">
        <h2>Chat Messages</h2>
        <div className="messages-list">
          {messages.length === 0 ? (
            <p className="no-messages">No messages yet.</p>
          ) : (
            messages.map((msg) => (
              <div key={msg.id} className="message-item">
                <span className="sender">{msg.sender}:</span>
                <span className="text">{msg.text}</span>
                <span className="time">
                  {new Date(msg.timestamp).toLocaleTimeString()}
                </span>
              </div>
            ))
          )}
        </div>

        {isLoggedIn && (
          <div className="message-input">
            <input
              type="text"
              placeholder="Type a message..."
              value={messageText}
              onChange={(e) => setMessageText(e.target.value)}
              onKeyDown={(e) => e.key === "Enter" && handleSendMessage()}
            />
            <button onClick={handleSendMessage}>Send</button>
          </div>
        )}
      </div>
    </div>
  );
}

export default App;
