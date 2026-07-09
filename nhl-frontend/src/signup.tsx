import "./App.css";

import { useState } from "react";

export function SignUp() {
  const [error, setError] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  //SIGN UP FOR EXISTING USER

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
      <h2>Sign Up</h2>
      <form onSubmit={handleCreateUser}>
        <div>
          <label>Email</label>
          <input type="email" value={email} onChange={(e) => setEmail(e.target.value)} placeholder="Email" />
        </div>
        <div>
          <label>Password</label>
          <input type="password" value={password} onChange={(e) => setPassword(e.target.value)} placeholder="Password" />
        </div>
        <button type="submit">Sign up</button>
      </form>
    </div>
  );
}
