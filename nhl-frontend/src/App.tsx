import "./App.css";
import { teamColours } from "./data/team_colours";
import { useEffect, useState } from "react";
import { Routes, Route, useNavigate, useSearchParams } from "react-router-dom";

const teamList = [
  { name: "Anaheim Ducks", tricode: "ANA" },
  { name: "Boston Bruins", tricode: "BOS" },
  { name: "Buffalo Sabres", tricode: "BUF" },
  { name: "Calgary Flames", tricode: "CGY" },
  { name: "Carolina Hurricanes", tricode: "CAR" },
  { name: "Chicago Blackhawks", tricode: "CHI" },
  { name: "Colorado Avalanche", tricode: "COL" },
  { name: "Columbus Blue Jackets", tricode: "CBJ" },
  { name: "Dallas Stars", tricode: "DAL" },
  { name: "Detroit Red Wings", tricode: "DET" },
  { name: "Edmonton Oilers", tricode: "EDM" },
  { name: "Florida Panthers", tricode: "FLA" },
  { name: "Los Angeles Kings", tricode: "LAK" },
  { name: "Minnesota Wild", tricode: "MIN" },
  { name: "Montréal Canadiens", tricode: "MTL" },
  { name: "Nashville Predators", tricode: "NSH" },
  { name: "New Jersey Devils", tricode: "NJD" },
  { name: "New York Islanders", tricode: "NYI" },
  { name: "New York Rangers", tricode: "NYR" },
  { name: "Ottawa Senators", tricode: "OTT" },
  { name: "Philadelphia Flyers", tricode: "PHI" },
  { name: "Pittsburgh Penguins", tricode: "PIT" },
  { name: "San Jose Sharks", tricode: "SJS" },
  { name: "Seattle Kraken", tricode: "SEA" },
  { name: "St. Louis Blues", tricode: "STL" },
  { name: "Tampa Bay Lightning", tricode: "TBL" },
  { name: "Toronto Maple Leafs", tricode: "TOR" },
  { name: "Utah Mammoth", tricode: "UTA" },
  { name: "Vancouver Canucks", tricode: "VAN" },
  { name: "Vegas Golden Knights", tricode: "VGK" },
  { name: "Washington Capitals", tricode: "WSH" },
  { name: "Winnipeg Jets", tricode: "WPG" },
];

type playerStats = {
  name: string;
  games_played: number;
  goals: number;
  assists: number;
  points: number;
  "p/gp": string;
  last_5_games_totals: string;
  playing_today: boolean;
  sweater_number: number;
  position: string;
  team_abbrev: string;
};

type teamStandings = {
  name: string;
  games_played: number;
  wins: number;
  losses: number;
  otl: number;
  points: number;
  rw: number;
  row: number;
  goal_differential: number;
  last_10: string;
};

async function handleFollowPlayer(id: string) {
  const data = await fetchHelper(`http://localhost:8080/api/players/${id}/follow`, "POST");

  console.log(data);
  console.log(`Player followed: ${data.PlayerID}`);
}

async function handleUnfollowPlayer(id: string) {
  const data = await fetchHelper(`http://localhost:8080/api/players/${id}/follow`, "DELETE");

  console.log(data);
  console.log(`Player unfollowed: ${id}`);
}

async function handleFollowTeam(tricode: string) {
  const data = await fetchHelper(`http://localhost:8080/api/teams/${tricode}/follow`, "POST");

  console.log(data);
  console.log(`Team followd: ${data.TriCode}`);
}

async function handleUnfollowTeam(tricode: string) {
  const data = await fetchHelper(`http://localhost:8080/api/teams/${tricode}/follow`, "DELETE");

  console.log(data);
  console.log(`Team unfollowd: ${tricode}`);
}

function Teams() {
  const navigate = useNavigate();
  const [followedTeam, setFollowedTeam] = useState();
  const [error, setError] = useState("");

  async function fetchTeam() {
    const data = await fetchHelper("http://localhost:8080/api/following");
    setFollowedTeam(data.team.TriCode);
  }

  async function handleFollowClick(tricode: string) {
    if (followedTeam !== tricode && followedTeam !== "") {
      setError("Only one team may be followed");
      return;
    }
    await handleFollowTeam(tricode);
    fetchTeam();
  }

  async function handleUnfollowClick(tricode: string) {
    await handleUnfollowTeam(tricode);
    fetchTeam();
    setError("");
  }

  useEffect(() => {
    fetchTeam();
  }, []);

  return (
    <div>
      <button onClick={() => navigate("/home")}>Home</button>
      {error && <p>{error}</p>}
      <pre>
        <div className="team-grid">
          {teamList.map((team) => (
            <div key={team.name} className="teams">
              <h3>{team.name}</h3>
              <img src={`https://assets.nhle.com/logos/nhl/svg/${team.tricode}_light.svg`} alt={team.name} className="team-logo" />
              <button onClick={() => (followedTeam === team.tricode ? handleUnfollowClick(team.tricode) : handleFollowClick(team.tricode))}>
                {followedTeam === team.tricode ? "Unfollow" : "Follow"}
              </button>
            </div>
          ))}
        </div>
      </pre>
    </div>
  );
}

function Following() {
  const [following, setFollowing] = useState<followee[]>([]);
  const [team, setTeam] = useState<followedTeam>();
  const navigate = useNavigate();

  type followee = {
    PlayerName: string;
    PlayerID: number;
  };

  type followedTeam = {
    TeamName: string;
    TriCode: string;
  };

  async function loadFollowing() {
    // HTTP request
    try {
      const data = await fetchHelper("http://localhost:8080/api/following");

      setFollowing(data.players);
      setTeam(data.team);
    } catch (err) {
      if (err instanceof Error && err.message === "REFRESH_EXPIRED") {
        navigate("/");
      }
    }
  }

  // Run effect
  useEffect(() => {
    loadFollowing();
  }, []);

  async function handleUnfollowPlayerClick(id: number) {
    await handleUnfollowPlayer(id.toString());
    loadFollowing();
  }

  async function handleUnfollowTeamClick(tricode: string | undefined) {
    if (tricode === undefined) {
      throw new Error("Undefined tricode");
    }
    await handleUnfollowTeam(tricode);
    loadFollowing();
  }

  if (team?.TeamName === "" && !following) {
    return (
      <div>
        <h1>Following</h1>
        <button onClick={() => navigate("/home")}>Home</button>
        <h3>None</h3>
      </div>
    );
  }

  return (
    <div>
      <h1>Following</h1>
      <button onClick={() => navigate("/home")}>Home</button>
      <pre>
        <h3>{team?.TeamName}</h3>
        {team?.TeamName && <button onClick={() => handleUnfollowTeamClick(team?.TriCode)}>Unfollow</button>}
      </pre>
      <pre>
        {following &&
          following.map((player) => (
            <div key={player.PlayerID}>
              <h3>{player.PlayerName}</h3>
              <button onClick={() => handleUnfollowPlayerClick(player.PlayerID)}>Unfollow</button>
            </div>
          ))}
      </pre>
    </div>
  );
}
function Login() {
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

async function newToken() {
  try {
    // HTTP request
    const response = await fetch("http://localhost:8080/api/refresh", {
      method: "POST",
      headers: {
        Authorization: `Bearer ${localStorage.getItem("refreshToken")}`,
      },
    });
    // Check if refresh token expired
    if (response.status === 401) {
      return false;
    }
    // Store new token
    const data = await response.json();
    localStorage.setItem("token", data.token);
  } catch (err) {
    console.error(`Request failed: ${err}`);
    throw err;
  }
  console.log("new token issued");
  return true;
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
      // Check if refresh token is expired
      const success = await newToken();
      if (!success) {
        throw new Error("REFRESH_EXPIRED");
      }
      // HTTP request again
      response = await fetch(url, {
        method: options,
        headers: {
          Authorization: `Bearer ${localStorage.getItem("token")}`,
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
  const [searchTerm, setSearchTerm] = useState("");
  const [followedTeam, setFollowedTeam] = useState("");
  const [homeData, setHomeData] = useState<playerStats[]>([]);
  const [standings, setStandings] = useState<teamStandings[]>([]);

  async function loadHome() {
    // HTTP request
    const data = await fetchHelper("http://localhost:8080/api/home");
    setHomeData(data.players);
    setStandings(data.standings);
  }

  async function fetchTeam() {
    const data = await fetchHelper("http://localhost:8080/api/following");
    setFollowedTeam(data.team.TeamName);
  }

  // Run effect
  useEffect(() => {
    (loadHome(), fetchTeam());
  }, []);

  async function handleSearch() {
    navigate(`/search?player=${encodeURIComponent(searchTerm)}`);
  }

  if (standings.length === 0 && homeData.length === 0) {
    return (
      <div>
        <h1>Home</h1>
        <button onClick={() => navigate("/teams")}>Teams</button>
        <button onClick={() => navigate("/following")}>Following</button>
        <form onSubmit={handleSearch}>
          <input type="search" value={searchTerm} onChange={(e) => setSearchTerm(e.target.value)} />
        </form>
        <h3>Follow some players!</h3>
      </div>
    );
  }

  return (
    <div className="home-background">
      <h1>Home</h1>
      <button onClick={() => navigate("/teams")}>Teams</button>
      <button onClick={() => navigate("/following")}>Following</button>
      <form onSubmit={handleSearch}>
        <input type="search" value={searchTerm} onChange={(e) => setSearchTerm(e.target.value)} />
      </form>
      <pre>
        <div className="players">
          {homeData.map((player) => (
            <div key={player.name} className="player-card">
              <div className="accent-bar" style={{ "--team-colour": teamColours[player.team_abbrev] } as React.CSSProperties}>
                {" "}
              </div>
              <h3 className="card-header">
                <span>{player.name}</span>
                <span>{player.team_abbrev}</span>
              </h3>
              <h4>
                <div className="status">
                  <span>{player.position}</span>
                  <span>#{player.sweater_number}</span>
                  <span className={`status-circle ${player.playing_today ? "playing" : "not-playing"}`} />
                </div>
              </h4>
              <div className="stats">
                <span>Games</span>
                <span>{player.games_played}</span>
                <span>Goals</span>
                <span>{player.goals}</span>
                <span>Assists</span>
                <span>{player.assists}</span>
                <span>Points</span>
                <span>{player.points}</span>
              </div>
              <div className="extra-stats">
                <span>P/PG</span>
                <span>{player["p/gp"]}</span>
                <span>Last 5</span>
                <span>{player.last_5_games_totals}</span>
              </div>
            </div>
          ))}
        </div>
        <table>
          <thead>
            <tr>
              <th>Team</th>
              <th>GP</th>
              <th>W</th>
              <th>L</th>
              <th>OTL</th>
              <th>PTS</th>
              <th>RW</th>
              <th>ROW</th>
              <th>GD</th>
              <th>L10</th>
            </tr>
          </thead>
          <tbody>
            {standings.map((team) => (
              <tr key={team.name} className={followedTeam === team.name ? "followed-team" : ""}>
                <td>{team.name}</td>
                <td>{team.games_played}</td>
                <td>{team.wins}</td>
                <td>{team.losses}</td>
                <td>{team.otl}</td>
                <td>{team.points}</td>
                <td>{team.rw}</td>
                <td>{team.row}</td>
                <td>{team.goal_differential}</td>
                <td>{team.last_10}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </pre>
    </div>
  );
}

function Search() {
  const [searchParams] = useSearchParams();
  const [results, setResults] = useState<searchedPlayer[]>([]);
  const [searchTerm, setSearchTerm] = useState("");
  const navigate = useNavigate();

  type searchedPlayer = {
    playerId: string;
    name: string;
    positionCode: string;
    teamAbbrev: string;
    height: string;
    weightInPounds: number;
    birthCountry: string;
    isFollowed: boolean;
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

  async function handleFollowClick(id: string) {
    await handleFollowPlayer(id);
    loadSearch();
  }

  async function handleUnfollowClick(id: string) {
    await handleUnfollowPlayer(id);
    loadSearch();
  }

  return (
    <div>
      <h1>Search Results</h1>
      <button onClick={() => navigate("/home")}>Home</button>
      <form onSubmit={handleSearch}>
        <input type="search" value={searchTerm} onChange={(e) => setSearchTerm(e.target.value)} />
      </form>
      <pre>
        {results.length === 0
          ? "No results found"
          : results.map((player) => (
              <div key={player.playerId}>
                <div>{player.name}</div>
                <div>{player.teamAbbrev}</div>
                <div>{player.positionCode}</div>
                <div>{player.height}</div>
                <div>{player.weightInPounds}</div>
                <div>{player.birthCountry}</div>
                <button onClick={() => (player.isFollowed ? handleUnfollowClick(player.playerId) : handleFollowClick(player.playerId))}>
                  {player.isFollowed ? "Unfollow" : "Follow"}
                </button>
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
      <Route path="/teams" element={<Teams />} />
    </Routes>
  );
}

export default App;
