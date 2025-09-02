package main

import (
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUploadThumbnail(w http.ResponseWriter, r *http.Request) {
	videoIDString := r.PathValue("videoID")
	videoID, err := uuid.Parse(videoIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	fmt.Println("uploading thumbnail for video", videoID, "by user", userID)

	maxMemory := 10 << 20
	err = r.ParseMultipartForm(int64(maxMemory))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't parse multipart form", err)
		return
	}

	file, fileHandler, err := r.FormFile("thumbnail")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't get thumbnail from form", err)
		return
	}
	defer file.Close()

	videoRow, err := cfg.db.GetVideo(videoID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't get video from database", err)
		return
	}

	if videoRow.UserID != userID {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized access", err)
		return
	}

	mediaType, _, err := mime.ParseMediaType(fileHandler.Header.Get("Content-Type"))
	if mediaType != "image/jpeg" && mediaType != "image/png" {
		respondWithError(w, http.StatusBadRequest, "Couldn't parse non image file type", err)
		return
	}

	fileExtension := "." + strings.Split(mediaType, "/")[1]
	videoServerFilePath := filepath.Join(cfg.assetsRoot, videoIDString+fileExtension)
	serverVideoFile, err := os.Create(videoServerFilePath)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't create file on server", err)
		return
	}

	if _, err := io.Copy(serverVideoFile, file); err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't copy file to server", err)
		return
	}

	videoUrlFilePath := fmt.Sprintf(
		"http://localhost:%s/%s",
		cfg.port,
		videoServerFilePath,
	)
	videoRow.ThumbnailURL = &videoUrlFilePath

	log.Print("Debug:", videoUrlFilePath)

	if err := cfg.db.UpdateVideo(videoRow); err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't update video in database", err)
		return
	}

	respondWithJSON(w, http.StatusOK, videoRow)
}
