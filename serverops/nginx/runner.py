import os.path
import subprocess
from os import PathLike
from typing import List

from .config import NginxConfig, generate_config

class Nginx:
    nginx_binary_path: PathLike[str]
    temp_config_filename: str
    temp_config: str
    prefix: PathLike[str]

    def __init__(self, nginx_binary_path: PathLike[str], tmp_dir: PathLike[str]):
        self.nginx_binary_path = nginx_binary_path
        self.prefix = tmp_dir
        self.temp_config_filename = "nginx.conf"
        self.temp_config = os.path.join(self.prefix, self.temp_config_filename)
        if not os.path.isdir(self.prefix):
            os.mkdir(self.prefix, mode=0o700)

    def assemble_args(self) -> List[str]:
        return [
            "-c", self.temp_config,
            "-e", "/dev/stderr",
            "-p", self.prefix
        ]

    def generate_config(self, config: NginxConfig):
        generate_config(self.temp_config, config)

    def serve(self, config: NginxConfig) -> subprocess.Popen:
        generate_config(self.temp_config, config)
        result = subprocess.Popen(
            [self.nginx_binary_path, *self.assemble_args()],
            shell=False) #, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
        return result

