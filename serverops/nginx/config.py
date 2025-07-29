from typing import List, Dict
from .mimetypes import load_default_mime_types
from .parser import NginxDumpableConfig
from ..certificates import prepare_certificate, Certificate


class NginxHTTPServerLocationConfig(NginxDumpableConfig):
    params: dict


class NginxHTTPServerConfig(NginxDumpableConfig):
    params: dict
    location: Dict[str, NginxHTTPServerLocationConfig]

    def __init__(self):
        self.location = dict()
        super().__init__()

    def add_static_location(self,
                            path: str,
                            content_location: str,
                            autoindex=False):
        self.location.update({path:
                              NginxHTTPServerLocationConfig()
                              .set_param("autoindex", "on"
                                         if autoindex else "off")
                              .set_param("index", "index.html")
                              .set_param("root", content_location)
                              .set_param("max_ranges", "1000000000000000")
                              })
        return self

    def add_proxy_location(self, path: str, proxy_pass: str):
        self.location.update({path:
                              NginxHTTPServerLocationConfig()
                              .set_param("proxy_pass", proxy_pass)
                              .set_param("proxy_set_header", "Host $host")
                              })
        return self

    def add_ssl(
            self,
            certificate_path: str,
            certificate_key_path: str,
            ssl_protocols: List[str] = ("TLSv1.1", "TLSv1.2", "TLSv1.3"),
            ssl_ciphers: List[str] = ("HIGH", "!aNULL", "!MD5")):
        (self
         .set_param("ssl_certificate", certificate_path)
         .set_param("ssl_certificate_key", certificate_key_path)
         .set_param("ssl_protocols", " ".join(ssl_protocols))
         .set_param("ssl_ciphers", ":".join(ssl_ciphers)))
        self.params["listen"] += " ssl"

    def add_http2(self):
        self.set_param("http2", "on")


class NginxHTTPMimeTypesConfig(NginxDumpableConfig):
    params: dict


class NginxHTTPConfig(NginxDumpableConfig):
    params: dict
    types: NginxHTTPMimeTypesConfig
    server: List[NginxHTTPServerConfig]

    def __init__(self):
        self.server = []
        self.types = NginxHTTPMimeTypesConfig()
        super().__init__()
        load_default_mime_types(self.types)
        self.set_param("access_log", "/dev/stdout")

    def add_static_server(
            self,
            server_name: str,
            listen: str,
            root: str,
            autoindex=False,
            client_max_body_size="100M",
            ssl_certificate: Certificate = None,
            http2: bool = True):
        new_server = (NginxHTTPServerConfig()
                      .set_param("server_name", server_name)
                      .set_param("listen", listen)
                      .set_param("client_max_body_size", client_max_body_size)
                      .add_static_location("/", root, autoindex))
        if ssl_certificate is not None:
            cert = ssl_certificate
            new_server.add_ssl(cert.certificate_path, cert.key_path)
        if http2:
            new_server.add_http2()
        self.server.append(new_server)
        return self

    def add_proxy_server(
            self,
            server_name: str,
            listen: str,
            proxy_pass: str,
            client_max_body_size="1G",
            ssl_certificate: Certificate = None,
            proxy_ssl_name: str = None):
        new_server = (NginxHTTPServerConfig()
                      .set_param("server_name", server_name)
                      .set_param("listen", listen)
                      .set_param("client_max_body_size", client_max_body_size)
                      .add_proxy_location("/", proxy_pass))
        if ssl_certificate is not None:
            cert = ssl_certificate
            new_server.add_ssl(cert.certificate_path, cert.key_path)
            if proxy_ssl_name is not None:
                for _, location in new_server.location.items():
                    location.set_param("proxy_ssl_name", proxy_ssl_name)
                    location.set_param("proxy_ssl_server_name", "on")
        self.server.append(new_server)
        return self


class NginxEventsConfig(NginxDumpableConfig):
    params: dict

    def __init__(self):
        super().__init__()
        self.set_param("worker_connections", "1024")


class NginxConfig(NginxDumpableConfig):
    params: dict
    events: NginxEventsConfig
    http: NginxHTTPConfig
    user: str

    def __init__(self, user: str = "www-data"):
        self.events = NginxEventsConfig()
        self.http = NginxHTTPConfig()
        self.user = user
        super().__init__()
        self.set_param("daemon", "off")
        self.set_param("user", self.user)
        self.set_param("worker_processes", "auto")
        self.set_param("error_log", "/dev/stderr")
        self.set_param("pid", "/tmp/nginx.pid")


def generate_config(location: str, config: NginxConfig):
    with open(location, "w") as f:
        f.write("\n".join(config.dump_formatted()))
