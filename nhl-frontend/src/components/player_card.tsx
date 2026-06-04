type Player = {
  name: string;
  goals: number;
  assists: number;
};

export function PlayerCard({ player }: { player: Player }) {
  return (
    <div>
      <h2>{player.name}</h2>
      <p>Goals: {player.goals}</p>
      <p>Assists: {player.assists}</p>
    </div>
  );
}
