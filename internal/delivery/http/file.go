package httpdelivery

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	fileapp "github.com/xcus33me/stash/internal/application/file"
)

const maxMultipartMemory = 32 << 20

type FileHandler struct {
	uc  *fileapp.Usecase
	log *slog.Logger
}

func NewFileHandler(uc *fileapp.Usecase, log *slog.Logger) *FileHandler {
	return &FileHandler{uc: uc, log: log.With("component", "file_handler")}
}

type fileView struct {
	ID            uuid.UUID  `json:"id"`
	OriginalName  string     `json:"original_name"`
	UploaderAlias *string    `json:"uploader_alias,omitempty"`
	SizeBytes     int64      `json:"size_bytes"`
	MimeType      string     `json:"mime_type"`
	UploadedAt    time.Time  `json:"uploaded_at"`
	ExpiresAt     *time.Time `json:"expires_at,omitempty"`
}

type listFilesResponse struct {
	Files []fileView `json:"files"`
}

type uploadResponse struct {
	FileID    uuid.UUID  `json:"file_id"`
	SHA256    string     `json:"sha256"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

type downloadResponse struct {
	URL          string     `json:"url"`
	OriginalName string     `json:"original_name"`
	MimeType     string     `json:"mime_type"`
	SizeBytes    int64      `json:"size_bytes"`
	ExpiresAt    *time.Time `json:"expires_at,omitempty"`
}

func (h *FileHandler) List(w http.ResponseWriter, r *http.Request) {
	roomID, err := uuid.Parse(chi.URLParam(r, "roomID"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid room id"})
		return
	}

	out, err := h.uc.ListFiles(r.Context(), fileapp.ListFilesInput{RoomID: roomID})
	if err != nil {
		writeError(w, h.log, err)
		return
	}

	views := make([]fileView, 0, len(out.Files))
	for _, f := range out.Files {
		views = append(views, fileView{
			ID:            f.ID,
			OriginalName:  f.OriginalName,
			UploaderAlias: f.UploaderAlias,
			SizeBytes:     f.SizeBytes,
			MimeType:      f.MimeType,
			UploadedAt:    f.UploadedAt,
			ExpiresAt:     f.ExpiresAt,
		})
	}

	writeJSON(w, http.StatusOK, listFilesResponse{Files: views})
}

func (h *FileHandler) Upload(w http.ResponseWriter, r *http.Request) {
	roomID, err := uuid.Parse(chi.URLParam(r, "roomID"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid room id"})
		return
	}

	if err := r.ParseMultipartForm(maxMultipartMemory); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid multipart"})
		return
	}

	f, header, err := r.FormFile("file")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "file field required"})
		return
	}
	defer f.Close()

	var alias *string
	if a := r.FormValue("uploader_alias"); a != "" {
		alias = &a
	}

	out, err := h.uc.UploadFile(r.Context(), fileapp.UploadFileInput{
		RoomID:        roomID,
		OriginalName:  header.Filename,
		UploaderAlias: alias,
		MimeType:      header.Header.Get("Content-Type"),
		SizeBytes:     header.Size,
		Reader:        f,
	})
	if err != nil {
		writeError(w, h.log, err)
		return
	}

	writeJSON(w, http.StatusCreated, uploadResponse{
		FileID:    out.FileID,
		SHA256:    out.SHA256,
		ExpiresAt: out.ExpiresAt,
	})
}

func (h *FileHandler) Download(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "fileID"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid file id"})
		return
	}

	out, err := h.uc.DownloadFile(r.Context(), fileapp.DownloadFileInput{FileID: id})
	if err != nil {
		writeError(w, h.log, err)
		return
	}

	writeJSON(w, http.StatusOK, downloadResponse{
		URL:          out.URL,
		OriginalName: out.OriginalName,
		MimeType:     out.MimeType,
		SizeBytes:    out.SizeBytes,
		ExpiresAt:    out.ExpiresAt,
	})
}

func (h *FileHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "fileID"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid file id"})
		return
	}
	token := r.Header.Get(ownerTokenHeader)
	if token == "" {
		writeJSON(w, http.StatusUnauthorized, errorResponse{Error: "missing owner token"})
		return
	}

	if err := h.uc.DeleteFile(r.Context(), fileapp.DeleteFileInput{
		FileID:     id,
		OwnerToken: token,
	}); err != nil {
		writeError(w, h.log, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
