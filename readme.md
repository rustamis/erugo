# erugo

[![Build](https://github.com/DeanWard/erugo/actions/workflows/build.yml/badge.svg)](https://github.com/DeanWard/erugo/actions/workflows/build.yml)

erugo is a secure, self-hosted alternative to WeTransfer, built with Go and Vue.js. It combines powerful file-sharing capabilities with a sleek user interface, all packaged in a single, easy-to-deploy binary.

## Screenshots

![Upload Interface](.github/images/upload.png)
*A clean, intuitive upload interface showing file selection and progress*

![Share Details](.github/images/share.png)
*Share details view with file listing and expiration information*

![Share URL](.github/images/share-url.png)
*Simple one-click share URL copying*


# Key Features

- **Effortless Deployment**: Single binary contains both backend and frontend
- **Human-Friendly Share Links**: Easy-to-read URLs like `https://yourdomains.com/shares/patient-haze-tiny-term`
- **Secure Access Control**: Only authorized users can create shares, while anyone with a share link can download
- **Simple Data Management**: SQLite database for efficient metadata storage
- **Flexible Configuration**: Customizable storage paths and share size limits
- **Interactive Setup**: User-friendly first-run configuration using Bubble Tea
- **Modern Interface**: Clean, intuitive web UI
- **Open Source**: MIT licensed and ready for white-labeling

## Quick Start

1. Download the appropriate binary for your platform from the [Releases](https://github.com/DeanWard/erugo/releases/) page

2. For Mac/Linux systems, make the binary executable:
   ```sh
   chmod +x erugo-darwin-arm64
   ```

3. Launch erugo:
   ```sh
   ./erugo
   ```

4. Complete the interactive first-run setup to create your admin account

5. Access the web interface at `http://localhost:9199`

## Configuration Options

erugo can be customized through a configuration file with the following options:

| Option              | Description                               | Default Value           |
|--------------------|-------------------------------------------|------------------------|
| `app_url`          | Application hosting URL                   | `http://localhost:9199` |
| `base_storage_path`| File storage location                     | `storage`              |
| `max_share_size`   | Maximum file size per share              | `2G`                   |
| `bind_port`        | Web server port                          | `9199`                 |
| `jwt_secret`       | JWT authentication secret                 | `change_me`            |

A default `config.json` file is automatically generated on first run.

## Using erugo

### Creating a Share
1. Log in to the web interface
2. Select files for upload
3. Share the generated link with your recipient

### Downloading Files
Recipients simply need to:
1. Click the share link
2. Download the files through the web interface

## Customization

As an open-source project, erugo can be tailored to your needs:
- Customize the UI to match your brand
- Modify URL structures and authentication methods
- Extend functionality through code modifications

## Roadmap

We're actively developing erugo with the following features planned:

- **UI-Based White-Labeling**: Brand customization through the web interface
- **Enhanced File Access**: Optional direct file downloads without ZIP packaging
- **Flexible Database Configuration**: Configurable database file location
- **Docker Support**: Containerized deployment option

## Contributing

We welcome community contributions! Feel free to:
- Submit bug reports and feature requests
- Create pull requests
- Engage in discussions

## License

erugo is released under the MIT License, ensuring maximum flexibility for both personal and commercial use.

---

🚀 **Ready to start? Download erugo and begin sharing files securely in minutes!**