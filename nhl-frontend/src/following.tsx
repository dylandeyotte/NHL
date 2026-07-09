import "./App.css";
import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { teamColours } from "./data/team_colours";
import { fetchHelper, handleUnfollowPlayer, handleUnfollowTeam, formatDate, formatHeight } from "./helpers";

type followee = {
  player_id: number;
  name: string;
  position: string;
  team_abbrev: string;
  birth_date: string;
  birth_city: string;
  birth_country: string;
  draft_year: number;
  draft_pos: number;
  weight: number;
  height: number;
};

type followedTeam = {
  TeamName: string;
  TriCode: string;
};

export function Following() {
  const [following, setFollowing] = useState<followee[]>([]);
  const [team, setTeam] = useState<followedTeam>();
  const navigate = useNavigate();

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
      <div className="following-container">
        <h1>Following</h1>
        <button onClick={() => navigate("/home")} className="home-button">
          Home
        </button>
        <pre>
          <div
            className={team?.TeamName === "" ? "" : "followed-team-box"}
            style={{ "--team-colour": team?.TriCode ? teamColours[team.TriCode] : "white" } as React.CSSProperties}
          >
            <h3>{team?.TeamName}</h3>
            <div className="unfollow-button">
              {team?.TeamName && (
                <button onClick={() => handleUnfollowTeamClick(team?.TriCode)}>
                  <span className="following">Following</span>
                  <span className="unfollow">Unfollow</span>
                </button>
              )}
            </div>
          </div>
        </pre>
      </div>
      <pre>
        <div className="followed-page">
          {following &&
            following.map((player) => (
              <div key={player.player_id} className="followed-player">
                <div className="accent-bar" style={{ "--team-colour": teamColours[player.team_abbrev] } as React.CSSProperties}>
                  {" "}
                </div>
                <h3 className="card-header">
                  <span>{player.name}</span>
                  <span>{player.team_abbrev}</span>
                </h3>
                <div className="card-bio">
                  {player.position} | {formatHeight(player.height)} | {player.weight} lbs
                </div>
                <div className="card-pob">
                  <span>{formatDate(player.birth_date)}</span>
                  <span className={player.birth_city.length > 17 ? "pob-long" : "pob"}>
                    {player.birth_city}, {player.birth_country === "CHE" ? "SUI" : player.birth_country}
                  </span>
                </div>
                <div className="draft-details">{player.draft_year === 0 ? "Undrafted" : `${player.draft_year} (#${player.draft_pos})`}</div>
                <div className="unfollow-button">
                  <button onClick={() => handleUnfollowPlayerClick(player.player_id)}>
                    <span className="following">Following</span>
                    <span className="unfollow">Unfollow</span>
                  </button>
                </div>
              </div>
            ))}
        </div>
      </pre>
      <p>{following.length} players followed</p>
    </div>
  );
}
