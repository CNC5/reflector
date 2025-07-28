import os
import subprocess
from os import PathLike
from typing import List

from .config import config, generate_config


class Xray:
    xray_binary_path: PathLike[str]
    temp_config_filename: str
    temp_config: str
    prefix: PathLike[str]

    def __init__(self, xray_binary_path: PathLike[str], tmp_dir: PathLike[str]):
        self.xray_binary_path = xray_binary_path
        self.prefix = tmp_dir
        self.temp_config_filename = "sing.json"
        self.temp_config = os.path.join(self.prefix, self.temp_config_filename)
        if not os.path.isdir(self.prefix):
            os.mkdir(self.prefix, mode=0o700)

    def assemble_args(self) -> List[str]:
        return [
            "run",
            "-c", self.temp_config
        ]

    def generate_config(self, config: config.XrayConfig):
        generate_config(self.temp_config, config)

    def serve(self, config: config.XrayConfig) -> subprocess.Popen:
        generate_config(self.temp_config, config)
        result = subprocess.Popen(
            [self.xray_binary_path, *self.assemble_args()],
            shell=False) #, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
        return result
