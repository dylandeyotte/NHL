import "./App.css";
import { useEffect, useState } from "react";
import { teamList } from "./data/team_list";
import { fetchHelper } from "./helpers";
import { useNavigate } from "react-router-dom";
import { teamColours } from "./data/team_colours";

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

export function Home() {
  const navigate = useNavigate();
  const [searchTerm, setSearchTerm] = useState("");
  const [followedTeam, setFollowedTeam] = useState("");
  const [homeData, setHomeData] = useState<playerStats[]>([]);
  const [standings, setStandings] = useState<teamStandings[]>([]);

  async function loadHome() {
    // HTTP request
    try {
      const data = await fetchHelper("http://localhost:8080/api/home");
      setHomeData(data.players);
      setStandings(data.standings);
    } catch (err) {
      navigate("/");
    }
  }

  async function fetchTeam() {
    const data = await fetchHelper("http://localhost:8080/api/following");
    setFollowedTeam(data.team.TeamName);
  }

  // Run effect
  useEffect(() => {
    (loadHome(), fetchTeam());
  }, []);

  function handleSearch() {
    navigate(`/search?player=${encodeURIComponent(searchTerm)}`);
  }

  async function handleLogout() {
    localStorage.removeItem("token");
    localStorage.removeItem("refreshToken");
    navigate("/");
  }

  return (
    <div className="home-page">
      <div className="home-title">
        <h1>Home</h1>
        <button onClick={() => handleLogout()} className="logout-button">
          Logout
        </button>
      </div>
      <div className="button-row">
        <button onClick={() => navigate("/teams")} className="home-team-button">
          Teams
        </button>
        <button onClick={() => navigate("/following")} className="home-following-button">
          Following
        </button>
      </div>
      <form onSubmit={handleSearch}>
        <input type="search" value={searchTerm} onChange={(e) => setSearchTerm(e.target.value)} className="search-bar" />
      </form>
      <pre>
        <h3>{standings.length === 0 && homeData.length === 0 && "Follow some players!"}</h3>
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
        {followedTeam != "" && (
          <table>
            <thead>
              <tr>
                <th>{teamList.map((team) => team.name === followedTeam && `${team.division} Division`)}</th>
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
        )}
      </pre>
    </div>
  );
}
