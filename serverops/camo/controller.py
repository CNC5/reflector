import logging
import os
import pwd
import shutil
import subprocess
from os import PathLike
from typing import Dict, List


def ensure_system_user(username: str) -> int:
    try:
        pwd.getpwnam(username)
        return 0
    except KeyError:
        logging.getLogger(__name__).info(
            f"user {username} does not exist, creating")
        if shutil.which("useradd"):
            return subprocess.Popen(["useradd", username]).wait()
        if shutil.which("adduser"):
            return subprocess.Popen([
                "adduser",
                "-D",  # Don't assign a password
                "-H",  # Don't create home directory
                username]).wait()
        logging.getLogger(__name__).error(
            "user management utils are unavailable")
        return -1


def chown_recursive(
        directory: PathLike[str],
        username: str,
        group: str = None) -> int:
    if directory in {"/", "/*"}:
        logging.getLogger(__name__).error("refusing to chown /")
        return -1
    if ensure_system_user(username) != 0:
        logging.getLogger(__name__).error(
            "chown failed as there is no requested user")
        return -1
    group = username if group is None else group
    return subprocess.Popen(
        ["chown", "-R", f"{username}:{group}", directory]).wait()


class CamoTemplate:
    template_dir: str

    def __init__(self, template_dir: str):
        self.template_dir = template_dir


class CamoController:
    available_camos: Dict[str, CamoTemplate]
    nginx_user: str

    def __init__(self,
                 templates_dir: PathLike[str],
                 tmp_dir: PathLike[str],
                 nginx_user: str):
        self.available_camos = dict()
        self.nginx_user = nginx_user
        templates_dir = os.path.realpath(templates_dir)
        tmp_dir = os.path.join(os.path.realpath(tmp_dir), "templates")
        shutil.copytree(templates_dir, tmp_dir, dirs_exist_ok=True)
        chown_recursive(tmp_dir, nginx_user)
        templates = os.listdir(tmp_dir)
        logging.getLogger(__name__).debug(f"templates: {templates}")
        for template_dir in templates:
            full_template_dir = os.path.join(tmp_dir, template_dir)
            self.available_camos.update(
                {template_dir: CamoTemplate(full_template_dir)})

    def get_available_camos(self) -> List[str]:
        return list(self.available_camos.keys())

    def prepare_camo(self, camo_name: str) -> str:
        if camo_name not in self.available_camos:
            raise Exception("camo not available")
        return self.available_camos[camo_name].template_dir
