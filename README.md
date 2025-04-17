# Building and Running the Go Project

## Prerequisites

- **Go installed:** Ensure you have Go installed on your system. You can download it from <https://go.dev/dl/>.
- This project use go version 1.23.2, so make you update or install the same or version above
- **Terminal:** You'll need a terminal or command prompt to execute the commands.

## Building the Project

1. **Open your terminal:** Open the terminal application on your operating system.

2. **Navigate to the project directory:** Use the `cd` command to go to the directory containing the Go code (the directory with the `main.go` file).

   ```
   cd project/path
   ```

   Replace `project/path` with the actual path to the project.

3. **Building for macOS:**

   - Set the target operating system:

     ```
     export GOOS=darwin
     ```

   - Set the target architecture (if needed, defaults to your machine's):

     ```
     export GOARCH=amd64 # For Intel Macs
     export GOARCH=arm64 # For Apple Silicon Macs
     ```

   - Build the executable:

     ```
     go build -o myapp
     ```

     The executable will be named `myapp`.

4. **Building for Windows:**

   - Set the target operating system:

     ```
     export GOOS=windows
     ```

   - Set the target architecture:

     ```
     export GOARCH=amd64 # For 64-bit Windows
     export GOARCH=386 # For 32-bit Windows
     ```

   - Build the executable:

     ```
     go build -o myapp.exe
     ```

     The executable will be named `myapp.exe`.

## Running the Project

After successfully building the project, you can run the executable:

1. **Open a terminal:** Open a terminal or command prompt.

2. **Navigate to the directory:** Use the `cd` command to go to the directory where the executable file is located (the same directory where you ran the `go build` command).

3. **Run the executable:**

   - On macOS, type `./` followed by the executable name:

     ```
     ./myapp
     ```

   - On Windows, type the executable name with the `.exe` extension:

     ```
     myapp.exe
     ```

Replace `myapp` with the actual name of your application.
