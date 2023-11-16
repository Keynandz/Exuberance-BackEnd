package handlers

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	register "exuberance-backend/app/register/repositories"
)

func DefaultPicture(id uint) error {
	loadErr := godotenv.Load()
	if loadErr != nil {
		log.Fatal("error loading file .env")
	}

	ssl, _ := strconv.ParseBool(os.Getenv("MINIO_SSL"))
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKeyID := os.Getenv("MINIO_ACCESS_KEY")
	secretAccessKey := os.Getenv("MINIO_SECRET_KEY")
	useSSL := ssl

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	userId := int(id)
	bucketName := os.Getenv("MINIO_BUCKET")
	fileName := "users/default.jpeg"

	objectInfo, err := minioClient.StatObject(context.Background(), bucketName, fileName, minio.StatObjectOptions{})
	if err != nil {
		return err
	}

	defaultImageData := objectInfo.Key

	if err := register.DefaultPicture(userId, defaultImageData); err != nil {
		return err
	}

	return nil
}
