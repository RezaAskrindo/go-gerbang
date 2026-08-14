package services

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"go-gerbang/config"
	"go-gerbang/handlers"
	"go-gerbang/types"

	"github.com/gofiber/fiber/v3"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func HandleConfigFile(c fiber.Ctx) error {
	u := new(types.ConfigServices)

	if err := handlers.ParseBody(c, u); err != nil {
		return handlers.BadRequestErrorResponse(c, err)
	}

	if err := handlers.SaveConfig(config.BasePath+config.ConfigPath, u); err != nil {
		return handlers.InternalServerErrorResponse(c, err)
	}

	return handlers.SuccessResponse(c, true, "success update config file", u, nil)
}

type S3Config struct {
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	Region    string `json:"region"`
	Endpoint  string `json:"endpoint"`
	Bucket    string `json:"bucket"`
}

func decodeConfig(encoded string) (*S3Config, error) {
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}

	var cfg S3Config

	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func createSession(accessKey string, secretAccessKey string, endpoint string) (*minio.Client, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretAccessKey, ""),
		Secure: true,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %w", err)
	}

	return client, nil
}

func handleS3Upload(
	accessKey string,
	secretAccessKey string,
	endpoint string,
	bucketName string,
	region string,
	rawLocation string,
	files []*multipart.FileHeader,
) error {
	minioClient, err := createSession(accessKey, secretAccessKey, endpoint)
	if err != nil {
		return err
	}

	ctx := context.Background()

	exists, err := minioClient.BucketExists(ctx, bucketName)
	if err != nil {
		return err
	}

	if !exists {
		if err := minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{
			Region: region,
		}); err != nil {
			return err
		}
	}

	rawLocation = filepath.ToSlash(rawLocation)
	rawLocation = strings.Trim(rawLocation, "/")

	err = createS3Backup(ctx, minioClient, bucketName, rawLocation)
	if err != nil {
		return err
	}

	err = pruneS3Backups(ctx, minioClient, bucketName, rawLocation, 10)
	if err != nil {
		return err
	}

	for _, fileHeader := range files {
		cleanPath := filepath.ToSlash(getRelativePath(fileHeader))
		cleanPath = strings.TrimPrefix(cleanPath, "/")

		parts := strings.SplitN(cleanPath, "/", 2)

		var objectName string
		if len(parts) == 2 {
			objectName = parts[1]
		} else {
			objectName = parts[0]
		}

		if rawLocation != "" {
			objectName = rawLocation + "/" + objectName
		}

		file, err := fileHeader.Open()
		if err != nil {
			return fmt.Errorf("open %s: %w", objectName, err)
		}

		contentType := fileHeader.Header.Get("Content-Type")
		if contentType == "" {
			contentType = "application/octet-stream"
		}

		_, err = minioClient.PutObject(
			ctx,
			bucketName,
			objectName,
			file,
			fileHeader.Size,
			minio.PutObjectOptions{
				ContentType: contentType,
				UserMetadata: map[string]string{
					"x-amz-acl": "public-read",
				},
			},
		)

		_ = file.Close()

		if err != nil {
			return fmt.Errorf("upload %s: %w", objectName, err)
		}
	}

	return nil
}

func createS3Backup(
	ctx context.Context,
	client *minio.Client,
	bucketName string,
	rawLocation string,
) error {
	rawLocation = strings.Trim(rawLocation, "/")

	backupName := fmt.Sprintf("backup-%s", time.Now().Format("20060102_150405"))

	backupPrefix := backupName + "/"
	if rawLocation != "" {
		backupPrefix = rawLocation + "/" + backupName + "/"
	}

	sourcePrefix := rawLocation
	if sourcePrefix != "" {
		sourcePrefix += "/"
	}

	objects := client.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Prefix:    sourcePrefix,
		Recursive: true,
	})

	for obj := range objects {
		if obj.Err != nil {
			return obj.Err
		}

		relativeKey := strings.TrimPrefix(obj.Key, sourcePrefix)
		// Skip existing backups
		if strings.HasPrefix(relativeKey, "backup-") {
			continue
		}

		destination := backupPrefix + relativeKey

		src := minio.CopySrcOptions{Bucket: bucketName, Object: obj.Key}
		dst := minio.CopyDestOptions{Bucket: bucketName, Object: destination}
		_, err := client.CopyObject(ctx, dst, src)
		if err != nil {
			return fmt.Errorf("backup %s failed: %w", obj.Key, err)
		}

		if err := client.RemoveObject(ctx, bucketName, obj.Key, minio.RemoveObjectOptions{}); err != nil {
			return fmt.Errorf("remove original %s after backup failed: %w", obj.Key, err)
		}
	}

	return nil
}

func unique(items []string) []string {
	seen := make(map[string]struct{})
	result := make([]string, 0, len(items))

	for _, item := range items {
		if _, exists := seen[item]; exists {
			continue
		}

		seen[item] = struct{}{}
		result = append(result, item)
	}

	return result
}

func pruneS3Backups(
	ctx context.Context,
	client *minio.Client,
	bucket string,
	rawLocation string,
	keep int,
) error {

	rawLocation = strings.Trim(rawLocation, "/")

	prefix := ""
	if rawLocation != "" {
		prefix = rawLocation + "/"
	}

	objects := client.ListObjects(
		ctx,
		bucket,
		minio.ListObjectsOptions{
			Prefix:    prefix + "backup-",
			Recursive: false,
		},
	)

	var backups []string

	for obj := range objects {

		if obj.Err != nil {
			return obj.Err
		}

		parts := strings.Split(obj.Key, "/")

		if len(parts) == 0 {
			continue
		}

		backupName := parts[len(parts)-2]

		backups = append(
			backups,
			backupName,
		)
	}

	backups = unique(backups)

	sort.Sort(
		sort.Reverse(
			sort.StringSlice(backups),
		),
	)

	if len(backups) <= keep {
		return nil
	}

	for _, old := range backups[keep:] {

		deletePrefix := ""

		if rawLocation != "" {
			deletePrefix = rawLocation + "/"
		}

		deletePrefix += old + "/"

		objects := client.ListObjects(
			ctx,
			bucket,
			minio.ListObjectsOptions{
				Prefix:    deletePrefix,
				Recursive: true,
			},
		)

		for obj := range objects {

			if obj.Err != nil {
				return obj.Err
			}

			err := client.RemoveObject(
				ctx,
				bucket,
				obj.Key,
				minio.RemoveObjectOptions{},
			)

			if err != nil {
				return err
			}
		}
	}

	return nil
}

// NOTE: HIGH SECURITY
func HandleFileUpload(c fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}

	files := form.File["files"]
	rawLocation := form.Value["file-location"][0]

	cfgValue := form.Value["s3_config"]
	if len(cfgValue) != 0 {
		cfg, err := decodeConfig(cfgValue[0])
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid config")
		}

		err = handleS3Upload(cfg.AccessKey, cfg.SecretKey, cfg.Endpoint, cfg.Bucket, cfg.Region, rawLocation, files)
		if err != nil {
			return handlers.InternalServerErrorResponse(c, err)
		}
		return handlers.SuccessResponse(c, true, "success upload file to s3 server", nil, nil)
	}

	cleanPath := filepath.ToSlash(rawLocation)
	uploadRoot := cleanPath
	_ = os.MkdirAll(uploadRoot, 0755)

	entries, err := os.ReadDir(uploadRoot)
	if err != nil {
		return err
	}

	timestamp := time.Now().Format("20060102-150405")
	backupDir := filepath.Join(uploadRoot, "backup-"+timestamp)
	_ = os.MkdirAll(backupDir, 0755)

	for _, entry := range entries {
		name := entry.Name()
		if strings.HasPrefix(name, "backup-") {
			continue
		}

		srcPath := filepath.Join(uploadRoot, name)
		dstPath := filepath.Join(backupDir, name)

		if err := os.Rename(srcPath, dstPath); err != nil {
			return err
		}
	}

	if err := pruneOldBackups(uploadRoot, 10); err != nil {
		log.Printf("warning: failed to prune old backups: %v", err)
	}

	for _, file := range files {
		cleanPath := filepath.ToSlash(getRelativePath(file))
		cleanPath = strings.TrimPrefix(cleanPath, "/")

		parts := strings.SplitN(cleanPath, "/", 2)
		var relativePath string
		if len(parts) == 2 {
			relativePath = parts[1]
		} else {
			relativePath = parts[0]
		}

		targetPath := filepath.Join(uploadRoot, relativePath)
		// fmt.Println("targetPath:", targetPath)

		_ = os.MkdirAll(filepath.Dir(targetPath), 0755)

		if err := saveFileNoExec(file, targetPath); err != nil {
			return err
		}
	}

	return handlers.SuccessResponse(c, true, "success upload file", nil, nil)
}

func getRelativePath(fh *multipart.FileHeader) string {
	cd := fh.Header.Get("Content-Disposition")

	re := regexp.MustCompile(`filename="([^"]+)"`)
	match := re.FindStringSubmatch(cd)
	if len(match) > 1 {
		path := filepath.ToSlash(match[1])
		path = strings.TrimPrefix(path, "/")
		return path
	}

	return fh.Filename
}

func saveFileNoExec(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	ext := filepath.Ext(dst) // "" if no extension, e.g. ".exe", ".txt"

	var perm os.FileMode
	if ext == "" || strings.EqualFold(ext, ".exe") || strings.EqualFold(ext, ".sh") {
		perm = 0755 // rwxr-xr-x — executable
	} else {
		perm = 0644 // rw-r--r-- — not executable
	}

	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}

func pruneOldBackups(uploadRoot string, keep int) error {
	entries, err := os.ReadDir(uploadRoot)
	if err != nil {
		return err
	}

	var backups []os.DirEntry
	for _, e := range entries {
		if e.IsDir() && strings.HasPrefix(e.Name(), "backup-") {
			backups = append(backups, e)
		}
	}

	// "backup-YYYYMMDD-HHMMSS" sorts lexically == chronologically
	sort.Slice(backups, func(i, j int) bool {
		return backups[i].Name() < backups[j].Name()
	})

	if len(backups) <= keep {
		return nil
	}

	toDelete := backups[:len(backups)-keep] // oldest first
	for _, b := range toDelete {
		path := filepath.Join(uploadRoot, b.Name())
		if err := os.RemoveAll(path); err != nil {
			return fmt.Errorf("failed to remove old backup %s: %w", b.Name(), err)
		}
	}
	return nil
}
