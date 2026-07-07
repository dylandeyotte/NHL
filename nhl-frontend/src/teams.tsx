import "./App.css";
import { teamList } from "./data/team_list";
import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { fetchHelper, handleFollowTeam, handleUnfollowTeam } from "./helpers";

export function Teams() {
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
              <img src={`https://assets.nhle.com/logos/nhl/svg/${team.tricode}_light.svg`} alt={team.name} className="team-logo" />
              <div className="team-name">{team.name}</div>
              <div className={followedTeam === team.tricode ? "unfollow-button" : "follow-button"}>
                <button
                  onClick={() => (followedTeam === team.tricode ? handleUnfollowClick(team.tricode) : handleFollowClick(team.tricode))}
                >
                  {followedTeam === team.tricode ? <span className="following">Following</span> : <span className="follow">Follow</span>}
                  <span className="unfollow">Unfollow</span>
                </button>
              </div>
            </div>
          ))}
        </div>
      </pre>
    </div>
  );
}
