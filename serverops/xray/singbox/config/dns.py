from os import PathLike
from typing import List, Dict

from .shared import DialFields, TLSOutbound
from .parser import XrayDumpableConfig

class DNSServer(XrayDumpableConfig):
    type: str
    tag: str

    class Local(DialFields, XrayDumpableConfig):
        pass

    class Hosts(DialFields, XrayDumpableConfig):
        path: List[PathLike]
        predefined: dict

    class TCP(DialFields, XrayDumpableConfig):
        server: str
        server_port: int

    class UDP(DialFields, XrayDumpableConfig):
        server: str
        server_port: int

    class TLS(DialFields, XrayDumpableConfig):
        server: str
        server_port: int
        tls = TLSOutbound()

    class QUIC(DialFields, XrayDumpableConfig):
        server: str
        server_port: int
        tls = TLSOutbound()

    class HTTPS(DialFields, XrayDumpableConfig):
        server: str
        server_port: int
        path: str
        headers: Dict[str, str]
        tls = TLSOutbound()

    class HTTP3(DialFields, XrayDumpableConfig):
        server: str
        server_port: int
        path: str
        headers: Dict[str, str]
        tls = TLSOutbound()

    class DHCP(DialFields, XrayDumpableConfig):
        interface: str

    class FakeIP:
        inet4_range: str
        inet6_range: str

    class Tailscale:
        endpoint: str
        accept_default_resolvers: bool = False

    class Resolved:
        service: str
        accept_default_resolvers: bool = False

    _type_mapping = {
        "local": Local,
        "hosts": Hosts,
        "tcp": TCP,
        "udp": UDP,
        "tls": TLS,
        "quic": QUIC,
        "https": HTTPS,
        "h3": HTTP3,
        "dhcp": DHCP,
        "fakeip": FakeIP,
        "tailscale": Tailscale,
        "resolved": Resolved
    }

    def __new__(cls, dns_server_type: str, dns_server_tag: str, *args, **kwargs):
        new_class = cls._type_mapping[dns_server_type]()
        new_class.type = dns_server_type
        new_class.tag = dns_server_tag
        return new_class

class DNSRule(XrayDumpableConfig):
    inbound: List[str]
    ip_version: int
    query_type: str | int
    network: str
    auth_user: [str]
    protocol: [str]
    domain: [str]
    domain_suffix: [str]
    domain_keyword: [str]
    domain_regex: [str]
    source_ip_cidr: [str]
    source_ip_is_private: bool = False
    ip_cidr: [str]
    ip_is_private: bool = False
    ip_accept_any: bool = False
    source_port: [int]
    source_port_range: [str]
    port: [int]
    port_range: [str]
    process_name: [str]
    process_path: [str]
    process_path_regex: [str]
    package_name: [str]
    user: [str]
    user_id: [int]
    clash_mode: str
    network_type: [str]
    network_is_expensive: bool = False
    network_is_constrained: bool = False
    wifi_ssid: [str]
    wifi_bssid: [str]
    rule_set: [str]
    rule_set_ip_cidr_match_source: bool = False
    rule_set_ip_cidr_accept_empty: bool = False
    invert: bool = False
    outbound: [str]
    action: str
    server: str

class DNS(XrayDumpableConfig):
    servers: List[DNSServer]
    rules: List[DNSRule]
    final: str
    strategy: str
    disable_cache: bool
    disable_expire: bool
    independent_cache: bool
    cache_capacity: int
    reverse_mapping: bool
    client_subnet: str

    def __init__(self):
        self.servers = list()
        self.servers.append(DNSServer("local", "default"))