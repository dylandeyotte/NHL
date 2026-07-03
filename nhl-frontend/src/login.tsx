import "./App.css";

import { useState } from "react";
import { useNavigate } from "react-router-dom";

export function Login() {
  const navigate = useNavigate();
  const [error, setError] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [loginEmail, setLoginEmail] = useState("");
  const [loginPassword, setLoginPassword] = useState("");

  const handleLogin = async (e: React.SyntheticEvent<HTMLFormElement>) => {
    e.preventDefault();

    try {
      // HTTP request
      const response = await fetch("http://localhost:8080/api/login", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          email: loginEmail,
          password: loginPassword,
        }),
      });

      if (response.status === 400 || response.status === 401) {
        setError("Incorrect email or password");
        return;
      }

      // Store tokens
      const data = await response.json();
      localStorage.setItem("token", data.token);
      localStorage.setItem("refreshToken", data.refresh_token);

      console.log("token stored");
      navigate("/home");
    } catch (err) {
      console.error(`Request failed: ${err}`);
      throw err;
    }
  };

  const handleCreateUser = async (e: React.SyntheticEvent<HTMLFormElement>) => {
    e.preventDefault();

    try {
      // HTTP request
      const response = await fetch("http://localhost:8080/api/users", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          email: email,
          password: password,
        }),
      });

      const data = await response.json();
      console.log(data);
    } catch (err) {
      console.error(`Request failed" ${err}`);
      throw err;
    }
  };

  return (
    <div>
      <h1>NHL Tracker</h1>
      <h2>Existing user?</h2>
      <form onSubmit={handleLogin}>
        <div>
          <label>Email</label>
          <input type="email" value={loginEmail} onChange={(e) => setLoginEmail(e.target.value)} />
        </div>
        <div>
          <label>Password</label>
          <input type="password" value={loginPassword} onChange={(e) => setLoginPassword(e.target.value)} />
          {error && <p>{error}</p>}
        </div>
        <button type="submit">Login</button>
      </form>
      <h2>New user?</h2>
      <form onSubmit={handleCreateUser}>
        <div>
          <label>Email</label>
          <input type="email" value={email} onChange={(e) => setEmail(e.target.value)} />
        </div>
        <div>
          <label>Password</label>
          <input type="password" value={password} onChange={(e) => setPassword(e.target.value)} />
        </div>
        <button type="submit">Sign up</button>
      </form>
    </div>
  );
}
