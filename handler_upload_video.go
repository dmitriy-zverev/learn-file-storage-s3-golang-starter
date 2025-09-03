package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUploadVideo(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<30)

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

	videoRow, err := cfg.db.GetVideo(videoID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't get video from database", err)
		return
	}

	if videoRow.UserID != userID {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized access", err)
		return
	}

	file, fileHandler, err := r.FormFile("video")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't parse file", err)
		return
	}
	defer file.Close()

	mediaType, _, err := mime.ParseMediaType(fileHandler.Header.Get("Content-Type"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't get media type", err)
		return
	}

	if mediaType != "video/mp4" {
		respondWithError(w, http.StatusBadRequest, "Cannot use non mp4 file format", err)
		return
	}

	tmpFile, err := os.CreateTemp("/tmp", "tubely-upload-*.mp4")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't create temp file", err)
		return
	}
	defer tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	if _, err := io.Copy(tmpFile, file); err != nil {
		tmpFile.Close()
		respondWithError(w, http.StatusBadRequest, "Couldn't copy data to temp file", err)
		return
	}

	if _, err := tmpFile.Seek(0, io.SeekStart); err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't set seek to start of temp file", err)
		return
	}

	randFileKey := make([]byte, 32)
	if _, err := rand.Read(randFileKey); err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't generate random file key", err)
		return
	}

	aspectRation, err := getVideoAspectRatio(tmpFile.Name())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't get aspect ratio of temp file", err)
		return
	}

	randFileName := base64.RawURLEncoding.EncodeToString(randFileKey) + ".mp4"
	fullObjectPath := ""
	switch aspectRation {
	case "16:9":
		fullObjectPath = "landscape/" + randFileName
	case "9:16":
		fullObjectPath = "portrait/" + randFileName
	case "other":
		fullObjectPath = "other/" + randFileName
	}

	processedUploadFilePath, err := processVideoForFastStart(tmpFile.Name())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't add fast start to video file", err)
		return
	}

	uploadFile, err := os.Open(processedUploadFilePath)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't open processed video file", err)
		return
	}
	defer uploadFile.Close()
	defer os.Remove(processedUploadFilePath)

	if _, err := cfg.s3Client.PutObject(
		context.TODO(),
		&s3.PutObjectInput{
			Bucket:      &cfg.s3Bucket,
			Key:         &fullObjectPath,
			Body:        uploadFile,
			ContentType: &mediaType,
		},
	); err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't load video file to aws", err)
		return
	}

	videoUrlFilePath := fmt.Sprintf(
		"https://storage.yandexcloud.net/%s/%s",
		cfg.s3Bucket,
		fullObjectPath,
	)
	videoRow.VideoURL = &videoUrlFilePath

	if err := cfg.db.UpdateVideo(videoRow); err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't update video in database", err)
		return
	}

	respondWithJSON(w, http.StatusOK, videoRow)
}
