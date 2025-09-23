package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var (
	endpointURL     = flag.String("endpoint-url", "", "R2 endpoint URL")
	region          = flag.String("region", "", "R2 region")
	accessKey       = flag.String("access-key", "", "R2 access key")
	secretAccessKey = flag.String("secret-access-key", "", "R2 secret access key")
	bucket          = flag.String("bucket", "", "R2 bucket name")
	imgURLPrefix    = flag.String("img-url-prefix", "", "Image URL prefix")
	directory       = flag.String("directory", "", "Directory path in S3 (optional)")
	configFile      = flag.String("config", "", "Configuration file path")
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
	} else {
		// Validate required flags
		if *endpointURL == "" || *region == "" || *accessKey == "" ||
			*secretAccessKey == "" || *bucket == "" {
			fmt.Fprintf(os.Stderr, "Error: All R2 credentials must be provided either via command line or config file\n")
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

	// Generate timestamped file name ONCE
	filename := filepath.Base(filePath)
	timestamp := time.Now().Unix()
	timestampedFileName := generateTimestampedFilename(filename, timestamp)

	// Upload the file
	err = uploadFile(&cfg.Default, filePath, timestampedFileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error uploading file: %v\n", err)
		os.Exit(1)
	}

	// Print the URL of the uploaded file
	printFileURL(&cfg.Default, filePath, timestampedFileName)
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

	// Create a new AWS session
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(cfg.Region),
		Endpoint:    aws.String(cfg.EndpointURL),
		Credentials: credentials.NewStaticCredentials(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
	})
	if err != nil {
		return err
	}

	// Create S3 service client
	svc := s3.New(sess)

	// Upload the file
	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(cfg.Bucket),
		Key:    aws.String(s3Key),
		Body:   file,
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
