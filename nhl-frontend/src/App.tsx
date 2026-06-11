import "./App.css";
import { useEffect, useState } from "react";
import { Routes, Route, useNavigate, useSearchParams } from "react-router-dom";

function Following() {
  const [following, setFollowing] = useState("");

  async function loadFollowing() {
    const data = await fetchHelper("http://localhost:8080/api/following");
    setFollowing(data);
  }

  useEffect(() => {
    loadFollowing();
  }, []);

  return (
    <div>
      <h1>Following</h1>
      <pre>{JSON.stringify(following, null, 2)}</pre>
    </div>
  );
}
function Login() {
  const navigate = useNavigate();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [loginEmail, setLoginEmail] = useState("");
  const [loginPassword, setLoginPassword] = useState("");

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
    localStorage.setItem("refreshToken", data.refresh_token);

    console.log("token stored");
    console.log(data);

    navigate("/home");
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

async function newToken() {
  // REFRESH EXPIRED CHECK
  const refreshToken = localStorage.getItem("refreshToken");
  const response = await fetch("http://localhost:8080/api/refresh", {
    method: "POST",
    headers: {
      Authorization: `Bearer ${refreshToken}`,
    },
  });
  const data = await response.json();
  localStorage.setItem("token", data.token);
  console.log("new token issued");
}

async function fetchHelper(url: string) {
  try {
    let response = await fetch(url, {
      headers: {
        Authorization: `Bearer ${localStorage.getItem("token")}`,
      },
    });
    if (response.status === 401) {
      await newToken();
      response = await fetch(url, {
        headers: {
          Authorization: `Bearer ${localStorage.getItem("token")}`, // CHEKC THIS
        },
      });
    }
    return response.json();
  } catch (err) {
    console.error(`Request failed: ${err}`);
  }
}

function Home() {
  const navigate = useNavigate();
  const [homeData, setHomeData] = useState(null);
  const [searchTerm, setSearchTerm] = useState("");

  async function loadHome() {
    const data = await fetchHelper("http://localhost:8080/api/home");
    setHomeData(data);
  }

  useEffect(() => {
    loadHome();
  }, []);

  async function handleSearch() {
    navigate(`/search?player=${encodeURIComponent(searchTerm)}`);
  }

  return (
    <div>
      <h1>Home</h1>
      <button onClick={() => navigate("/following")}>Following</button>
      <form onSubmit={handleSearch}>
        <input type="search" value={searchTerm} onChange={(e) => setSearchTerm(e.target.value)} />
      </form>
      <pre>{JSON.stringify(homeData, null, 2)}</pre>
    </div>
  );
}

function Search() {
  const [searchParams] = useSearchParams();
  const [results, setResults] = useState([]);

  const player = searchParams.get("player");

  async function loadSearch() {
    const data = await fetchHelper(`http://localhost:8080/api/players/search?player=${player}`);
    setResults(data);
  }

  useEffect(() => {
    loadSearch();
  }, []);

  return (
    <div>
      <h1>Search Results</h1>
      <pre>{JSON.stringify(results, null, 2)}</pre>
    </div>
  );
}

function App() {
  return (
    <Routes>
      <Route path="/" element={<Login />} />
      <Route path="/following" element={<Following />} />
      <Route path="/home" element={<Home />} />
      <Route path="/search" element={<Search />} />
    </Routes>
  );
}

export default App;
