package uploadService

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2"
	"hunter-backend/util"
	"io"
)

func (u uploadService) HandlerUploadFile(c *fiber.Ctx) error {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		panic(errors.New("file is required"))
	}

	file, err := fileHeader.Open()
	if err != nil {
		panic(err)
	}
	defer file.Close()

	buf := bytes.NewBuffer(nil)
	size, err := io.Copy(buf, file)
	if err != nil {
		panic(err)
	}

	filename, err := util.GenerateRandomFilename(fileHeader.Filename)
	if err != nil {
		panic(err)
	}

	_, err = u.s3Client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket:        aws.String("uploader"),
		Key:           aws.String(filename),
		Body:          bytes.NewReader(buf.Bytes()),
		ContentLength: &size,
		ContentType:   aws.String(fileHeader.Header.Get("Content-Type")),
	})
	if err != nil {
		panic(err)
	}

	return c.JSON(fiber.Map{
		"url":  fmt.Sprintf("%s/%s", u.config.S3Config.S3PublicEndpoint, filename),
		"file_name": filename,
	})
}
