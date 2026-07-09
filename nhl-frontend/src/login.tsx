import "./App.css";

import { useState } from "react";
import { useNavigate } from "react-router-dom";

export function Login() {
  const navigate = useNavigate();
  const [error, setError] = useState("");
  const [status, setStatus] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [loginEmail, setLoginEmail] = useState("");
  const [loginPassword, setLoginPassword] = useState("");

  //SIGN UP FOR EXISTING USER

  const handleLogin = async (e: React.SyntheticEvent<HTMLFormElement>) => {
    e.preventDefault();
    setStatus("Login");

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
  if (status === "SignUp") {
    return (
      <div>
        <h1>NHL Tracker</h1>
        <h2>Sign Up</h2>
        <form onSubmit={handleCreateUser}>
          <div className="form-container">
            <label className="form-label">Email</label>
            <input type="email" value={email} onChange={(e) => setEmail(e.target.value)} placeholder="Email" className="form-input" />
          </div>
          <div className="form-container">
            <label className="form-label">Password</label>
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="Password"
              className="form-input"
            />
          </div>
          <button type="submit">Sign up</button>
        </form>
        <button onClick={() => setStatus("Login")}>Existing user?</button>
      </div>
    );
  }
  return (
    <div className="login-page">
      <h1>NHL Tracker</h1>
      <h2>Login</h2>
      <form onSubmit={handleLogin}>
        <div className="form-container">
          <label className="form-label">Email</label>
          <input
            type="email"
            value={loginEmail}
            onChange={(e) => setLoginEmail(e.target.value)}
            placeholder="Email"
            className="form-input"
          />
        </div>
        <div className="form-container">
          <label className="form-label">Password</label>
          <input
            type="password"
            value={loginPassword}
            onChange={(e) => setLoginPassword(e.target.value)}
            placeholder="Password"
            className="form-input"
          />
          {error && <p>{error}</p>}
        </div>
        <button type="submit">Login</button>
      </form>
      <button onClick={() => setStatus("SignUp")}>Create new account</button>
    </div>
  );
}
