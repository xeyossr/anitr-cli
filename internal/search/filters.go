package search

import (
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/xeyossr/anitr-cli/internal"
)

// SearchFilters represents search filter options
type SearchFilters struct {
	Genre      string
	Year       int
	Status     string // "ongoing", "completed", "upcoming"
	Type       string // "tv", "movie", "ova", "special"
	SortBy     string // "name", "year", "rating", "popularity"
	SortOrder  string // "asc", "desc"
	MinRating  float64
	MaxRating  float64
}

// SearchResult represents an enhanced search result
type SearchResult struct {
	ID           int                    `json:"id"`
	Name         string                 `json:"name"`
	OriginalName string                 `json:"original_title"`
	Type         string                 `json:"title_type"`
	Poster       string                 `json:"poster"`
	Year         int                    `json:"year,omitempty"`
	Rating       float64                `json:"rating,omitempty"`
	Genres       []string               `json:"genres,omitempty"`
	Status       string                 `json:"status,omitempty"`
	Description  string                 `json:"description,omitempty"`
	EpisodeCount int                    `json:"episode_count,omitempty"`
	RawData      map[string]interface{} `json:"-"`
}

// FilterManager handles search filtering operations
type FilterManager struct {
	filters SearchFilters
}

// NewFilterManager creates a new filter manager
func NewFilterManager() *FilterManager {
	return &FilterManager{
		filters: SearchFilters{
			SortBy:    "name",
			SortOrder: "asc",
			MinRating: 0.0,
			MaxRating: 10.0,
		},
	}
}

// SetFilters sets the search filters
func (fm *FilterManager) SetFilters(filters SearchFilters) {
	fm.filters = filters
}

// GetFilters returns current filters
func (fm *FilterManager) GetFilters() SearchFilters {
	return fm.filters
}

// ConvertToSearchResults converts raw API data to SearchResult structs
func (fm *FilterManager) ConvertToSearchResults(rawData []map[string]interface{}) []SearchResult {
	results := make([]SearchResult, 0, len(rawData))

	for _, item := range rawData {
		result := SearchResult{
			ID:           int(internal.GetFloat64(item, "id")),
			Name:         internal.GetString(item, "name"),
			OriginalName: internal.GetString(item, "original_title"),
			Type:         internal.GetString(item, "title_type"),
			Poster:       internal.GetString(item, "poster"),
			RawData:      item,
		}

		// Extract additional fields if available
		if year := internal.GetString(item, "year"); year != "" {
			if y, err := strconv.Atoi(year); err == nil {
				result.Year = y
			}
		}

		if rating := internal.GetFloat64(item, "rating"); rating > 0 {
			result.Rating = rating
		}

		if genres := internal.GetString(item, "genres"); genres != "" {
			result.Genres = strings.Split(genres, ",")
			for i := range result.Genres {
				result.Genres[i] = strings.TrimSpace(result.Genres[i])
			}
		}

		result.Status = internal.GetString(item, "status")
		result.Description = internal.GetString(item, "description")
		result.EpisodeCount = int(internal.GetFloat64(item, "episode_count"))

		results = append(results, result)
	}

	return results
}

// ApplyFilters applies the current filters to search results
func (fm *FilterManager) ApplyFilters(results []SearchResult) []SearchResult {
	filtered := make([]SearchResult, 0, len(results))

	for _, result := range results {
		// Apply type filter
		if fm.filters.Type != "" && !strings.EqualFold(result.Type, fm.filters.Type) {
			continue
		}

		// Apply year filter
		if fm.filters.Year > 0 && result.Year != fm.filters.Year {
			continue
		}

		// Apply status filter
		if fm.filters.Status != "" && !strings.EqualFold(result.Status, fm.filters.Status) {
			continue
		}

		// Apply genre filter
		if fm.filters.Genre != "" {
			genreMatch := false
			for _, genre := range result.Genres {
				if strings.Contains(strings.ToLower(genre), strings.ToLower(fm.filters.Genre)) {
					genreMatch = true
					break
				}
			}
			if !genreMatch {
				continue
			}
		}

		// Apply rating filter
		if result.Rating > 0 && (result.Rating < fm.filters.MinRating || result.Rating > fm.filters.MaxRating) {
			continue
		}

		filtered = append(filtered, result)
	}

	// Apply sorting
	fm.sortResults(filtered)

	return filtered
}

// sortResults sorts the results based on current sort settings
func (fm *FilterManager) sortResults(results []SearchResult) {
	sort.Slice(results, func(i, j int) bool {
		var less bool

		switch fm.filters.SortBy {
		case "name":
			less = strings.ToLower(results[i].Name) < strings.ToLower(results[j].Name)
		case "year":
			less = results[i].Year < results[j].Year
		case "rating":
			less = results[i].Rating < results[j].Rating
		case "popularity":
			// For now, sort by ID (assuming lower ID = more popular)
			less = results[i].ID < results[j].ID
		default:
			less = strings.ToLower(results[i].Name) < strings.ToLower(results[j].Name)
		}

		if fm.filters.SortOrder == "desc" {
			return !less
		}
		return less
	})
}

// GetAvailableGenres returns a list of common anime genres
func (fm *FilterManager) GetAvailableGenres() []string {
	return []string{
		"Aksiyon", "Macera", "Komedi", "Drama", "Fantastik",
		"Korku", "Gizem", "Romantik", "Bilim Kurgu", "Slice of Life",
		"Spor", "Supernatural", "Gerilim", "Savaş", "Mecha",
		"Müzik", "Okul", "Harem", "Josei", "Seinen",
		"Shoujo", "Shounen", "Yaoi", "Yuri", "Ecchi",
	}
}

// GetAvailableYears returns a list of years for filtering
func (fm *FilterManager) GetAvailableYears() []int {
	currentYear := time.Now().Year()
	years := make([]int, 0, currentYear-1960+1)

	for year := currentYear; year >= 1960; year-- {
		years = append(years, year)
	}

	return years
}

// GetAvailableTypes returns available anime types
func (fm *FilterManager) GetAvailableTypes() []string {
	return []string{"TV", "Movie", "OVA", "Special", "ONA"}
}

// GetAvailableStatuses returns available anime statuses
func (fm *FilterManager) GetAvailableStatuses() []string {
	return []string{"Devam Ediyor", "Tamamlandı", "Yakında"}
}

// GetSortOptions returns available sort options
func (fm *FilterManager) GetSortOptions() []string {
	return []string{"İsim", "Yıl", "Puan", "Popülerlik"}
}

// ResetFilters resets all filters to default values
func (fm *FilterManager) ResetFilters() {
	fm.filters = SearchFilters{
		SortBy:    "name",
		SortOrder: "asc",
		MinRating: 0.0,
		MaxRating: 10.0,
	}
}