package history

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// WatchHistory represents a watch history entry
type WatchHistory struct {
	AnimeID       int       `json:"anime_id"`
	AnimeName     string    `json:"anime_name"`
	EpisodeIndex  int       `json:"episode_index"`
	EpisodeName   string    `json:"episode_name"`
	SeasonIndex   int       `json:"season_index"`
	WatchedAt     time.Time `json:"watched_at"`
	Duration      int       `json:"duration,omitempty"` // in seconds
	Progress      float64   `json:"progress,omitempty"` // 0.0 to 1.0
	Completed     bool      `json:"completed"`
	PosterURL     string    `json:"poster_url,omitempty"`
}

// HistoryManager handles watch history operations
type HistoryManager struct {
	configDir string
	filePath  string
}

// NewHistoryManager creates a new history manager
func NewHistoryManager() (*HistoryManager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("kullanıcı ana dizini alınamadı: %w", err)
	}

	configDir := filepath.Join(homeDir, ".config", "anitr-cli")
	filePath := filepath.Join(configDir, "history.json")

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("config dizini oluşturulamadı: %w", err)
	}

	return &HistoryManager{
		configDir: configDir,
		filePath:  filePath,
	}, nil
}

// LoadHistory loads watch history from file
func (hm *HistoryManager) LoadHistory() ([]WatchHistory, error) {
	if _, err := os.Stat(hm.filePath); os.IsNotExist(err) {
		return []WatchHistory{}, nil
	}

	data, err := os.ReadFile(hm.filePath)
	if err != nil {
		return nil, fmt.Errorf("geçmiş dosyası okunamadı: %w", err)
	}

	var history []WatchHistory
	if err := json.Unmarshal(data, &history); err != nil {
		return nil, fmt.Errorf("geçmiş parse edilemedi: %w", err)
	}

	return history, nil
}

// SaveHistory saves watch history to file
func (hm *HistoryManager) SaveHistory(history []WatchHistory) error {
	data, err := json.MarshalIndent(history, "", "  ")
	if err != nil {
		return fmt.Errorf("geçmiş serialize edilemedi: %w", err)
	}

	if err := os.WriteFile(hm.filePath, data, 0644); err != nil {
		return fmt.Errorf("geçmiş dosyası yazılamadı: %w", err)
	}

	return nil
}

// AddWatchEntry adds a new watch entry to history
func (hm *HistoryManager) AddWatchEntry(animeID int, animeName string, episodeIndex int, episodeName string, seasonIndex int, posterURL string) error {
	history, err := hm.LoadHistory()
	if err != nil {
		return err
	}

	// Remove existing entry for the same episode if exists
	for i := len(history) - 1; i >= 0; i-- {
		if history[i].AnimeID == animeID && history[i].EpisodeIndex == episodeIndex {
			history = append(history[:i], history[i+1:]...)
			break
		}
	}

	newEntry := WatchHistory{
		AnimeID:      animeID,
		AnimeName:    animeName,
		EpisodeIndex: episodeIndex,
		EpisodeName:  episodeName,
		SeasonIndex:  seasonIndex,
		WatchedAt:    time.Now(),
		Completed:    true,
		PosterURL:    posterURL,
	}

	// Add to the beginning of the slice (most recent first)
	history = append([]WatchHistory{newEntry}, history...)

	// Keep only last 100 entries
	if len(history) > 100 {
		history = history[:100]
	}

	return hm.SaveHistory(history)
}

// GetLastWatchedEpisode returns the last watched episode for an anime
func (hm *HistoryManager) GetLastWatchedEpisode(animeID int) (*WatchHistory, error) {
	history, err := hm.LoadHistory()
	if err != nil {
		return nil, err
	}

	for _, entry := range history {
		if entry.AnimeID == animeID {
			return &entry, nil
		}
	}

	return nil, nil // No history found
}

// GetRecentHistory returns recent watch history (last N entries)
func (hm *HistoryManager) GetRecentHistory(limit int) ([]WatchHistory, error) {
	history, err := hm.LoadHistory()
	if err != nil {
		return nil, err
	}

	if len(history) <= limit {
		return history, nil
	}

	return history[:limit], nil
}

// GetHistoryForAnime returns all watch history for a specific anime
func (hm *HistoryManager) GetHistoryForAnime(animeID int) ([]WatchHistory, error) {
	history, err := hm.LoadHistory()
	if err != nil {
		return nil, err
	}

	var animeHistory []WatchHistory
	for _, entry := range history {
		if entry.AnimeID == animeID {
			animeHistory = append(animeHistory, entry)
		}
	}

	return animeHistory, nil
}

// ClearHistory clears all watch history
func (hm *HistoryManager) ClearHistory() error {
	return hm.SaveHistory([]WatchHistory{})
}

// GetHistoryNames returns formatted history entries for UI
func (hm *HistoryManager) GetHistoryNames(limit int) ([]string, []WatchHistory, error) {
	history, err := hm.GetRecentHistory(limit)
	if err != nil {
		return nil, nil, err
	}

	names := make([]string, len(history))
	for i, entry := range history {
		names[i] = fmt.Sprintf("%s - %s (%s)", 
			entry.AnimeName, 
			entry.EpisodeName, 
			entry.WatchedAt.Format("02/01/2006 15:04"))
	}

	return names, history, nil
}

// UpdateProgress updates the progress of a watch entry
func (hm *HistoryManager) UpdateProgress(animeID, episodeIndex int, progress float64, duration int) error {
	history, err := hm.LoadHistory()
	if err != nil {
		return err
	}

	for i, entry := range history {
		if entry.AnimeID == animeID && entry.EpisodeIndex == episodeIndex {
			history[i].Progress = progress
			history[i].Duration = duration
			history[i].Completed = progress >= 0.9 // Consider 90%+ as completed
			return hm.SaveHistory(history)
		}
	}

	return nil // Entry not found, not an error
}