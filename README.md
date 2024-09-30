### Project for Counting Unique IP Addresses in a Large File

#### Project Description
This project implements a highly efficient algorithm for counting unique IP addresses in very large files using the HyperLogLog data structure. The program is written in Go and optimized for concurrent data processing, enabling fast processing of large volumes of information. Before running the main code (`main.go`), you need to create a test file with IP addresses using the `generator.py` script.

#### Project Structure
1. **`main.go`** — the main Go file that implements the logic for processing and counting unique IP addresses.
2. **`generator.py`** — an auxiliary Python script that generates a file with IP addresses for testing.
3. **`execution.log`** — a log file automatically created when running the `main.go` program (this file is in `.gitignore` and not tracked in the repository).

#### Preparation for Launch
Before running the main program, you need to create a test file with IP addresses:

1. Open the `generator.py` file and, if necessary, modify the desired file size:

   ```python
   TARGET_SIZE = 100 * 1024 * 1024 * 1024  # Change this value to the desired size (in bytes)
   ```

2. Run the generator to create the test file:

   ```bash
   python3 generator.py
   ```

   After running this command, a file named `ip_addresses` with the specified size will appear in the current directory.

#### Running the Main Program
Once the file with IP addresses is created, you can run the Go program to count the unique IP addresses:

```bash
go run cmd/ip-counter/main.go
```

#### Generating Profiling Reports
After the profiles are created, you can visualize them using the `pprof` tool. To do this, run the following commands:

   ```bash
   go tool pprof -http=:8080 cpu_profile.prof
   ```

   This report will help visually analyze performance bottlenecks and optimize CPU usage.

If you have any questions or suggestions for improving the code, feel free to make changes or create Issues on GitHub.
