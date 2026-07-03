export async function handleFollowPlayer(id: string) {
  const data = await fetchHelper(`http://localhost:8080/api/players/${id}/follow`, "POST");

  console.log(data);
  console.log(`Player followed: ${data.PlayerID}`);
}

export async function handleUnfollowPlayer(id: string) {
  const data = await fetchHelper(`http://localhost:8080/api/players/${id}/follow`, "DELETE");

  console.log(data);
  console.log(`Player unfollowed: ${id}`);
}

export async function handleFollowTeam(tricode: string) {
  const data = await fetchHelper(`http://localhost:8080/api/teams/${tricode}/follow`, "POST");

  console.log(data);
  console.log(`Team followd: ${data.TriCode}`);
}

export async function handleUnfollowTeam(tricode: string) {
  const data = await fetchHelper(`http://localhost:8080/api/teams/${tricode}/follow`, "DELETE");

  console.log(data);
  console.log(`Team unfollowd: ${tricode}`);
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

export async function fetchHelper(url: string, options?: string) {
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

export function formatDate(date: string) {
  const [year, month, day] = date.split("-").map(Number);

  return new Date(year, month - 1, day).toLocaleDateString("en-CA", {
    year: "numeric",
    month: "short",
    day: "2-digit",
  });
}

export function formatHeight(height: number) {
  const ft = Math.floor(height / 12);
  const inch = height % 12;

  return `${ft}'${inch}"`;
}
