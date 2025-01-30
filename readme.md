# erugo

[![Build](https://github.com/DeanWard/erugo/actions/workflows/build.yml/badge.svg)](https://github.com/DeanWard/erugo/actions/workflows/build.yml)

erugo is a **self-hosted WeTransfer alternative** written in **Go**, with a **Vue.js** frontend. It provides a seamless file-sharing experience that is both secure and user-friendly. Designed to be easy to deploy and use, erugo runs from a **single binary** containing both the backend and frontend.

## Features
- **Single Binary Deployment** - No complex setup required; just download and run.
- **Secure & User-Friendly URLs** - Generates human-readable share links like `https://yourdomains.com/shares/patient-haze-tiny-term`.
- **Authorization for Creating Shares** - Only authorized users can create shares, but anyone with a share URL can download.
- **SQLite for Storage** - Uses SQLite to store share metadata for simplicity and efficiency.
- **Configurable Storage & Limits** - Customize storage path and max share size to suit your needs.
- **Interactive First-Run Setup** - Guides you through creating the first user account using **Bubble Tea**.
- **Modern UI** - Provides a sleek and user-friendly web interface.
- **Fully Open-Source & White-Label Ready** - Customize the branding and modify the code as needed.

## Installation
1. **Download the binary** for your platform from the [Releases](#) page.
2. **Run the binary:**
   ```sh
   ./erugo
   ```
3. **Follow the first-run interactive setup** to create your admin account.
4. **Access the web UI** at `http://localhost:9199` (default port).
5. **Start sharing files!**

## Important Note for macOS Users
**Warning:** macOS users might encounter issues with Gatekeeper when running the binary. For more information and steps to resolve this, please refer to the [Gatekeeper Issues](#gatekeeper-issues) section below.


## Configuration
erugo allows you to customize various settings via a config file or environment variables.

### Available Options:
| Option              | Description                                          | Default Value |
|---------------------|------------------------------------------------------|-------------------------|
| `app_url`           | URL where the app is hosted                          | `http://localhost:9199` |
| `base_storage_path` | Path where uploaded files are stored                 | `storage`               |
| `max_share_size`    | Maximum file size per share in human readable format | `2G`                    |
| `bind_port`         | Port for the web server                              | `9199`                  |
| `jwt_secret`        | Secret for JWT authentication                        | `change_me`             |

To run with custom options, you can edit the `config.json` file or set environment variables. This file will be created when you run the app for the first time.

## Usage
### Uploading a File (Authorized Users Only)
1. **Log in to the web UI**.
2. **Upload your file** and generate a shareable link.
3. **Send the link** to the recipient.

### Downloading a File
- Simply visit the provided share URL in a browser

## White-Labeling
Since erugo is fully open-source, you can:
- Modify the UI to match your branding.
- Customize URLs and authentication.
- Extend functionality as needed.

## Contributing
We welcome contributions! Feel free to submit issues, feature requests, or pull requests.

## License
erugo is licensed under the **MIT License**, allowing free and open-source usage.

---

### Get Started Now
🚀 **Download erugo, set it up in minutes, and start sharing files securely!**

