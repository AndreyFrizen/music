package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID                uuid.UUID `json:"-" db:"id"`
	Username          string    `json:"username" db:"username"`
	Password          string    `json:"password,omitempty"`
	EncryptedPassword string    `json:"encrypted_password" db:"password"`
	Email             string    `json:"email" db:"email"`
}

func (u *User) EncryptPassword() error {

	if len(u.Password) > 0 {
		enc, err := encryptedPassword(u.Password)
		if err != nil {
			return err
		}
		u.EncryptedPassword = string(enc)
	}

	return nil
}

func encryptedPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

type Artist struct {
	ID   uuid.UUID `json:"id" db:"id"`
	Name string    `json:"name" db:"name"`
}

type Album struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	ArtistID    uuid.UUID `json:"artist_id" db:"artist_id"`
	ReleaseDate time.Time `json:"release_date" db:"release_date"`
}

type Track struct {
	ID       uuid.UUID `json:"id" db:"id"`
	Title    string    `json:"title" db:"title"`
	Duration time.Time `json:"duration" db:"duration"` // в секундах
	AudioURL string    `json:"audio_url" db:"audio_url"`
}

type Playlist struct {
	ID     uuid.UUID `json:"id" db:"id"`
	UserID uuid.UUID `json:"user_id" db:"user_id"`
	Title  string    `json:"title" db:"title"`
}

type PlaylistTrack struct {
	PlaylistID uuid.UUID `json:"playlist_id" db:"playlist_id"`
	TrackID    uuid.UUID `json:"track_id" db:"track_id"`
	Position   int       `json:"position" db:"position"`
}

func main() {

	// Initialize database connection
	db, err := sql.Open("sqlite3", "/home/andrey/golang-proj/music/musicdrevier/storage/storage.db")
	if err != nil {
		log.Fatal("error opening database", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("error pinging database", err)
	}

	// Initialize router connection
	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.POST("/register", func(c *gin.Context) {
		RegisterUser(c, db)
	})
	r.POST("/auth", func(c *gin.Context) {
		AuthUser(c, db)
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// CreateUser(db)
	// AddArtist(db)
	// CreateAlbum(db)
	// CreatePlaylist(db)
	// CreatePlaylistTrack(db)

	log.Println("Database created successfully")
	r.Run("localhost:8080")
}

// RegisterUser ...

func RegisterUser(c *gin.Context, db *sql.DB) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := CreateUser(db, &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func CreateUser(db *sql.DB, u *User) error {
	u.EncryptPassword()
	query := fmt.Sprintf("INSERT INTO users VALUES ('%s', '%s', '%s', '%s')",
		uuid.New().String(), u.Username, u.EncryptedPassword, u.Email)

	_, err := db.Exec(query)

	if err != nil {
		return err
	}

	return nil
}

// Auth ...

func AuthUser(c *gin.Context, db *sql.DB) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := AuthenticateUser(db, &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Auth user successfully"})
}

func AuthenticateUser(db *sql.DB, u *User) error {
	query := fmt.Sprintf("SELECT password FROM users WHERE email = '%s'", u.Email)

	rows, err := db.Query(query)
	if err != nil {
		return err
	}

	defer rows.Close()

	if rows.Next() {
		var Password string

		if err := rows.Scan(&Password); err != nil {
			return err
		}

		if err := bcrypt.CompareHashAndPassword([]byte(Password), []byte(u.Password)); err != nil {
			return errors.New("invalid password")
		}

		return nil
	}

	return errors.New("error")
}

// func AddArtist(db *sql.DB) error {
// 	query := fmt.Sprintf("INSERT INTO artists VALUES ('%s', '%s')", uuid.New().String(), "John Doe")

// 	log.Print(query)

// 	_, err := db.Exec(query)

// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func CreateAlbum(db *sql.DB) error {
// 	query := fmt.Sprintf("INSERT INTO albums VALUES ('%s', '%s', '%s')", uuid.New().String(), "Album Title", time.Now().Format(time.RFC3339))

// 	log.Print(query)

// 	_, err := db.Exec(query)

// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func CreatePlaylist(db *sql.DB) error {
// 	query := fmt.Sprintf("INSERT INTO playlists VALUES ('%s', '%s', '%s')", uuid.New().String(), "Playlist Title", time.Now().Format(time.RFC3339))

// 	log.Print(query)

// 	_, err := db.Exec(query)

// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func CreatePlaylistTrack(db *sql.DB) error {
// 	query := fmt.Sprintf("INSERT INTO playlist_tracks VALUES ('%s', '%s', '%s')", uuid.New().String(), "Playlist ID", "Track ID")

// 	log.Print(query)

// 	_, err := db.Exec(query)

// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func streamTrack(w http.ResponseWriter, r *http.Request) {
// 	trackID := gin.URLParam(r, "id")
// 	userID := r.Context().Value("userID").(string)

// 	track, err := h.trackService.StreamTrack(r.Context(), userID, trackID)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// В реальном сервисе здесь был бы стриминг файла
// 	json.NewEncoder(w).Encode(track)
// }

// func (bd *sql.DB) FindByID(ctx context.Context, id uuid.UUID) (*domain.Track, error) {
// 	query := `
//         SELECT id, title, album_id, duration, track_number, audio_url, plays_count
//         FROM tracks WHERE id = $1
//     `

// 	var track domain.Track
// 	err := r.db.QueryRowContext(ctx, query, id).Scan(
// 		&track.ID, &track.Title, &track.AlbumID, &track.Duration,
// 		&track.TrackNumber, &track.AudioURL, &track.PlaysCount,
// 	)

// 	if err != nil {
// 		return nil, err
// 	}

// 	return &track, nil
// }

// func (s *TrackService) StreamTrack(ctx context.Context, userID, trackID string) (*domain.Track, error) {
// 	// Проверяем кэш
// 	cacheKey := fmt.Sprintf("track:%s", trackID)
// 	cachedTrack, err := s.redis.Get(ctx, cacheKey)
// 	if err == nil && cachedTrack != "" {
// 		// Возвращаем из кэша
// 	}

// 	// Получаем трек из БД
// 	trackUUID, _ := uuid.Parse(trackID)
// 	track, err := s.trackRepo.FindByID(ctx, trackUUID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Инкрементируем счетчик прослушиваний
// 	go func() {
// 		s.trackRepo.IncrementPlays(context.Background(), trackUUID)
// 	}()

// 	// Сохраняем в историю
// 	history := &domain.StreamingHistory{
// 		ID:       uuid.New(),
// 		UserID:   userUUID,
// 		TrackID:  trackUUID,
// 		PlayedAt: time.Now(),
// 		Duration: 0, // обновляется при завершении прослушивания
// 	}

// 	// Сохраняем в кэш
// 	s.redis.Set(ctx, cacheKey, track, 10*time.Minute)

// 	return track, nil
// }
