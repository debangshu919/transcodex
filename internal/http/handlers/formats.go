package handlers

import (
	"encoding/json"
	"net/http"
)

type Format struct {
	Name         string   `json:"name"`
	Extensions   []string `json:"extensions"`
	MimeTypes    []string `json:"mime_types"`
	Category     string   `json:"category"`
	CanConvertTo []string `json:"can_convert_to"`
}

type FormatResponse struct {
	Name     string `json:"name"`
	Category string `json:"category"`
}

type FormatsResponse struct {
	Formats []FormatResponse `json:"formats"`
}

var formats = []Format{
	{
		Name:         "mp4",
		Extensions:   []string{".mp4"},
		MimeTypes:    []string{"video/mp4"},
		Category:     "video",
		CanConvertTo: []string{"mkv", "webm"},
	},
	{
		Name:         "mkv",
		Extensions:   []string{".mkv"},
		MimeTypes:    []string{"video/x-matroska"},
		Category:     "video",
		CanConvertTo: []string{"mp4", "webm"},
	},
	{
		Name:         "webm",
		Extensions:   []string{".webm"},
		MimeTypes:    []string{"video/webm"},
		Category:     "video",
		CanConvertTo: []string{"mp4", "mkv"},
	},
}

func ListFormats(w http.ResponseWriter, r *http.Request) {
	response := FormatsResponse{
		Formats: make([]FormatResponse, 0, len(formats)),
	}

	for _, format := range formats {
		response.Formats = append(response.Formats, FormatResponse{
			Name:     format.Name,
			Category: format.Category,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func GetFormat(w http.ResponseWriter, r *http.Request) {
	formatName := r.PathValue("format")
	for _, format := range formats {
		if format.Name == formatName {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(format)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}
