let matches = [];
let currentMatchIndex = 0;
const ROTATION_INTERVAL = 6000; // 6 seconds per match

async function fetchMatches() {
    try {
        const response = await fetch('/api/laliga/schedule');
        const data = await response.json();
        matches = data.matches || [];
        
        if (matches.length > 0) {
            showCurrentMatch();
            if (matches.length > 1) {
                setInterval(rotateMatches, ROTATION_INTERVAL);
            }
        } else {
            showNoMatches();
        }
    } catch (error) {
        console.error('Error fetching matches:', error);
        showError();
    }
}

function showCurrentMatch() {
    const match = matches[currentMatchIndex];
    const matchDate = new Date(match.date);
    const dateStr = matchDate.toLocaleDateString('en-US', { 
        weekday: 'short', 
        month: 'short', 
        day: 'numeric' 
    });
    const timeStr = matchDate.toLocaleTimeString('en-US', { 
        hour: '2-digit', 
        minute: '2-digit' 
    });

    const content = document.getElementById('content');
    content.innerHTML = `
        <div class="match-container fade-in">
            <div class="team">
                <img src="${match.home_team.crest}" alt="${match.home_team.name}" class="team-crest">
                <div class="team-name">${match.home_team.name.replace(' CF', '').replace(' FC', '').replace(' SAD', '')}</div>
            </div>
            <div class="vs">VS</div>
            <div class="team">
                <img src="${match.away_team.crest}" alt="${match.away_team.name}" class="team-crest">
                <div class="team-name">${match.away_team.name.replace(' CF', '').replace(' FC', '').replace(' SAD', '')}</div>
            </div>
        </div>
        <div class="match-info">
            ${dateStr} • ${timeStr} • Matchday ${match.matchday}
        </div>
    `;
}

function rotateMatches() {
    const container = document.querySelector('.match-container');
    if (container) {
        container.classList.add('fade-out');
        
        setTimeout(() => {
            currentMatchIndex = (currentMatchIndex + 1) % matches.length;
            showCurrentMatch();
            
            setTimeout(() => {
                const newContainer = document.querySelector('.match-container');
                if (newContainer) {
                    newContainer.classList.remove('fade-out');
                    newContainer.classList.add('fade-in');
                }
            }, 50);
        }, 250);
    }
}

function showNoMatches() {
    document.getElementById('content').innerHTML = `
        <div class="no-matches">No upcoming matches this week</div>
    `;
}

function showError() {
    document.getElementById('content').innerHTML = `
        <div class="no-matches">Error loading matches</div>
    `;
}

// Start fetching matches when page loads
document.addEventListener('DOMContentLoaded', fetchMatches);