帮我用go，使用aws s3 sdk 写一个图片上传程序`cloudflare-r2-uploader`

参数使用同级目录下的配置文件或者通过命令行指定。

通过`fmt.Println("msg")`打印结果：

- 使用错误输出流，当响应值不是200的时打印错误。

- 使用标准输出流，当响应值是200的时候打印image中的url属性。

命令行运行方式：

```
cloudflare-r2-uploader --endpoint-url <YOUR_ENDPOINT_URL> --region <YOUR_REGION> --access-key <YOUR_ACCESS_KEY> --secret-access-key <YOUR_SECRET_ACCESS_KEY> --bucket <YOUR_BUCKET_NAME> --img-url-prefix <YOUR_IMG_URL_PREFIX> <FILE_PATH>
```

配置文件方式：

```
cloudflare-r2-uploader --config=config
```

配置文件名称"config.toml", 内容：

```
[default]
endpoint_url = https://dacf4433278b05d1e393438ec680646a.r2.cloudflarestorage.com
region = auto
access_key_id = fe268961365d11144c26c65e79a28bec
secret_access_key = 5ec9eb046abf03fb4c434eca1891af51430e3f183f30b505de9575f8df054d44
bucket = chilix
img_url_prefix = https://r2.chilix.qzz.io
```

输出的内容是imgUrlPrefix+filePath

```
https://r2.chilix.qzz.io/chilix/blog.png
```

使用Makefile编译各个平台的二进制文件，实现多平台支持。

## Qwen Added Memories
- 开发了一个名为 cloudflare-r2-uploader 的 Go 程序，用于将图片上传到 Cloudflare R2 存储服务。程序支持通过命令行参数或 TOML 配置文件两种方式配置连接参数。输出格式为 imgUrlPrefix + filePath。程序支持多平台交叉编译，可通过 Makefile 生成 Darwin (macOS)、Linux 和 Windows 的二进制文件。
- 程序使用 AWS S3 SDK (github.com/aws/aws-sdk-go) 连接到 Cloudflare R2 服务。配置文件使用 TOML 格式解析 (github.com/BurntSushi/toml)。主要功能包括命令行参数解析、配置文件加载、文件上传到 R2 存储桶以及按指定格式输出文件 URL。
- 项目结构包含以下文件：main.go（主程序）、config.go（配置解析）、config.toml（配置示例）、Makefile（构建和交叉编译）、README.md（使用说明）、go.mod 和 go.sum（Go 模块依赖）。程序支持通过 make cross-compile 命令生成多平台二进制文件，分别位于 dist 目录下的各平台子目录中。
- Cloudflare R2 图片上传程序现在会将上传的文件重命名为当前时间戳格式（Unix 时间戳），例如 1758620387.png。这个功能确保了文件名的唯一性，避免了文件名冲突的问题。程序仍然支持 --directory 参数来保留目录结构，同时将文件名替换为时间戳格式。
- Cloudflare R2 图片上传程序现在支持在配置文件中设置 directory 参数。用户可以在 config.toml 文件中添加 directory = true/false 来控制是否保留目录结构。命令行参数 --directory 会覆盖配置文件中的设置。文件名仍然会重命名为当前时间戳格式（Unix 时间戳），例如 1758620387.png。
- Cloudflare R2 Uploader 程序已完成所有功能开发：1) 支持命令行参数和 TOML 配置文件两种配置方式；2) 支持 --directory 参数控制是否保留目录结构；3) 自动将上传文件重命名为 Unix 时间戳格式确保唯一性；4) 支持多平台交叉编译；5) 集成 GitHub Actions 自动发布 Release。
- Cloudflare R2 Uploader 程序已更新，现在支持在可执行文件目录下自动查找配置文件 config.toml。--directory 参数已修改为接受目录路径字符串，而不是布尔值。当命令行和配置文件同时指定参数时，命令行参数优先级更高。程序会优先使用命令行指定的目录，如果没有指定则检查配置文件中的 directory 设置，如果都没有则上传到根目录。
- Cloudflare R2 Uploader 程序已更新，配置文件中的 directory 参数现在用于指定文件上传到 S3 的目录路径，而不是布尔值。程序会在可执行文件目录下自动查找配置文件 config.toml。--directory 命令行参数会覆盖配置文件中的设置。当两者都未设置时，文件会上传到 S3 存储桶的根目录。
- Cloudflare R2 Uploader 程序已修复路径分隔符问题，现在无论在什么操作系统上运行，都会在 S3 key 和 URL 中使用正斜杠 `/` 而不是反斜杠 `\`。这是通过使用 filepath.ToSlash 函数实现的，确保了 URL 的正确性。
- Cloudflare R2 Uploader 程序已修复文件名时间戳不一致的问题。现在时间戳只在程序开始时生成一次，并在上传文件和打印 URL 时使用相同的值，确保两者完全一致。同时保留了路径分隔符的修复，确保在所有平台上都使用正斜杠。
- Cloudflare R2 Uploader 程序已完成单元测试，包含对时间戳文件名生成、S3 key 生成和 URL 路径生成的测试。测试覆盖了 logo.png 和 dir/logo.png 两种文件路径情况。程序通过了所有测试，并且实际功能测试也正常工作。添加了测试说明到 README.md 文件。
