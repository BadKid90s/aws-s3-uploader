# Typora 集成指南

本文档介绍了如何将 aws-s3-uploader 与 Typora 集成，实现一键上传图片到 S3 存储服务。

## 配置步骤

### 1. 构建或下载 aws-s3-uploader

首先确保你已经构建了 aws-s3-uploader 工具，或者从 GitHub Release 页面下载了适用于你操作系统的预编译版本。

### 2. 配置 aws-s3-uploader

创建配置文件`config.toml`，填入你的 S3 凭证信息：

```toml
[default]
endpoint_url = "https://s3.amazonaws.com"
region = "us-east-1"
access_key_id = "your-access-key-id"
secret_access_key = "your-secret-access-key"
bucket = "your-bucket-name"
img_url_prefix = "https://your-bucket-name.s3.amazonaws.com"
directory = "images"  # 可选，指定上传目录
```

### 3. 配置 Typora

1. 打开 Typora
2. 进入 `偏好设置` → `图像`
3. 在 `上传服务` 部分，选择 `自定义命令`
4. 在命令输入框中输入以下命令：

```bash
# macOS/Linux
/path/to/aws-s3-uploader $imageFilePath

# Windows
C:\path\to\aws-s3-uploader.exe $imageFilePath
```

请将 `/path/to/aws-s3-uploader` 替换为你实际的 aws-s3-uploader 可执行文件路径。

### 4. 测试配置

1. 在 Typora 中插入一张图片
2. 右键点击图片，选择 `上传图片` 或使用快捷键 `Ctrl+U` (Windows/Linux) 或 `Cmd+U` (macOS)
3. 如果配置正确，图片将被上传到 S3，并且图片链接会自动更新为 S3 URL

## 高级配置

### 使用不同的上传目录

如果需要为不同项目使用不同的上传目录，可以通过命令行参数指定：

```bash
/path/to/aws-s3-uploader --directory=project1/images $imageFilePath
```

### 使用配置文件

如果需要为不同项目使用不同的配置，可以通过 `--config` 参数指定：

```bash
/path/to/aws-s3-uploader --config=/path/to/project1-config.toml $imageFilePath
```

## 故障排除

### 1. 权限问题

确保 aws-s3-uploader 可执行文件具有执行权限：

```bash
chmod +x /path/to/aws-s3-uploader
```

### 2. 配置文件未找到

确保 [config.toml](file:///Users/wry/IdeaProjects/typroa-s3-uploader/config.toml) 文件与 aws-s3-uploader 可执行文件位于同一目录，或者使用 `--config` 参数指定配置文件路径。

### 3. 上传失败

检查以下几点：
- 网络连接是否正常
- S3 凭证是否正确
- 是否有权限上传到指定的 bucket 和目录
- bucket 名称是否正确

### 4. URL 格式不正确

检查 [config.toml](file:///Users/wry/IdeaProjects/typroa-s3-uploader/config.toml) 中的 `img_url_prefix` 配置是否正确。

## 常见问题

### Q: 如何在 Windows 上使用？

A: 在 Windows 上使用时，确保使用 `.exe` 扩展名，并使用反斜杠或双反斜杠作为路径分隔符。

### Q: 如何自定义图片名称？

A: aws-s3-uploader 会自动将图片重命名为 Unix 时间戳格式，这是为了防止文件名冲突。

### Q: 支持哪些云存储服务？

A: aws-s3-uploader 支持任何兼容 S3 API 的存储服务，包括 AWS S3、Cloudflare R2、MinIO 等。