from typing import List, Dict

from .parser import XrayDumpableConfig


class DialFields:
    detour: str
    bind_interface: str
    inet4_bind_address: str
    inet6_bind_address: str
    routing_mark: int
    reuse_addr: bool
    netns: str
    connect_timeout: str
    tcp_fast_open: bool
    tcp_multi_path: bool
    udp_fragment: bool
    domain_resolver: str | dict
    network_strategy: str
    network_type: List[str]
    fallback_network_type: List[str]
    fallback_delay: str


class ExternalAccountFields(XrayDumpableConfig):
    key_id: str = ""
    mac_key: str = ""


class DNS01ChallengeFields(XrayDumpableConfig):
    provider: str = "cloudflare"
    api_token: str = ""


class ACMEFields(XrayDumpableConfig):
    domain: List[str]
    data_directory: str = ""
    default_server_name: str = ""
    email: str = ""
    provider: str = ""
    disable_http_challenge: bool = False
    disable_tls_alpn_challenge: bool = False
    alternative_http_port: int = 0
    alternative_tls_port: int = 0
    external_account = ExternalAccountFields()
    dns01_challenge = DNS01ChallengeFields()


class RealityHandshakeFields(DialFields, XrayDumpableConfig):
    server: str = "google.com"
    server_port: int = 443


class TLSInbound(XrayDumpableConfig):
    enabled: bool = True
    server_name: str
    alpn: List[str]
    min_version: str
    max_version: str
    cipher_suites: List[str]
    certificate: List[str]
    certificate_path: str
    key: List[str]
    key_path: str
    acme: ACMEFields

    class ECHFields(XrayDumpableConfig):
        enabled: bool = False
        key: List[str]
        key_path: str
    ech: ECHFields

    class RealityFields(XrayDumpableConfig):
        enabled: False
        handshake: RealityHandshakeFields
        private_key: str
        short_id: List[str]
    reality: RealityFields
    max_time_difference: str = "1m"


class TLSOutbound(XrayDumpableConfig):
    enabled: bool = True
    disable_sni: bool = False
    server_name: str
    insecure: bool = False
    alpn: List[str]
    min_version: str
    max_version: str
    cipher_suites: List[str]
    certificate: str
    certificate_path: str
    fragment: bool = False
    fragment_fallback_delay: str
    record_fragment: bool = False

    class ECHFields(XrayDumpableConfig):
        enabled: False
        config: List[str]
        config_path: str
    ech: ECHFields

    class UTLSFields(XrayDumpableConfig):
        enabled: bool = False
        fingerprint: str
    utls: UTLSFields

    class RealityFields(XrayDumpableConfig):
        enabled: bool = False
        public_key: str
        short_id: str
    reality: RealityFields


class ListenFields:
    listen: str
    listen_port: int
    bind_interface: str
    routing_mark: int
    reuse_addr: bool = False
    netns: str
    tcp_fast_open: bool = False
    tcp_multi_path: bool = False
    udp_fragment: bool = False
    udp_timeout: str
    detour: str


class TCPBrutal(XrayDumpableConfig):
    enabled: bool = True
    up_mbps: int
    down_mbps: int


class MultiplexInbound(XrayDumpableConfig):
    enabled: bool = True
    padding: bool = False
    brutal: TCPBrutal


class MultiplexOutbound(XrayDumpableConfig):
    enabled: bool = True
    protocol: str
    max_connections: int
    min_streams: int
    max_streams: int
    padding: bool = False
    brutal: TCPBrutal


class V2RayTransport:
    type: str

    class HTTP(XrayDumpableConfig):
        host: List[str]
        path: str
        method: str
        headers: Dict[str, str]
        idle_timeout: str
        ping_timeout: str

    class WebSocket(XrayDumpableConfig):
        path: str
        headers: Dict[str, str]
        max_early_data: int
        early_data_header_name: str

    class QUIC(XrayDumpableConfig):
        pass
        # No additional encryption support:
        # It's basically duplicate encryption.
        # And Xray-core is not compatible with v2ray-core in here.

    class GRPC(XrayDumpableConfig):
        service_name: str
        idle_timeout: str
        ping_timeout: str
        permit_without_stream: bool = False

    class HTTPUpgrade(XrayDumpableConfig):
        host: str
        path: str
        headers: Dict[str, str]

    _type_mapping = {
        "http": HTTP,
        "ws": WebSocket,
        "quic": QUIC,
        "grpc": GRPC,
        "httpupgrade": HTTPUpgrade
    }

    def __new__(cls,
                dns_server_type: str,
                dns_server_tag: str,
                *args, **kwargs):
        new_class = cls._type_mapping[dns_server_type]()
        new_class.type = dns_server_type
        return new_class


class UDPoverTCP(XrayDumpableConfig):
    enabled: bool
    version: int
