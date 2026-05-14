package services

import (
	"io"
	"io/fs"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"go-gerbang/config"
	"go-gerbang/handlers"
	"go-gerbang/types"

	"github.com/gofiber/fiber/v3"
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

func HandleFileUpload(c fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}

	rawLocation := form.Value["file-location"][0]
	cleanPath := filepath.ToSlash(rawLocation)
	uploadRoot := cleanPath
	os.MkdirAll(uploadRoot, 0755)

	entries, err := os.ReadDir(uploadRoot)
	if err != nil {
		return err
	}

	timestamp := time.Now().Format("20060102-150405")
	backupDir := filepath.Join(uploadRoot, "backup-"+timestamp)
	os.MkdirAll(backupDir, 0755)

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

	files := form.File["files"]

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

		os.MkdirAll(filepath.Dir(targetPath), 0755)

		if err := c.SaveFile(file, targetPath); err != nil {
			return err
		}
	}

	return handlers.SuccessResponse(c, true, "success upload file", nil, nil)
}

func copyAll(src, dst string) error {
	return filepath.Walk(src, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		target := filepath.Join(dst, rel)

		if info.IsDir() {
			if strings.HasPrefix(info.Name(), "backup-") {
				return filepath.SkipDir
			}
			return os.MkdirAll(target, 0755)
		}

		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			return err
		}

		in, err := os.Open(path)
		if err != nil {
			return err
		}
		defer in.Close()

		out, err := os.Create(target)
		if err != nil {
			return err
		}
		defer out.Close()

		_, err = io.Copy(out, in)
		return err
	})
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
