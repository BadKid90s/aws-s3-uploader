package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	endpointURL     = flag.String("endpoint-url", "https://s3.amazonaws.com", "S3 endpoint URL")
	region          = flag.String("region", "auto", "S3 region")
	accessKey       = flag.String("access-key", "", "S3 access key")
	secretAccessKey = flag.String("secret-access-key", "", "S3 secret access key")
	bucket          = flag.String("bucket", "", "S3 bucket name")
	imgURLPrefix    = flag.String("img-url-prefix", "", "Image URL prefix")
	directory       = flag.String("directory", "", "Directory path in S3 (optional)")
	configFile      = flag.String("config", "", "Configuration file path")
	renameFile      = flag.Bool("rename-file", false, "Whether to rename file with timestamp") // 新增命令行标志
)

// generateTimestampedFilename generates a timestamped filename based on the original filename
func generateTimestampedFilename(filename string, timestamp int64) string {
	fileExt := filepath.Ext(filename)
	return fmt.Sprintf("%d%s", timestamp, fileExt)
}

// generateS3Key generates the S3 key based on the timestamped filename and directory
func generateS3Key(timestampedFilename, directory string) string {
	if directory == "" {
		return timestampedFilename
	}
	return filepath.ToSlash(filepath.Join(directory, timestampedFilename))
}

// generateURLPath generates the URL path based on the timestamped filename and directory
func generateURLPath(timestampedFilename, directory string) string {
	if directory == "" {
		return timestampedFilename
	}
	return filepath.ToSlash(filepath.Join(directory, timestampedFilename))
}

func main() {
	flag.Parse()

	// Check if config file is provided
	var cfg *Config
	var err error

	// Keep track of which flags were explicitly set via command line
	explicitlySet := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) {
		explicitlySet[f.Name] = true
	})

	// Determine the directory where the executable is located
	execPath, err := os.Executable()
	if err != nil {
		// Fallback to current directory if we can't determine executable path
		execPath = "."
	}
	execDir := filepath.Dir(execPath)

	// If no config file is specified, try to load default config file from executable directory
	configFilePath := ""
	if *configFile != "" {
		configFilePath = *configFile
	} else {
		// Try to load default config file from executable directory
		defaultConfigPath := filepath.Join(execDir, "config.toml")
		if _, err := os.Stat(defaultConfigPath); err == nil {
			configFilePath = defaultConfigPath
		}
	}

	if configFilePath != "" {
		cfg, err = LoadConfigFromFile(configFilePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config file: %v\n", err)
			os.Exit(1)
		}

		// Override config values with command line arguments if provided
		if !explicitlySet["endpoint-url"] && cfg.Default.EndpointURL != "" {
			*endpointURL = cfg.Default.EndpointURL
		}
		if !explicitlySet["region"] && cfg.Default.Region != "" {
			*region = cfg.Default.Region
		}
		if !explicitlySet["access-key"] && cfg.Default.AccessKeyID != "" {
			*accessKey = cfg.Default.AccessKeyID
		}
		if !explicitlySet["secret-access-key"] && cfg.Default.SecretAccessKey != "" {
			*secretAccessKey = cfg.Default.SecretAccessKey
		}
		if !explicitlySet["bucket"] && cfg.Default.Bucket != "" {
			*bucket = cfg.Default.Bucket
		}
		if !explicitlySet["img-url-prefix"] && cfg.Default.ImgURLPrefix != "" {
			*imgURLPrefix = cfg.Default.ImgURLPrefix
		}
		// Use directory from config file only if not explicitly set via command line
		if !explicitlySet["directory"] && cfg.Default.Directory != "" {
			*directory = cfg.Default.Directory
		}
		// Use rename_file from config file only if not explicitly set via command line
		if !explicitlySet["rename-file"] && cfg.Default.RenameFile {
			*renameFile = cfg.Default.RenameFile
		}
	} else {
		// Validate required flags
		if *endpointURL == "" || *region == "" || *accessKey == "" ||
			*secretAccessKey == "" || *bucket == "" {
			fmt.Fprintf(os.Stderr, "Error: All S3 credentials must be provided either via command line or config file\n")
			flag.Usage()
			os.Exit(1)
		}

		cfg = &Config{
			Default: DefaultConfig{
				EndpointURL:     *endpointURL,
				Region:          *region,
				AccessKeyID:     *accessKey,
				SecretAccessKey: *secretAccessKey,
				Bucket:          *bucket,
				ImgURLPrefix:    *imgURLPrefix,
				Directory:       *directory,
				RenameFile:      *renameFile,
			},
		}
	}

	// Get file path from command line arguments
	args := flag.Args()
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "Error: File path is required\n")
		flag.Usage()
		os.Exit(1)
	}

	filePath := args[0]

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: File %s does not exist\n", filePath)
		os.Exit(1)
	}

	// Generate timestamped file name if renameFile is true
	filename := filepath.Base(filePath)
	timestampedFileName := filename // 默认使用原始文件名

	// 如果启用了重命名功能，则使用时间戳重命名
	if *renameFile {
		timestamp := time.Now().Unix()
		timestampedFileName = generateTimestampedFilename(filename, timestamp)
	}

	// Upload the file
	err = uploadFile(&cfg.Default, filePath, timestampedFileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error uploading file: %v\n", err)
		os.Exit(1)
	}

	// Print the URL of the uploaded file
	printFileURL(&cfg.Default, filePath, timestampedFileName)
}

// detectContentType detects the content type of a file based on its extension or content
func detectContentType(filePath string) string {
	// First, try to detect content type by file extension
	contentType := mime.TypeByExtension(filepath.Ext(filePath))
	if contentType != "" {
		return contentType
	}

	// If that fails, try to detect by opening the file and reading the first 512 bytes
	file, err := os.Open(filePath)
	if err != nil {
		return "application/octet-stream" // default binary type
	}
	defer file.Close()

	// Read the first 512 bytes (standard buffer size for DetectContentType)
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil && err != io.EOF {
		return "application/octet-stream"
	}

	// Use Go's built-in content type detection
	contentType = http.DetectContentType(buffer)
	if contentType != "" {
		return contentType
	}

	return "application/octet-stream" // default binary type
}

func uploadFile(cfg *DefaultConfig, filePath, timestampedFileName string) error {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Determine the S3 key based on the directory flag
	s3Key := generateS3Key(timestampedFileName, *directory)

	// Detect content type
	contentType := detectContentType(filePath)

	// Create a new AWS config
	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.SecretAccessKey, "")),
		config.WithRegion(cfg.Region),
	)
	if err != nil {
		return err
	}

	// Create S3 service client with custom endpoint
	svc := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(cfg.EndpointURL)
	})

	// Upload the file with content type
	_, err = svc.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(cfg.Bucket),
		Key:         aws.String(s3Key),
		Body:        file,
		ContentType: aws.String(contentType),
	})

	return err
}

func printFileURL(cfg *DefaultConfig, filePath, timestampedFileName string) {
	// Determine the URL path based on the directory flag
	urlPath := generateURLPath(timestampedFileName, *directory)

	// If img_url_prefix is provided, use it; otherwise construct from endpoint
	if cfg.ImgURLPrefix != "" {
		fmt.Printf("%s/%s\n", cfg.ImgURLPrefix, urlPath)
	} else {
		fmt.Printf("%s/%s/%s\n", cfg.EndpointURL, cfg.Bucket, urlPath)
	}
}
