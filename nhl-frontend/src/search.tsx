import "./App.css";
import { useEffect, useState } from "react";
import { teamColours } from "./data/team_colours";
import { useNavigate, useSearchParams } from "react-router-dom";
import { fetchHelper, handleUnfollowPlayer, handleFollowPlayer } from "./helpers";

type searchedPlayer = {
  playerId: string;
  name: string;
  positionCode: string;
  teamAbbrev: string;
  height: string;
  weightInPounds: number;
  birthCity: string;
  birthCountry: string;
  isFollowed: boolean;
};

export function Search() {
  const [searchParams] = useSearchParams();
  const [results, setResults] = useState<searchedPlayer[]>([]);
  const [searchTerm, setSearchTerm] = useState("");
  const navigate = useNavigate();

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
        <div className="search-page">
          {results.length === 0
            ? "No results found"
            : results.map((player) => (
                <div key={player.playerId} className="search-cards">
                  <div className="accent-bar" style={{ "--team-colour": teamColours[player.teamAbbrev] } as React.CSSProperties}>
                    {" "}
                  </div>
                  <div className="search-card-header">
                    <span>{player.name}</span>
                    <span>{player.teamAbbrev}</span>
                  </div>
                  <div className="card-pob">
                    {player.birthCity}, {player.birthCountry === "CHE" ? "SUI" : player.birthCountry}
                  </div>
                  <div className="card-bio">
                    {player.positionCode} | {player.height} | {player.weightInPounds} lbs
                    <div className="search-follow-button">
                      <button
                        onClick={() => (player.isFollowed ? handleUnfollowClick(player.playerId) : handleFollowClick(player.playerId))}
                      >
                        {player.isFollowed ? "Unfollow" : "Follow"}
                      </button>
                    </div>
                  </div>
                </div>
              ))}
        </div>
      </pre>
    </div>
  );
}
