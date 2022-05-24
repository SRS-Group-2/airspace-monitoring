from __future__ import annotations

import os
import re
import subprocess
from typing import Sequence

tags_file: str = "terraform/image-versions.auto.tfvars" 

def main(argv: Sequence[str] | None = None) -> int:
    process = subprocess.Popen(['git', 'status'],
                     stdout=subprocess.PIPE, 
                     stderr=subprocess.PIPE)
    stdout, stderr = process.communicate()
    process.kill()

    if tags_file in stdout.decode('utf-8'):
        # file has been modified by user, not touching it
        return 0
    else:
        # file has not been modified by user, updating it
        docker_directories = [x[0] for x in os.walk('.') if 'Dockerfile' in x[2] and '.devcontainer' not in x[0]]
        original_file_lines = re.sub("[ \t]+"," ", open(tags_file, 'r').read()).split("\n")
        new_lines = []
        for directory in docker_directories:
            process2 = subprocess.Popen(['git-log', '-n', '1', '--pretty=format:%H', '--', directory],
                                        stdout=subprocess.PIPE, 
                                        stderr=subprocess.PIPE)
            commit_bytes, stderr = process2.communicate()
            process2.kill()
            if "/" in directory:
                # unix based
                service = directory.split('/')[-1]
            else:
                # windows
                service = directory.split('\\')[-1]
            commit = commit_bytes.decode('utf-8')
            new_line = service + "_tag = \"" + commit + "\""
            if new_line in original_file_lines:
                original_file_lines.remove(new_line)
            new_line = new_line + "\n"
            new_lines.append(new_line)
        
        print(str(len(original_file_lines) - 1) + " services changed")

        if len(original_file_lines) == 1: # only the empty string is present
            # the file has not changed
            return 0
        else:
            # the file requires changes
            with open(tags_file, 'w+') as file:
                file.writelines(new_lines)
            file.close()
            return 1

if __name__ == "__main__":
    main()
