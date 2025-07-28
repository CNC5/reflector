import dataclasses
from typing import List, Dict

from .parser import XrayDumpableConfig
from .shared import ListenFields, TLSInbound, MultiplexInbound, V2RayTransport, DialFields


class User(XrayDumpableConfig):
    username: str
    password: str

class ShadowsocksDestination(XrayDumpableConfig):
    name: str
    server: str
    server_port: int
    password: str

class Inbound(XrayDumpableConfig):
    type: str
    tag: str

    class Direct(ListenFields, XrayDumpableConfig):
        network: str
        override_address: str
        override_port: int

    class Mixed(ListenFields, XrayDumpableConfig):
        users: List[User]
        set_system_proxy: bool

    class SOCKS(ListenFields, XrayDumpableConfig):
        users: List[User]

    class HTTP(ListenFields, XrayDumpableConfig):
        users: List[User]
        tls: TLSInbound
        set_system_proxy: bool

    class Shadowsocks(ListenFields, XrayDumpableConfig):
        method: str
        password: str
        class User(XrayDumpableConfig):
            name: str
            password: str
        users: List[User]
        destinations: List[ShadowsocksDestination]
        multiplex: MultiplexInbound

    class VMESS(ListenFields, XrayDumpableConfig):
        class User(XrayDumpableConfig):
            name: str
            uuid: str
            alterId: int
        users: List[User]
        tls: TLSInbound
        multiplex: MultiplexInbound
        transport: V2RayTransport

    class Trojan(ListenFields, XrayDumpableConfig):
        class User(XrayDumpableConfig):
            name: str
            password: str
        users: List[User]
        tls: TLSInbound
        class Fallback(XrayDumpableConfig):
            server: str
            server_port: int
        fallback: Fallback
        fallback_for_alpn = {
            "http/1.1": { # name not suitable for a variable
                "server": "127.0.0.1",
                "server_port": 8081
            }
        }
        multiplex: MultiplexInbound
        transport: V2RayTransport

    class Naive(ListenFields, XrayDumpableConfig):
        network: str
        users: List[User]
        tls: TLSInbound

    class Hysteria(ListenFields, XrayDumpableConfig):
        up: str
        up_mbps: int
        down: str
        down_mbps: int
        obfs: str
        class User(XrayDumpableConfig):
            name: str
            auth: str
            auth_str: str
        users: List[User]
        recv_window_conn: int
        recv_window_client: int
        max_conn_client: int
        disable_mtu_discovery: bool
        tls: TLSInbound

    class ShadowTLS(ListenFields, XrayDumpableConfig):
        version: int
        password: str
        class User(XrayDumpableConfig):
            name: str
            password: str
        users: List[User]
        class Handshake(DialFields, XrayDumpableConfig):
            server: str
            server_port: int
        handshake: Handshake
        class HandshakeServerName(DialFields, XrayDumpableConfig):
            server: str
            server_port: str
        handshake_for_server_name: Dict[str, HandshakeServerName]
        strict_mode: bool
        wildcard_sni: str

    class TUIC(ListenFields, XrayDumpableConfig):
        class User(XrayDumpableConfig):
            name: str
            uuid: str
            password: str
        users: List[User]
        congestion_control: str
        auth_timeout: str
        zero_rtt_handshake: bool
        heartbeat: str
        tls: TLSInbound

    class Hysteria2(ListenFields, XrayDumpableConfig):
        up_mbps: int
        down_mbps: int
        class OBFS(XrayDumpableConfig):
            type: str
            password: str
        obfs = OBFS()
        class User(XrayDumpableConfig):
            name: str
            password: str
        users: List[User]
        ignore_client_bandwidth: bool
        tls: TLSInbound
        class Masquerade(XrayDumpableConfig):
            type: str
            directory: str
            url: str
            rewrite_host: str
            status_code: str
            headers: Dict[str, str]
            content: str
        masquerade: str | Masquerade
        brutal_debug: bool

    class VLESS(ListenFields, XrayDumpableConfig):
        class User(XrayDumpableConfig):
            name: str
            uuid: str
            flow: str
        users: List[User]
        tls: TLSInbound
        multiplex: MultiplexInbound
        transport: V2RayTransport

    class AnyTLS(ListenFields, XrayDumpableConfig):
        class User(XrayDumpableConfig):
            name: str
            password: str
        users: List[User]
        padding_scheme: List[str]
        tls: TLSInbound

    class Tun(ListenFields, XrayDumpableConfig):
        interface_name: str
        address: List[str]
        mtu: int
        auto_route: bool
        iproute2_table_index: int
        iproute2_rule_index: int
        auto_redirect: bool
        auto_redirect_input_mark: str
        auto_redirect_output_mark: str
        loopback_address: List[str]
        strict_route: bool
        route_address: List[str]
        route_exclude_address: List[str]
        route_address_set: List[str]
        route_exclude_address_set: List[str]
        endpoint_independent_nat: bool
        udp_timeout: str
        stack: str
        include_interface: List[str]
        exclude_interface: List[str]
        include_uid: List[int]
        include_uid_range: List[str]
        exclude_uid: List[int]
        exclude_uid_range: List[str]
        include_android_user: List[int]
        include_package: List[str]
        exclude_package: List[str]
        class Platform(XrayDumpableConfig):
            class HTTPProxy(XrayDumpableConfig):
                enabled: bool
                server: str
                server_port: int
                bypass_domain: []
                match_domain: []
            http_proxy: HTTPProxy
        platform: Platform

    class Redirect(ListenFields, XrayDumpableConfig):
        pass

    class TProxy(ListenFields, XrayDumpableConfig):
        network: str

    _type_mapping = {
        "direct": Direct,
        "mixed": Mixed,
        "socks": SOCKS,
        "http": HTTP,
        "shadowsocks": Shadowsocks,
        "vmess": VMESS,
        "trojan": Trojan,
        "naive": Naive,
        "hysteria": Hysteria,
        "shadowtls": ShadowTLS,
        "tuic": TUIC,
        "hysteria2": Hysteria2,
        "vless": VLESS,
        "anytls": AnyTLS,
        "tun": Tun,
        "redirect": Redirect,
        "tproxy": TProxy,
    }

    def __new__(cls, inbound_type: str, inbound_tag: str, *args, **kwargs):
        new_class = cls._type_mapping[inbound_type]()
        new_class.type = inbound_type
        new_class.tag = inbound_tag
        return new_class

