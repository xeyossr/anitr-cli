package favorites

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// FavoriteAnime represents a favorite anime entry
type FavoriteAnime struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	PosterURL   string    `json:"poster_url"`
	Type        string    `json:"type"`
	AddedAt     time.Time `json:"added_at"`
	LastWatched time.Time `json:"last_watched,omitempty"`
}

// FavoritesManager handles favorite anime operations
type FavoritesManager struct {
	configDir string
	filePath  string
}

// NewFavoritesManager creates a new favorites manager
func NewFavoritesManager() (*FavoritesManager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("kullanıcı ana dizini alınamadı: %w", err)
	}

	configDir := filepath.Join(homeDir, ".config", "anitr-cli")
	filePath := filepath.Join(configDir, "favorites.json")

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("config dizini oluşturulamadı: %w", err)
	}

	return &FavoritesManager{
		configDir: configDir,
		filePath:  filePath,
	}, nil
}

// LoadFavorites loads favorites from file
func (fm *FavoritesManager) LoadFavorites() ([]FavoriteAnime, error) {
	if _, err := os.Stat(fm.filePath); os.IsNotExist(err) {
		return []FavoriteAnime{}, nil
	}

	data, err := os.ReadFile(fm.filePath)
	if err != nil {
		return nil, fmt.Errorf("favoriler dosyası okunamadı: %w", err)
	}

	var favorites []FavoriteAnime
	if err := json.Unmarshal(data, &favorites); err != nil {
		return nil, fmt.Errorf("favoriler parse edilemedi: %w", err)
	}

	return favorites, nil
}

// SaveFavorites saves favorites to file
func (fm *FavoritesManager) SaveFavorites(favorites []FavoriteAnime) error {
	data, err := json.MarshalIndent(favorites, "", "  ")
	if err != nil {
		return fmt.Errorf("favoriler serialize edilemedi: %w", err)
	}

	if err := os.WriteFile(fm.filePath, data, 0644); err != nil {
		return fmt.Errorf("favoriler dosyası yazılamadı: %w", err)
	}

	return nil
}

// AddFavorite adds an anime to favorites
func (fm *FavoritesManager) AddFavorite(id int, name, posterURL, animeType string) error {
	favorites, err := fm.LoadFavorites()
	if err != nil {
		return err
	}

	// Check if already exists
	for _, fav := range favorites {
		if fav.ID == id {
			return fmt.Errorf("anime zaten favorilerde")
		}
	}

	newFavorite := FavoriteAnime{
		ID:        id,
		Name:      name,
		PosterURL: posterURL,
		Type:      animeType,
		AddedAt:   time.Now(),
	}

	favorites = append(favorites, newFavorite)
	return fm.SaveFavorites(favorites)
}

// RemoveFavorite removes an anime from favorites
func (fm *FavoritesManager) RemoveFavorite(id int) error {
	favorites, err := fm.LoadFavorites()
	if err != nil {
		return err
	}

	for i, fav := range favorites {
		if fav.ID == id {
			favorites = append(favorites[:i], favorites[i+1:]...)
			return fm.SaveFavorites(favorites)
		}
	}

	return fmt.Errorf("anime favorilerde bulunamadı")
}

// IsFavorite checks if an anime is in favorites
func (fm *FavoritesManager) IsFavorite(id int) (bool, error) {
	favorites, err := fm.LoadFavorites()
	if err != nil {
		return false, err
	}

	for _, fav := range favorites {
		if fav.ID == id {
			return true, nil
		}
	}

	return false, nil
}

// UpdateLastWatched updates the last watched time for a favorite
func (fm *FavoritesManager) UpdateLastWatched(id int) error {
	favorites, err := fm.LoadFavorites()
	if err != nil {
		return err
	}

	for i, fav := range favorites {
		if fav.ID == id {
			favorites[i].LastWatched = time.Now()
			return fm.SaveFavorites(favorites)
		}
	}

	return nil // Not an error if not in favorites
}

// GetFavoriteNames returns a slice of favorite anime names for UI
func (fm *FavoritesManager) GetFavoriteNames() ([]string, []int, error) {
	favorites, err := fm.LoadFavorites()
	if err != nil {
		return nil, nil, err
	}

	names := make([]string, len(favorites))
	ids := make([]int, len(favorites))

	for i, fav := range favorites {
		names[i] = fmt.Sprintf("%s (ID: %d)", fav.Name, fav.ID)
		ids[i] = fav.ID
	}

	return names, ids, nil
}