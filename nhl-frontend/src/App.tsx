import "./App.css";
import { useEffect, useState } from "react";

function App() {
  const [loggedIn, setLoggedIn] = useState(false);

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [loginEmail, setLoginEmail] = useState("");
  const [loginPassword, setLoginPassword] = useState("");

  const [homeData, setHomeData] = useState(null);

  const handleLogin = async (e: React.SyntheticEvent<HTMLFormElement>) => {
    e.preventDefault();

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

    const data = await response.json();

    localStorage.setItem("token", data.token);

    console.log("token stored");
    console.log(data);

    setLoggedIn(true);
  };

  const handleCreateUser = async (e: React.SyntheticEvent<HTMLFormElement>) => {
    e.preventDefault();

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
  };

  function Home() {
    async function fetchPlayers() {
      const token = localStorage.getItem("token");

      const response = await fetch("http://localhost:8080/api/home", {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });
      const data = await response.json();

      setHomeData(data);
    }

    useEffect(() => {
      fetchPlayers();
    }, []);

    return (
      <div>
        <h1>Home</h1>

        <pre>{JSON.stringify(homeData, null, 2)}</pre>
      </div>
    );
  }

  if (loggedIn) {
    return <Home />;
  }

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

export default App;
