# AWS S3 Uploader v1.2.0 发布说明

我们很高兴地宣布 AWS S3 Uploader v1.2.0 正式发布！此版本带来了重要的新功能和改进，提升了用户体验和文件上传的灵活性。

## 🔧 新增功能

### 文件重命名选项
- 添加了新的配置项 `rename_file`（默认值：false）
- 添加了命令行参数 `--rename-file` 来控制是否将上传的文件重命名为时间戳格式
- 当设置为 true 时，文件将被重命名为 Unix 时间戳格式（例如：1758686243.png）
- 当设置为 false 时，文件将保持原始文件名

### Content-Type 自动检测
- 改进了文件上传功能，现在会自动检测并设置正确的 Content-Type
- 首先尝试通过文件扩展名检测 MIME 类型
- 如果扩展名检测失败，则读取文件内容进行类型推断
- 确保上传到 S3 的文件具有正确的 Content-Type，提升浏览器兼容性

## 📋 配置变更

### config.go
- 在 `DefaultConfig` 结构体中添加了 `RenameFile` 字段

### example.config.toml
- 添加了 `rename_file = false` 配置项示例

## 🚀 命令行参数

### 新增参数
- `--rename-file`: 控制是否重命名文件为时间戳格式

## 📊 完整更新日志

查看 [CHANGELOG.md](CHANGELOG.md) 了解完整版本更新信息。

## 📦 安装

您可以从 [Releases](https://github.com/BadKid90s/aws-s3-uploader/releases) 页面下载适用于您操作系统的预编译二进制文件，或者使用 Go 进行安装：

```bash
go install github.com/BadKid90s/aws-s3-uploader@v1.2.0
```

## 📖 使用说明

请参考 [README.md](README.md) 了解更多使用方法和配置选项。

## ❤️ 感谢

感谢所有为这个项目做出贡献的开发者和用户！