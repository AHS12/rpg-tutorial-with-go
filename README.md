# RPG Tutorial Project

This is an RPG tutorial project built using Go and Ebiten v2, a popular 2D game library for Go.

## Features

- **2D Game Engine:** Built with Ebiten v2 for smooth rendering and game logic.
- **Cross-Platform Support:** Runs on Windows, macOS, Linux, and WebAssembly.
- **Simple Game Mechanics:** Perfect for learning game development with Go.
- **Modular Codebase:** Designed with maintainability and extensibility in mind.

## Prerequisites

Before you can run this project, make sure you have the following installed:

### Required Software

- **Go:** You need a working Go environment. Download it from [https://go.dev/dl/](https://go.dev/dl/). It's recommended to use Go version **1.24.1** or later, as specified in the `go.mod` file.
- **Git:** Git is used for version control and dependency management. Download it from [https://git-scm.com/downloads](https://git-scm.com/downloads).

### Optional (But Recommended)

- **Visual Studio Code (VS Code):** A lightweight and powerful IDE for Go development.
- **Ebiten Examples & Documentation:** Visit the official site at [https://ebitengine.org/](https://ebitengine.org/) for guides and examples.

## Setup

Follow these steps to get started with the project:

1. **Clone the Repository:**

    ```bash
    git clone https://github.com/AHS12/rpg-tutorial-with-go.git
    cd rpg-tutorial-with-go
    ```

2. **Download Dependencies:**

    Go automatically manages dependencies based on the `go.mod` file. Run the following command to download and install all required modules:

    ```bash
    go mod download
    ```
    Dependency Download Errors: If go mod download fails, check your internet connection and make sure you have sufficient permissions to write to the project directory. You can try running go mod tidy as an alternative.

    ```bash
    go mod tidy
    ```
    

## Running the Project

To run the project, use the following command in your terminal from the root directory (where `go.mod` is located):

```bash
go run .
```

If everything is set up correctly, the game window should open and display the current implementation of the RPG.

## Troubleshooting

If you encounter any issues, consider the following:

- Make sure you have installed the correct version of Go.
- Run `go mod tidy` to clean up unused dependencies.
- Check the [Ebiten documentation](https://ebitengine.org/) for additional setup help.
- If an error occurs when running the game, check the terminal output for hints.

## Contributing

Contributions are welcome! If you find a bug or want to improve the game, feel free to fork the repository and submit a pull request.

## License

This project is licensed under the MIT License. See `LICENSE` for more details.

---

Enjoy developing your RPG with Go and Ebiten! 🚀

