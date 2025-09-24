package main

import (
	"testing"
)

func TestGenerateTimestampedFilename(t *testing.T) {
	// Test with a PNG file
	filename := "logo.png"
	timestamp := int64(1758632000)
	expected := "1758632000.png"
	result := generateTimestampedFilename(filename, timestamp)

	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}

	// Test with a JPG file
	filename = "image.jpg"
	expected = "1758632000.jpg"
	result = generateTimestampedFilename(filename, timestamp)

	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}

	// Test with a file without extension
	filename = "README"
	expected = "1758632000"
	result = generateTimestampedFilename(filename, timestamp)

	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestGenerateS3Key(t *testing.T) {
	timestampedFilename := "1758632000.png"

	// Test with empty directory
	directory := ""
	expected := "1758632000.png"
	result := generateS3Key(timestampedFilename, directory)

	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}

	// Test with a directory
	directory = "uploads"
	expected = "uploads/1758632000.png"
	result = generateS3Key(timestampedFilename, directory)

	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}

	// Test with a nested directory
	directory = "images/logos"
	expected = "images/logos/1758632000.png"
	result = generateS3Key(timestampedFilename, directory)

	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestGenerateURLPath(t *testing.T) {
	timestampedFilename := "1758632000.png"

	// Test with empty directory
	directory := ""
	expected := "1758632000.png"
	result := generateURLPath(timestampedFilename, directory)

	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}

	// Test with a directory
	directory = "uploads"
	expected = "uploads/1758632000.png"
	result = generateURLPath(timestampedFilename, directory)

	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}

	// Test with a nested directory
	directory = "images/logos"
	expected = "images/logos/1758632000.png"
	result = generateURLPath(timestampedFilename, directory)

	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

// TestRenameFeature 测试重命名功能
func TestRenameFeature(t *testing.T) {
	// 测试启用重命名功能时的时间戳文件名生成
	filename := "test.txt"
	timestamp := int64(1758686243)

	// 测试 generateTimestampedFilename 函数
	expected := "1758686243.txt"
	result := generateTimestampedFilename(filename, timestamp)

	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}

	// 测试不同扩展名
	filename = "image.png"
	expected = "1758686243.png"
	result = generateTimestampedFilename(filename, timestamp)

	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}

	// 测试没有扩展名的文件
	filename = "README"
	expected = "1758686243"
	result = generateTimestampedFilename(filename, timestamp)

	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestIntegrationWithFiles(t *testing.T) {
	// Test with logo.png
	filename := "logo.png"
	timestamp := int64(1758632000)
	expected := "1758632000.png"
	result := generateTimestampedFilename(filename, timestamp)

	if result != expected {
		t.Errorf("For logo.png: expected %s, got %s", expected, result)
	}

	// Test with dir/logo.png (we only use the base filename)
	filename = "logo.png" // Base filename from dir/logo.png
	expected = "1758632000.png"
	result = generateTimestampedFilename(filename, timestamp)

	if result != expected {
		t.Errorf("For dir/logo.png: expected %s, got %s", expected, result)
	}

	// Test S3 key generation for logo.png with no directory
	timestampedFilename := "1758632000.png"
	directory := ""
	expected = "1758632000.png"
	result = generateS3Key(timestampedFilename, directory)

	if result != expected {
		t.Errorf("For logo.png with no directory: expected %s, got %s", expected, result)
	}

	// Test S3 key generation for logo.png with directory
	directory = "uploads"
	expected = "uploads/1758632000.png"
	result = generateS3Key(timestampedFilename, directory)

	if result != expected {
		t.Errorf("For logo.png with directory: expected %s, got %s", expected, result)
	}

	// Test S3 key generation for dir/logo.png with no directory
	// (should be the same as above since we only use the base filename)
	directory = ""
	expected = "1758632000.png"
	result = generateS3Key(timestampedFilename, directory)

	if result != expected {
		t.Errorf("For dir/logo.png with no directory: expected %s, got %s", expected, result)
	}

	// Test S3 key generation for dir/logo.png with directory
	directory = "images"
	expected = "images/1758632000.png"
	result = generateS3Key(timestampedFilename, directory)

	if result != expected {
		t.Errorf("For dir/logo.png with directory: expected %s, got %s", expected, result)
	}
}
