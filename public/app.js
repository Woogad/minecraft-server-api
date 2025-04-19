$(async function(){
    await loadStatus();
});

async function loadStatus() {
  await fetch("/api/status")
    .then(res => res.json())
    .then(data => {
      const statusDiv = document.getElementById("status");
      let playerListHTML = "No players online";

      if (data.player_list && data.player_list.length > 0) {
        playerListHTML = `
          <div class="player-list">
            <strong>Players Online:</strong>
            <ul>
              ${data.player_list.map(name => `<li>${name}</li>`).join("")}
            </ul>
          </div>
        `;
      }

      statusDiv.innerHTML = `
        <p>Status: <span class="${data.online ? 'online' : 'offline'}">${data.online ? 'online' : 'offline'}</span></p>
        <p>Host: ${data.host}</p>
        <p>Version: ${data.version || 'N/A'}</p>
        <p>Players: ${data.player_count} / ${data.max_players}</p>
        ${playerListHTML}
      `;
    })
    .catch(err => {
      document.getElementById("status").innerText = "Error fetching server status.";
      console.error(err);
    });
}