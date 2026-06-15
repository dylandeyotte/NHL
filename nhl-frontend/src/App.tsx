import "./App.css";
import { useEffect, useState } from "react";
import { Routes, Route, useNavigate, useSearchParams } from "react-router-dom";

function Following() {
  const [following, setFollowing] = useState("");

  async function loadFollowing() {
    // HTTP request
    const data = await fetchHelper("http://localhost:8080/api/following");
    setFollowing(data);
  }

  // Run effect
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
  const navigate = useNavigate();

  try {
    // HTTP request
    const response = await fetch("http://localhost:8080/api/refresh", {
      method: "POST",
      headers: {
        Authorization: `Bearer ${localStorage.getItem("refreshToken")}`,
      },
    });
    // Check if token expired
    if (response.status === 401) {
      navigate("/");
    }
    // Store new token
    const data = await response.json();
    localStorage.setItem("token", data.token);
  } catch (err) {
    console.error(`Request failed: ${err}`);
    throw err;
  }
  console.log("new token issued");
}

async function fetchHelper(url: string, options?: string) {
  try {
    // HTTP request
    let response = await fetch(url, {
      method: options,
      headers: {
        Authorization: `Bearer ${localStorage.getItem("token")}`,
      },
    });
    // If failed, get new token and request again
    if (response.status === 401) {
      await newToken();
      response = await fetch(url, {
        method: options,
        headers: {
          Authorization: `Bearer ${localStorage.getItem("token")}`, // CHEKC THIS
        },
      });
    }
    return response.json();
  } catch (err) {
    console.error(`Request failed: ${err}`);
    throw err;
  }
}

function Home() {
  const navigate = useNavigate();
  const [homeData, setHomeData] = useState(null);
  const [searchTerm, setSearchTerm] = useState("");

  async function loadHome() {
    // HTTP request
    const data = await fetchHelper("http://localhost:8080/api/home");
    setHomeData(data);
  }

  // Run effect
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
  const [results, setResults] = useState<searchedPlayer[]>([]);
  const [searchTerm, setSearchTerm] = useState("");
  const [followed, setFollowed] = useState<string[]>([]);
  const navigate = useNavigate();

  type searchedPlayer = {
    playerId: string;
    name: string;
    positionCode: string;
    teamAbbrev: string;
    height: string;
    weightInPounds: number;
    birthCountry: string;
  };

  // Get search parameter that user typed
  const player = searchParams.get("player");

  async function loadSearch() {
    // HTTP request
    const data = await fetchHelper(`http://localhost:8080/api/players/search?player=${player}`);
    setResults(data);
  }

  // Run effect
  useEffect(() => {
    loadSearch();
  }, [player]);

  async function handleSearch(e: React.SyntheticEvent<HTMLFormElement>) {
    e.preventDefault();
    navigate(`/search?player=${encodeURIComponent(searchTerm)}`);
  }

  async function handleFollow(id: string) {
    // HTTP request
    const data = await fetchHelper(`http://localhost:8080/api/players/${id}/follow`, "POST");

    console.log(data);
    console.log("player followed");
    setFollowed([...followed, id]); // LOOK INTO THIS
  }

  return (
    <div>
      <h1>Search Results</h1>
      <button onClick={() => navigate("/home")}>Home</button>
      <form onSubmit={handleSearch}>
        <input type="search" value={searchTerm} onChange={(e) => setSearchTerm(e.target.value)} />
      </form>
      <pre>
        {results.map((player) => (
          <div key={player.playerId}>
            <div>{player.name}</div>
            <div>{player.teamAbbrev}</div>
            <div>{player.positionCode}</div>
            <div>{player.height}</div>
            <div>{player.weightInPounds}</div>
            <div>{player.birthCountry}</div>
            <button onClick={() => handleFollow(player.playerId)}>{followed.includes(player.playerId) ? "Following" : "Follow"}</button>
          </div>
        ))}
      </pre>
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
