# Cloudflare R2 Uploader

这是一个使用 Go 和 AWS SDK 上传图片到 Cloudflare R2 的命令行工具。上传的文件会被重命名为当前时间戳格式（Unix 时间戳）。

## 安装

确保你已经安装了 Go 1.16 或更高版本。

构建程序：
```bash
go build -o cloudflare-r2-uploader
```

或者使用 Makefile：
```bash
make build
```

## 配置文件

程序使用 TOML 格式的配置文件。你可以复制项目中的 `example.config.toml` 文件并将其重命名为 `config.toml` 来创建你的配置文件。

```bash
cp example.config.toml config.toml
```

然后编辑 `config.toml` 文件，填入你的 Cloudflare R2 凭证和配置信息：

```toml
[default]
endpoint_url = "https://your-account-id.r2.cloudflarestorage.com"
region = "auto"
access_key_id = "your-access-key-id"
secret_access_key = "your-secret-access-key"
bucket = "your-bucket-name"
img_url_prefix = "https://your-image-url-prefix"
directory = "uploads"  # 可选，指定文件上传到 S3 的目录路径
```

程序会在可执行文件所在目录自动查找 `config.toml` 文件。

## 运行测试

程序包含单元测试，可以使用以下命令运行：

```bash
go test -v
```

## 使用方法

### 通过命令行参数

```bash
./cloudflare-r2-uploader --endpoint-url <YOUR_ENDPOINT_URL> --region <YOUR_REGION> --access-key <YOUR_ACCESS_KEY> --secret-access-key <YOUR_SECRET_ACCESS_KEY> --bucket <YOUR_BUCKET_NAME> --img-url-prefix <YOUR_IMG_URL_PREFIX> [--directory <DIRECTORY_PATH>] <FILE_PATH>
```

### 通过配置文件

你也可以通过 `--config` 参数指定配置文件路径：

```bash
./cloudflare-r2-uploader --config=<CONFIG_FILE_PATH> <FILE_PATH>
```

## 参数说明

- `--endpoint-url`: R2 endpoint URL
- `--region`: R2 region (通常为 "auto")
- `--access-key`: R2 access key
- `--secret-access-key`: R2 secret access key
- `--bucket`: R2 bucket name
- `--img-url-prefix`: 图片 URL 前缀（可选）
- `--directory`: 上传到 S3 的目录路径（可选，默认为空，可通过配置文件设置）
- `--config`: 配置文件路径（可选，程序会自动在可执行文件目录下查找 config.toml）
- `<FILE_PATH>`: 要上传的文件路径

## 输出

- 成功时，程序会输出上传文件的 URL 到标准输出，格式为：`imgUrlPrefix + filePath`
- 失败时，错误信息会输出到标准错误流

## 示例

使用命令行参数：
```bash
./cloudflare-r2-uploader --endpoint-url <YOUR_ENDPOINT_URL> --region <YOUR_REGION> --access-key <YOUR_ACCESS_KEY> --secret-access-key <YOUR_SECRET_ACCESS_KEY> --bucket <YOUR_BUCKET_NAME> --img-url-prefix <YOUR_IMG_URL_PREFIX> [--directory <DIRECTORY_PATH>] <FILE_PATH>
```

使用配置文件：
```bash
# 使用配置文件中指定的目录
./cloudflare-r2-uploader <FILE_PATH>

# 上传到指定目录（通过命令行参数）
./cloudflare-r2-uploader --directory <DIRECTORY_PATH> <FILE_PATH>

# 上传到根目录（通过命令行参数覆盖配置文件）
./cloudflare-r2-uploader --directory="" <FILE_PATH>
```

输出结果：
```
# 上传到根目录
https://your-image-url-prefix/1758620387.png

# 上传到指定目录
https://your-image-url-prefix/path/to/file/1758620390.png
```

## 多平台支持

### 本地构建

使用 Makefile 可以交叉编译多个平台的二进制文件：
```bash
make cross-compile
```

这将生成以下平台的二进制文件：
- Darwin (macOS) AMD64
- Darwin (macOS) ARM64
- Linux AMD64
- Linux ARM64
- Windows AMD64

生成的二进制文件将位于 `dist/` 目录中。

### GitHub Actions 自动构建

该项目配置了 GitHub Actions，在创建 Release 时会自动构建并上传各平台的可执行文件。

要创建新的 Release：
1. 在本地创建并推送一个新的 tag（例如 `git tag v1.0.0 && git push origin v1.0.0`）
2. GitHub Actions 会自动触发，构建各平台的二进制文件并创建 Release
3. 在 GitHub 的 Release 页面可以下载各平台的可执行文件