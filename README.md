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
go run main.go
```

The program uses concurrency and automatically adjusts the number of workers based on the file size to ensure maximum performance.

#### Enabling Profiling
If you want to enable CPU and memory profiling, open the `main.go` file and uncomment the corresponding lines.

After uncommenting and running `main.go`, two profiles will be created in the current directory:
- `cpu_profile.prof` — CPU usage profile.
- `heap_profile.prof` — memory usage profile.

#### Generating Profiling Reports
After the profiles are created, you can visualize them using the `pprof` tool. To do this, run the following commands:

1. **Create a PNG file for the CPU profile**:

   ```bash
   go tool pprof -png cpu_profile.prof > cpu_profile.png
   ```

   This command generates a call graph as an image file named `cpu_profile.png`.

2. **Create a text report for the CPU profile**:

   ```bash
   go tool pprof -text cpu_profile.prof > cpu_profile.txt
   ```

   The text report `cpu_profile.txt` contains detailed information on function execution times and resource usage.

3. **Create a PNG file for the memory (heap) profile**:

   ```bash
   go tool pprof -png heap_profile.prof > heap_profile.png
   ```

4. **Create a text report for the memory (heap) profile**:

   ```bash
   go tool pprof -text heap_profile.prof > heap_profile.txt
   ```

These reports will help visually analyze performance bottlenecks and optimize memory usage.

#### Optimization Notes
- The `main.go` code includes dynamic adaptation of the number of workers depending on the size of the input file.
- If the file is small (less than 2 GB), only one worker is used to avoid excessive resource usage.
- If necessary, you can modify the `baseChunkSize` to adjust the block size processed by each worker:

  ```go
  const baseChunkSize int64 = 2 * 1024 * 1024 * 1024  // 2 GB
  ```

#### Conclusion
The project allows flexible and efficient processing of very large files with IP addresses by automatically adapting to the size of the input data and using concurrency. Enabling profiling and analyzing the reports will help you gain a deeper understanding of the program's performance and improve the data processing algorithms.

### Commands for Running and Analysis
1. **Creating a file with IP addresses**:

   ```bash
   python3 generator.py
   ```

2. **Running the main program**:

   ```bash
   go run main.go
   ```

3. **Profiling analysis (after uncommenting in `main.go`)**:

   ```bash
   go tool pprof -png cpu_profile.prof > cpu_profile.png
   go tool pprof -text cpu_profile.prof > cpu_profile.txt
   go tool pprof -png heap_profile.prof > heap_profile.png
   go tool pprof -text heap_profile.prof > heap_profile.txt
   ```

If you have any questions or suggestions for improving the code, feel free to make changes or create Issues on GitHub.
