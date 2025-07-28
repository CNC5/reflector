import json
from os import PathLike
from typing import List
from urllib.parse import urlparse, parse_qs

from .certificate import Certificate
from .dns import DNS
from .endpoint import Endpoint
from .inbound import Inbound
from .ntp import NTP
from .outbound import Outbound
from .parser import XrayDumpableConfig
import os

from .route import Route, RouteRule
from .shared import TLSOutbound

class Log(XrayDumpableConfig):
    disabled: bool
    level: str
    output: str
    timestamp: bool

    def __init__(self):
        self.disabled = False
        self.level = "info"
        self.output = "xray.log"
        self.timestamp = True

class XrayConfig(XrayDumpableConfig):
    log: Log
    dns: DNS
    ntp: NTP
    certificate: Certificate
    endpoints: List[Endpoint]
    inbounds: List[Inbound]
    outbounds: List[Outbound]
    route: Route
    # TODO
    #services: List[XrayServiceConfig]
    #experimental: XrayExperimentalConfig

    def __init__(self):
        self.log = Log()
        self.log.output = "stdout"
        self.log.timestamp = True
        self.dns = DNS()
        self.endpoints = list()
        self.inbounds = list()
        self.outbounds = list()
        self.route = Route()
        self.route.rules = []

    def set_log_level(self, level: str):
        self.log.level = level

    def add_inbound(self, inbound_type: str, inbound_tag: str):
        new_inbound = Inbound(inbound_type, inbound_tag)
        self.inbounds.append(new_inbound)
        return new_inbound

    def add_outbound(self, outbound_type: str, outbound_tag: str):
        new_outbound = Outbound(outbound_type, outbound_tag)
        self.outbounds.append(new_outbound)
        return new_outbound

    def add_outbound_from_link(self, link: str, outbound_tag: str) -> Outbound:
        parsed = urlparse(link)
        params = parse_qs(parsed.query)
        outbound_type = parsed.scheme
        new_outbound = Outbound(outbound_type, outbound_tag)
        new_outbound.server = parsed.hostname
        new_outbound.server_port = parsed.port
        nb_annotations = new_outbound.__annotations__
        if "username" in nb_annotations:
            new_outbound.username = parsed.username
        else:
            new_outbound.uuid = parsed.username
        if "password" in nb_annotations:
            new_outbound.password = parsed.password
        if params["security"][0] == "none":
            return new_outbound
        tls = TLSOutbound()
        new_outbound.tls = tls
        tls.enabled = True
        if "fp" in params:
            tls.utls = tls.UTLSFields()
            tls.utls.enabled = True
            tls.utls.fingerprint = params["fp"][0]
        tls.disable_sni = False
        tls.server_name = params["sni"][0]
        tls.insecure = False
        tls.reality = tls.RealityFields()
        tls.reality.enabled = True
        tls.reality.public_key = params["pbk"][0]
        tls.reality.short_id = params["sid"][0]
        self.outbounds.append(new_outbound)
        return new_outbound

    def add_route(self, outbound_tag: str):
        new_route = RouteRule()
        new_route.outbound = outbound_tag
        new_route.auth_user = []
        self.route.rules.append(new_route)
        return new_route

def generate_config(location: str, config: XrayConfig):
        with open(location, "w") as f:
            f.write(json.dumps(config.dump()))
