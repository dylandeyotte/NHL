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

      if (response.status === 400) {
        setError("Email already in use");
        return;
      }

      const data = await response.json();
      console.log(data);
    } catch (err) {
      console.error(`Request failed" ${err}`);
      throw err;
    }
  };

  return (
    <div className="login-page">
      <div className="title">
        <h1>{status === "SignUp" ? "Create account" : "Welcome"}</h1>
      </div>
      <form onSubmit={status === "SignUp" ? handleCreateUser : handleLogin}>
        <div className="form-container">
          <label className="form-label">Email</label>
          <input
            type="email"
            value={status === "SignUp" ? email : loginEmail}
            onChange={status === "SignUp" ? (e) => setEmail(e.target.value) : (e) => setLoginEmail(e.target.value)}
            placeholder="Email"
            className="form-input"
          />
        </div>
        <div className="form-container">
          <label className="form-label">Password</label>
          <input
            type="password"
            value={status === "SignUp" ? password : loginPassword}
            onChange={status === "SignUp" ? (e) => setPassword(e.target.value) : (e) => setLoginPassword(e.target.value)}
            placeholder={status === "SignUp" ? "Create a Password" : "Password"}
            className="form-input"
          />
          <div className="error">{error && <p>{error}</p>}</div>
        </div>
        <div className="login-button">
          <button type="submit">{status === "SignUp" ? "Sign up" : "Login"}</button>
        </div>
      </form>
      <div className="sign-up-button">
        <button onClick={status === "SignUp" ? () => setStatus("Login") : () => setStatus("SignUp")}>
          {status === "SignUp" ? "Existing user?" : "Create new account"}
        </button>
      </div>
    </div>
  );
}
