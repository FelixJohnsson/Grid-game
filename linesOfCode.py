import os

def count_lines_of_code(directory=".", extensions=None, exclude_comments=True, ignored_dirs=None):
    if extensions is None:
        # Default file extensions to consider as code files
        extensions = ['.py', '.js', '.tsx', '.ts', '.go']

    if ignored_dirs is None:
        # Directories to ignore, e.g., 'node_modules', '.git', etc.
        ignored_dirs = ['node_modules']

    total_lines = 0
    total_files = 0

    for root, dirs, files in os.walk(directory):
        # Skip ignored directories
        dirs[:] = [d for d in dirs if d not in ignored_dirs]

        for file in files:
            if any(file.endswith(ext) for ext in extensions):
                file_path = os.path.join(root, file)
                with open(file_path, 'r', encoding='utf-8', errors='ignore') as f:
                    lines = f.readlines()
                    total_files += 1
                    for line in lines:
                        stripped_line = line.strip()
                        if exclude_comments:
                            # Skip empty lines and lines that start with comment characters
                            if stripped_line and not stripped_line.startswith(('#', '//', '/*', '*', '--')):
                                total_lines += 1
                        else:
                            # Count all non-empty lines
                            if stripped_line:
                                total_lines += 1

    return total_files, total_lines

if __name__ == "__main__":
    directory = "."
    total_files, total_lines = count_lines_of_code(directory)

    print(f"Total files scanned: {total_files}")
    print(f"Total lines of code: {total_lines}")


# 25/08/2024
# Total files scanned: 18
# Total lines of code: 1217

# 28/08/2024
#Total files scanned: 20
#Total lines of code: 2035

# 29/08/2024
#Total files scanned: 25
#Total lines of code: 2601

# 30/08/2024
#Total files scanned: 27
#Total lines of code: 3089

# 31/08/2024
#Total files scanned: 29
#Total lines of code: 3381

# 02/09/2024
#Total files scanned: 28
#Total lines of code: 3384