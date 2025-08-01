import logging
import os
import signal
import subprocess
import time
from copy import deepcopy
from os import PathLike
from typing import Dict, Set, List

from .camo import CamoController
from .certificates import Certificate, prepare_certificate
from .config import ConfigV1, load_config
from .nginx import Nginx, NginxConfig
from .utils import (
    check_existence_on_disk,
    find_free_port,
    reset_allocated_ports)
from .xray import Xray, XrayConfig


class Operator:
    config: ConfigV1
    port: str

    nginx_config: NginxConfig
    nginx: Nginx
    nginx_pid: int
    camo_controller: CamoController

    xray_config: XrayConfig
    xray: Xray
    xray_pid: int

    pid_file: PathLike[str]
    inbound_ports: List[str]

    def __init__(self,
                 config_location: PathLike[str],
                 tmp_dir: PathLike[str],
                 nginx_bin: PathLike[str],
                 xray_bin: PathLike[str],
                 camo_dir: PathLike[str],
                 pid_file: PathLike[str],
                 nginx_user: str = "www-data"):

        check_existence_on_disk(
            config_location,
            nginx_bin,
            xray_bin,
            camo_dir
        )

        os.makedirs(tmp_dir, exist_ok=True)
        self.config_location = config_location
        self.config = load_config(self.config_location)
        self.nginx = Nginx(nginx_bin, tmp_dir)
        self.nginx_config = NginxConfig()
        self.camo_controller = CamoController(camo_dir, tmp_dir, nginx_user)
        self.xray = Xray(xray_bin, tmp_dir)
        self.xray_config = XrayConfig()
        self.pid_file = os.path.join(tmp_dir, pid_file)
        self.inbound_ports = []

        # outbound_tag -> {users}
        # '*' user marks outbounds that can accept any user
        # this is used for route mapping and validation
        self.outbounds_to_users_mapping: Dict[str, Set[str]] = {}

    def _add_user_mapping(self, outbound: str, user: str):
        if outbound not in self.outbounds_to_users_mapping.keys():
            self.outbounds_to_users_mapping.update({outbound: set()})
        self.outbounds_to_users_mapping[outbound].add(user)

    def configure_nginx_local_camo(self,
                                   camo_template_name: str,
                                   local_camo_https_port: int,
                                   ssl_domain: str,
                                   ssl_certificate: Certificate):
        camo_path = self.camo_controller.prepare_camo(camo_template_name)
        self.nginx_config.http.add_static_server(
            server_name=ssl_domain,
            listen=f"127.0.0.1:{local_camo_https_port}",
            root=camo_path,
            ssl_certificate=ssl_certificate)
        logging.getLogger(__name__).debug(
            f"{camo_template_name} camo configured")

    def configure_nginx_xray_proxy(self,
                                   listen: str,
                                   proxy_pass: str,
                                   ssl_domain: str,
                                   ssl_certificate: Certificate,
                                   proxy_ssl_name: str):
        self.nginx_config.http.add_proxy_server(
            server_name=ssl_domain,
            listen=listen,
            proxy_pass=proxy_pass,
            ssl_certificate=ssl_certificate,
            proxy_ssl_name=proxy_ssl_name)

    # Putting parsing here is bad, TODO
    def parse_config(self):
        config = self.config
        for conf_inb in config.spec.inbounds:
            local_camo_port = find_free_port()

            # NGX
            certificate = prepare_certificate(
                conf_inb.camo.issuer.type,
                conf_inb.camo.fqdn,
                conf_inb.camo.issuer.email
                if conf_inb.camo.issuer.type == "letsencrypt" else None
            )
            logging.getLogger(__name__).debug(
                f"processing inbound {conf_inb.name}")
            self.configure_nginx_local_camo(
                camo_template_name=conf_inb.camo.template,
                local_camo_https_port=local_camo_port,
                ssl_domain=conf_inb.camo.fqdn,
                ssl_certificate=certificate
            )
            self.inbound_ports.append(str(conf_inb.listen_port))

            # XRAY
            new_inb = self.xray_config.add_inbound(
                inbound_type=conf_inb.type,
                inbound_tag=conf_inb.name)
            if conf_inb.type == "vless":
                new_inb.listen = "127.0.0.1"
                new_inb.listen_port = conf_inb.listen_port
                tls = new_inb.tls = new_inb.__annotations__["tls"]()
                tls.enabled = True
                tls.server_name = conf_inb.camo.fqdn
                reality = tls.reality = tls.__annotations__["reality"]()
                reality.enabled = True
                reality.handshake = reality.__annotations__["handshake"]()
                reality.handshake.server = "localhost"
                reality.handshake.server_port = local_camo_port
                reality.private_key = conf_inb.private_key
                new_inb.users = []
                reality.short_id = []
                for user in conf_inb.users:
                    new_user = new_inb.User()
                    new_user.name = user.name
                    new_user.uuid = user.uuid
                    new_user.flow = user.flow
                    new_inb.users.append(new_user)
                    reality.short_id.append(user.short_id)

        for conf_out in config.spec.outbounds:
            logging.getLogger(__name__).debug(
                f"processing outbound {conf_out.name}")
            if conf_out.type == "link":
                new_out = self.xray_config.add_outbound_from_link(
                    conf_out.link,
                    conf_out.name)
                self._add_user_mapping(new_out.tag, "*")
            elif conf_out.type == "vless":
                for user in conf_out.users:
                    out_name = f"{user.name}@{conf_out.name}"
                    new_out = self.xray_config.add_outbound(
                        conf_out.type,
                        out_name)
                    new_out.server = conf_out.server
                    new_out.server_port = conf_out.server_port
                    new_out.uuid = user.uuid
                    new_out.flow = user.flow
                    tls = new_out.tls = new_out.__annotations__["tls"]()
                    tls.enabled = True
                    tls.disable_sni = False
                    tls.server_name = conf_out.server_name \
                        if conf_out.server_name else conf_out.server
                    if conf_out.fingerprint is not None:
                        utls = tls.utls = tls.__annotations__["utls"]()
                        utls.enabled = True
                        utls.fingerprint = conf_out.fingerprint
                    tls.insecure = False
                    reality = tls.reality = tls.__annotations__["reality"]()
                    reality.enabled = True
                    reality.public_key = conf_out.public_key
                    reality.short_id = user.short_id
                    new_out.packet_encoding = ""
                    self._add_user_mapping(conf_out.name, user.name)
            elif conf_out.type == "direct":
                new_out = self.xray_config.add_outbound(
                    conf_out.type,
                    conf_out.name)
                self._add_user_mapping(new_out.tag, "*")

        for conf_rou in config.spec.routes:
            logging.getLogger(__name__).debug(
                f"processing route for user {conf_rou.user}")
            possible_users = self.outbounds_to_users_mapping.get(
                conf_rou.outbound)
            if possible_users is None:
                logging.getLogger(__name__).info(
                    f"outbound {conf_rou.outbound} was not found, "
                    f"skipping route for user {conf_rou.user}")
                continue
            if (
                    conf_rou.user not in possible_users) and (
                    "*" not in possible_users):
                logging.getLogger(__name__).info(
                    f"{conf_rou.user} user is not allowed to egress from "
                    f"{conf_rou.outbound} outbound and "
                    f"outbound is not wildcard")
                continue
            if "*" in possible_users:
                new_rou = self.xray_config.add_route(conf_rou.outbound)
                new_rou.auth_user.append(conf_rou.user)
                logging.getLogger(__name__).debug(
                    f"route created '{conf_rou.user}'->'{conf_rou.outbound}'")
                continue
            if conf_rou.user in possible_users:
                new_rou = self.xray_config.add_route(
                    f"{conf_rou.user}@{conf_rou.outbound}")
                new_rou.auth_user.append(conf_rou.user)
                logging.getLogger(__name__).debug(
                    f"route created '{conf_rou.user}'->'{conf_rou.outbound}'")
                continue
            logging.getLogger(__name__).info(
                f"route did not match any criteria "
                f"(user: {conf_rou.user}; outbound: {conf_rou.outbound})")

    def validate_updated_config(self) -> bool:
        prev_config = deepcopy(self.config)
        prev_nginx_config = deepcopy(self.nginx_config)
        prev_xray_config = deepcopy(self.xray_config)
        valid = False
        self.nginx_config = NginxConfig()
        self.xray_config = XrayConfig()
        try:
            self.config = load_config(self.config_location)
            self.parse_config()
            logging.getLogger(__name__).debug("config validation succeeded")
            valid = True
        except Exception as e:
            logging.getLogger(__name__).error(f"config validation failed: {e}")
            valid = False
        finally:
            self.config = prev_config
            self.nginx_config = prev_nginx_config
            self.xray_config = prev_xray_config
            return valid

    def sighup_handler(self, signum: int, frame):
        if not self.validate_updated_config():
            logging.getLogger(__name__).error("reload aborted")
            return
        self.config = load_config(self.config_location)
        self.nginx_config = NginxConfig()
        self.xray_config = XrayConfig()
        reset_allocated_ports()
        # TODO: Only reload nginx if it's config actually changed
        #  (all changes are likely to always affect sing-box but not nginx)
        self.parse_config()
        self.nginx.generate_config(self.nginx_config)
        self.xray.generate_config(self.xray_config)
        os.kill(self.nginx_pid, signal.SIGHUP)
        os.kill(self.xray_pid, signal.SIGHUP)
        logging.getLogger(__name__).info("sighup reload succeeded")

    def run(self):
        nginx_process = None
        xray_process = None
        try:
            self.parse_config()
            logging.getLogger(__name__).debug("reflector booting")
            pid = os.getpid()
            logging.getLogger(__name__).debug(f"reflector pid: {pid}")
            with open(self.pid_file, "w") as pf:
                pf.write(str(pid))
            nginx_process = self.nginx.serve(self.nginx_config)
            xray_process = self.xray.serve(self.xray_config)
            self.nginx_pid = nginx_process.pid
            self.xray_pid = xray_process.pid

            boot_grace_seconds = 1
            time.sleep(boot_grace_seconds)
            logging.getLogger(__name__).info(
                "serving on ports: " + ", ".join(self.inbound_ports))
            signal.signal(signal.SIGHUP, self.sighup_handler)
            poll_interval_seconds = 0.1
            while True:
                nginx_code = nginx_process.poll()
                xray_code = xray_process.poll()
                if nginx_code is not None:
                    logging.getLogger(__name__).error(
                        f"nginx exited with {nginx_code}")
                    logging.getLogger(__name__).critical(
                        "can't operate without nginx, shutting down")
                    break
                if xray_code is not None:
                    logging.getLogger(__name__).error(
                        f"xray exited with {xray_code}")
                    logging.getLogger(__name__).critical(
                        "can't operate without xray, shutting down")
                    break
                time.sleep(poll_interval_seconds)
        except KeyboardInterrupt:
            logging.getLogger(__name__).info("shutdown on SIGINT")
        finally:
            graceful_timeout_seconds = 0.5
            if nginx_process:
                nginx_process.terminate()
                try:
                    nginx_process.wait(timeout=graceful_timeout_seconds)
                    logging.getLogger(__name__).debug(
                        "nginx graceful termination")
                except subprocess.TimeoutExpired:
                    nginx_process.kill()
                    logging.getLogger(__name__).debug("nginx kill")
            if xray_process:
                xray_process.terminate()
                try:
                    xray_process.wait(timeout=graceful_timeout_seconds)
                    logging.getLogger(__name__).debug(
                        "xray graceful termination")
                except subprocess.TimeoutExpired:
                    xray_process.kill()
                    logging.getLogger(__name__).debug("xray kill")

    def send_reload_signal(self):
        try:
            with open(self.pid_file) as pf:
                pid = pf.read().strip()
            os.kill(int(pid), signal.SIGHUP)
        except ValueError:
            logging.getLogger(__name__).critical("pid file is invalid")
