// Конфигурация
const API_BASE_URL = "http://localhost:8080";
const ENDPOINTS = {
  TRACKS: "/track/all",
  PLAY: "/track/play",
  PAUSE: "/track/pause",
  SELECT: "/track/select",
  SEEK: "/track/seek",
  PROFILE: "/profile",
};

// Состояние плеера
let tracks = [];
let currentTrackIndex = 0;
let isPlaying = false;
let isMuted = false;
let previousVolume = 0.7;

// DOM элементы
const audioPlayer = document.getElementById("audioPlayer");
const playIcon = document.getElementById("playIcon");
const playPauseBtn = document.getElementById("playPauseBtn");
const currentTrackTitle = document.getElementById("currentTrackTitle");
const currentTrackArtist = document.getElementById("currentTrackArtist");
const currentTimeEl = document.getElementById("currentTime");
const totalTimeEl = document.getElementById("totalTime");
const progressFill = document.getElementById("progressFill");
const progressBar = document.getElementById("progressBar");
const playlistEl = document.getElementById("playlist");
const statusText = document.getElementById("statusText");
const statusDot = document.getElementById("statusDot");
const tracksCountEl = document.getElementById("tracksCount");
const volumeSlider = document.getElementById("volumeSlider");
const body = document.getElementById("bg-body");

// Инициализация
document.addEventListener("DOMContentLoaded", () => {
  initPlayer();
  setupAudioListeners();
});

// Настройка слушателей аудио
function setupAudioListeners() {
  audioPlayer.addEventListener("timeupdate", updateProgress);
  audioPlayer.addEventListener("loadedmetadata", updateTotalTime);
  audioPlayer.addEventListener("ended", nextTrack);
  audioPlayer.addEventListener("play", () => {
    isPlaying = true;
    updatePlayButton();
  });
  audioPlayer.addEventListener("pause", () => {
    isPlaying = false;
    updatePlayButton();
  });
  audioPlayer.addEventListener("error", handleAudioError);
}

// Инициализация плеера
async function initPlayer() {
  updateStatus("connecting", "Подключение к серверу...");
  await loadTracks();
  checkServerHealth();
}

// Проверка здоровья сервера
async function checkServerHealth() {
  try {
    const response = await fetch(`${API_BASE_URL}/health`, { mode: "no-cors" });
    updateStatus("connected", "Сервер доступен");
  } catch (error) {
    updateStatus("error", "Сервер недоступен");
    showError("Не удалось подключиться к серверу. Проверьте localhost:8080");
  }
}

// Загрузка треков с сервера
async function loadTracks() {
  try {
    updateStatus("connecting", "Загрузка треков...");

    // Пытаемся загрузить с сервера
    const response = await fetch(`${API_BASE_URL}${ENDPOINTS.TRACKS}`);

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const data = await response.json();

    // Если данные пришли, используем их
    if (data && data.length > 0) {
      tracks = data;
    } else {
      // Если нет данных, используем тестовые
      useMockTracks();
    }
  } catch (error) {
    console.log("Failed to load tracks from server, using mock data:", error);
    useMockTracks();
  }

  renderPlaylist();
  updateTrackInfo();
  updateStatus("connected", `${tracks.length} треков загружено`);
  tracksCountEl.textContent = tracks.length;
}

// Тестовые данные на случай недоступности сервера
function useMockTracks() {
  tracks = [
    {
      id: 1,
      title: "Fox step",
      artist: "Polar Fox",
      duration: 124,
      cover: "fox",
      bgColor: "#FF6B6B",
    },
    {
      id: 2,
      title: "Bamboo dream",
      artist: "Panda Chill",
      duration: 108,
      cover: "tree",
      bgColor: "#4A90E2",
    },
    {
      id: 3,
      title: "Tropical beak",
      artist: "Toucan Groove",
      duration: 92,
      cover: "feather",
      bgColor: "#9B6B9C",
    },
    {
      id: 4,
      title: "Lazy breeze",
      artist: "Sloth vibes",
      duration: 142,
      cover: "clock",
      bgColor: "#6B8E5C",
    },
  ];
}

// Отрисовка плейлиста
function renderPlaylist() {
  if (!tracks || tracks.length === 0) {
    playlistEl.innerHTML =
      '<div class="error-message">Нет доступных треков</div>';
    return;
  }

  const icons = [
    "fa-fox",
    "fa-tree",
    "fa-feather",
    "fa-clock",
    "fa-music",
    "fa-headphones",
  ];

  playlistEl.innerHTML = tracks
    .map(
      (track, index) => `
        <div class="playlist-item ${index === currentTrackIndex ? "active" : ""}"
             data-index="${index}"
             onclick="selectTrack(${index})">
            <div class="item-cover">
                <i class="fas ${icons[index % icons.length]}"></i>
            </div>
            <div class="item-info">
                <div class="item-title">${track.title || "Без названия"}</div>
                <div class="item-sub">${track.artist || "Неизвестный исполнитель"}</div>
            </div>
            <span class="duration">${formatTime(track.duration || 0)}</span>
            <i class="fas fa-play play-icon"></i>
        </div>
    `,
    )
    .join("");
}

// Форматирование времени
function formatTime(seconds) {
  if (!seconds || isNaN(seconds)) return "0:00";
  const mins = Math.floor(seconds / 60);
  const secs = Math.floor(seconds % 60);
  return `${mins}:${secs < 10 ? "0" : ""}${secs}`;
}

// Обновление информации о текущем треке
function updateTrackInfo() {
  const track = tracks[currentTrackIndex];
  if (!track) return;

  currentTrackTitle.textContent = track.title || "Без названия";
  currentTrackArtist.textContent = track.artist || "Неизвестный исполнитель";
  totalTimeEl.textContent = formatTime(track.duration);

  // Меняем цвет фона
  if (track.bgColor) {
    body.style.background = `linear-gradient(145deg, ${track.bgColor} 0%, ${adjustColor(track.bgColor, 20)} 100%)`;
  }
}

// Осветление цвета
function adjustColor(hex, percent) {
  // Для простоты возвращаем тот же цвет
  return hex;
}

// Обновление кнопки воспроизведения
function updatePlayButton() {
  playIcon.className = isPlaying ? "fas fa-pause" : "fas fa-play";
}

// Обновление прогресса
function updateProgress() {
  if (!audioPlayer.duration) return;

  const current = audioPlayer.currentTime;
  const duration = audioPlayer.duration;
  const progress = (current / duration) * 100;

  progressFill.style.width = `${progress}%`;
  currentTimeEl.textContent = formatTime(current);
}

// Обновление общего времени
function updateTotalTime() {
  totalTimeEl.textContent = formatTime(audioPlayer.duration);
  if (tracks[currentTrackIndex]) {
    tracks[currentTrackIndex].duration = audioPlayer.duration;
  }
}

// Обработка ошибок аудио
function handleAudioError() {
  showNotification("❌ Ошибка воспроизведения");
  updateStatus("error", "Ошибка воспроизведения");
}

// Переключение воспроизведения
async function togglePlay() {
  if (!tracks.length) return;

  const track = tracks[currentTrackIndex];

  if (isPlaying) {
    audioPlayer.pause();
    await sendRequest(`${ENDPOINTS.PAUSE}`, "POST", { trackId: track.id });
    showNotification(`⏸ Пауза: ${track.title}`);
  } else {
    try {
      // Воспроизводим через эндпоинт /track/play/{id}
      const audioUrl = `${API_BASE_URL}${ENDPOINTS.PLAY}/${track.id}`;
      audioPlayer.src = audioUrl;
      await audioPlayer.play();

      // Отправляем статистику на сервер
      await sendRequest(`${ENDPOINTS.PLAY}/${track.id}`, "GET");

      showNotification(`▶ Играет: ${track.title}`);
    } catch (error) {
      console.error("Playback error:", error);
      showNotification("❌ Ошибка воспроизведения");
    }
  }
}

// Выбор трека
async function selectTrack(index) {
  if (index < 0 || index >= tracks.length) return;

  const wasPlaying = isPlaying;
  const previousTrack = currentTrackIndex;

  // Останавливаем текущее воспроизведение
  if (isPlaying) {
    audioPlayer.pause();
  }

  currentTrackIndex = index;

  // Отправляем информацию о выборе трека
  await sendRequest(ENDPOINTS.SELECT, "POST", {
    from: previousTrack,
    to: index,
    trackId: tracks[index].id,
    track: tracks[index],
  });

  // Обновляем UI
  updateTrackInfo();
  renderPlaylist();

  // Если играло, начинаем воспроизведение нового трека
  if (wasPlaying) {
    try {
      const audioUrl = `${API_BASE_URL}${ENDPOINTS.PLAY}/${tracks[index].id}`;
      audioPlayer.src = audioUrl;
      await audioPlayer.play();
    } catch (error) {
      console.error("Playback error:", error);
    }
  }

  showNotification(`🎵 Выбран: ${tracks[index].title}`);
}

// Предыдущий трек
function prevTrack() {
  if (!tracks.length) return;
  let newIndex = currentTrackIndex - 1;
  if (newIndex < 0) newIndex = tracks.length - 1;
  selectTrack(newIndex);
}

// Следующий трек
function nextTrack() {
  if (!tracks.length) return;
  let newIndex = (currentTrackIndex + 1) % tracks.length;
  selectTrack(newIndex);
}

// Обработка клика по прогресс-бару
async function handleProgressClick(event) {
  if (!audioPlayer.duration) return;

  const rect = progressBar.getBoundingClientRect();
  const clickX = event.clientX - rect.left;
  const width = rect.width;
  const clickPercentage = clickX / width;

  const newTime = clickPercentage * audioPlayer.duration;
  audioPlayer.currentTime = newTime;

  // Отправляем информацию о перемотке
  await sendRequest(ENDPOINTS.SEEK, "POST", {
    from: audioPlayer.currentTime,
    to: newTime,
    percentage: clickPercentage * 100,
    trackId: tracks[currentTrackIndex].id,
  });

  showNotification(`⏩ Перемотка к ${formatTime(newTime)}`);
}

// Установка громкости
function setVolume(value) {
  audioPlayer.volume = value;
  volumeSlider.value = value;
}

// Включение/выключение звука
function toggleMute() {
  if (isMuted) {
    audioPlayer.volume = previousVolume;
    volumeSlider.value = previousVolume;
    document.querySelector(".volume-control i").className = "fas fa-volume-up";
  } else {
    previousVolume = audioPlayer.volume;
    audioPlayer.volume = 0;
    volumeSlider.value = 0;
    document.querySelector(".volume-control i").className =
      "fas fa-volume-mute";
  }
  isMuted = !isMuted;
}

// Открытие профиля
function openProfile() {
  sendRequest(ENDPOINTS.PROFILE, "GET");
  showNotification("👤 Профиль");
}

// Отправка запроса на сервер
async function sendRequest(endpoint, method = "GET", data = null) {
  const url = `${API_BASE_URL}${endpoint}`;

  updateStatus("sending", "Отправка запроса...");

  const options = {
    method: method,
    headers: {
      "Content-Type": "application/json",
    },
  };

  if (data) {
    options.body = JSON.stringify(data);
  }

  try {
    const response = await fetch(url, options);

    if (response.ok) {
      updateStatus("connected", "Запрос выполнен");
    } else {
      updateStatus("error", `Ошибка ${response.status}`);
    }

    return response;
  } catch (error) {
    console.error("Request failed:", error);
    updateStatus("error", "Сервер недоступен");
    return null;
  }
}

// Обновление статуса
function updateStatus(type, message) {
  statusText.textContent = message;
  statusDot.className = "status-dot " + type;
}

// Показ уведомления
function showNotification(message) {
  const notification = document.createElement("div");
  notification.textContent = message;
  notification.style.cssText = `
        position: fixed;
        bottom: 20px;
        left: 50%;
        transform: translateX(-50%);
        background: rgba(0, 0, 0, 0.3);
        backdrop-filter: blur(8px);
        color: white;
        padding: 10px 20px;
        border-radius: 30px;
        font-size: 14px;
        z-index: 1000;
        animation: slideUp 0.3s ease, fadeOut 0.3s ease 1.7s forwards;
        border: 1px solid rgba(255, 255, 255, 0.3);
    `;

  document.body.appendChild(notification);

  setTimeout(() => {
    notification.remove();
  }, 2000);
}

// Показ ошибки
function showError(message) {
  const errorEl = document.createElement("div");
  errorEl.className = "error-message";
  errorEl.textContent = message;
  errorEl.style.cssText = `
        position: fixed;
        top: 20px;
        left: 50%;
        transform: translateX(-50%);
        background: rgba(244, 67, 54, 0.3);
        backdrop-filter: blur(8px);
        color: white;
        padding: 10px 20px;
        border-radius: 30px;
        font-size: 14px;
        z-index: 1000;
        border: 1px solid rgba(255, 255, 255, 0.3);
    `;

  document.body.appendChild(errorEl);

  setTimeout(() => {
    errorEl.remove();
  }, 3000);
}

// Добавляем стили для анимаций уведомлений
const style = document.createElement("style");
style.textContent = `
    @keyframes slideUp {
        from { opacity: 0; transform: translate(-50%, 20px); }
        to { opacity: 1; transform: translate(-50%, 0); }
    }
    @keyframes fadeOut {
        to { opacity: 0; transform: translate(-50%, -10px); }
    }
`;
document.head.appendChild(style);

// Глобальные функции для вызова из HTML
window.togglePlay = togglePlay;
window.prevTrack = prevTrack;
window.nextTrack = nextTrack;
window.selectTrack = selectTrack;
window.handleProgressClick = handleProgressClick;
window.setVolume = setVolume;
window.toggleMute = toggleMute;
window.openProfile = openProfile;
window.sendRequest = sendRequest;
