import dataclasses
from typing import Union, Dict, List

from .endpoint import WireguardPeer
from .parser import XrayDumpableConfig
from .shared import DialFields, UDPoverTCP, TLSOutbound, MultiplexOutbound, V2RayTransport


class Outbound(XrayDumpableConfig):
    type: str
    tag: str

    class Direct(DialFields, XrayDumpableConfig):
        pass

    class Block(XrayDumpableConfig):
        pass

    class SOCKS(DialFields, XrayDumpableConfig):
        server: str
        server_port: int
        version: str
        username: str
        password: str
        network: str
        udp_over_tcp: bool | UDPoverTCP

    class HTTP(DialFields, XrayDumpableConfig):
        server: str
        server_port: int
        username: str
        password: str
        path: str
        headers: Dict[str, str]
        tls: TLSOutbound

    class Shadowsocks(DialFields, XrayDumpableConfig):
        server: str
        server_port: int
        method: str
        password: str
        plugin: str
        plugin_opts: str
        network: str
        udp_over_tcp: bool | UDPoverTCP
        multiplex: MultiplexOutbound

    class VMESS(DialFields, XrayDumpableConfig):
        server: str
        server_port: int
        uuid: str
        security: str
        alter_id: int
        global_padding: bool
        authenticated_length: bool
        network: str
        tls: TLSOutbound
        packet_encoding: str
        transport: V2RayTransport
        multiplex: MultiplexOutbound

    class Trojan(DialFields, XrayDumpableConfig):
        server: str
        server_port: int
        password: str
        network: str
        tls: TLSOutbound
        multiplex: MultiplexOutbound
        transport: V2RayTransport

    class Wireguard(DialFields, XrayDumpableConfig):
        server: str
        server_port: int
        system_interface: bool
        interface_name: str
        local_address: List[str]
        private_key: str
        class WireguardPeer(XrayDumpableConfig):
            server: str
            server_port: int
            public_key: str
            pre_shared_key: str
            allowed_ips: List[str]
            reserved: List[int]
        peers: WireguardPeer
        peer_public_key: str
        pre_shared_key: str
        reserved: List[int]
        workers: int
        mtu: int
        network: str

    class Hysteria(DialFields, XrayDumpableConfig):
        server: str
        server_port: int
        server_ports: List[str]
        hop_interval: str
        up: str
        up_mbps: int
        down: str
        down_mbps: int
        obfs: str
        auth: str
        auth_str: str
        recv_window_conn: int
        recv_window: int
        disable_mtu_discovery: bool
        network: str
        tls: TLSOutbound

    class VLESS(DialFields, XrayDumpableConfig):
        server: str
        server_port: int
        uuid: str
        flow: str
        network: str
        tls: TLSOutbound
        packet_encoding: str
        multiplex: MultiplexOutbound
        transport: V2RayTransport

    class ShadowTLS(DialFields, XrayDumpableConfig):
        server: str
        server_port: int
        version: int
        password: str
        tls: TLSOutbound

    class TUIC(DialFields, XrayDumpableConfig):
        server: str
        server_port: int
        uuid: str
        password: str
        congestion_control: str
        udp_relay_mode: str
        udp_over_stream: bool
        zero_rtt_handshake: bool
        heartbeat: str
        network: str
        tls: TLSOutbound

    class Hysteria2(DialFields, XrayDumpableConfig):
        server: str
        server_port: int
        server_ports: List[str]
        hop_interval: str
        up_mbps: int
        down_mbps: int
        class OBFS(XrayDumpableConfig):
            type: str
            password: str
        obfs: OBFS
        password: str
        network: str
        tls: TLSOutbound
        brutal_debug: bool

    class AnyTLS(DialFields, XrayDumpableConfig):
        server: str
        server_port: int
        password: str
        idle_session_check_interval: str
        idle_session_timeout: str
        min_idle_session: int
        tls: TLSOutbound

    class Tor(DialFields, XrayDumpableConfig):
        executable_path: str
        extra_args: List[str]
        data_directory: str
        torrc: Dict[str, any]

    class SSH(DialFields, XrayDumpableConfig):
        server: str
        server_port: int
        user: str
        password: str
        private_key: str
        private_key_path: str
        private_key_passphrase: str
        host_key: List[str]
        host_key_algorithms: List[str]
        client_version: str

    class Selector(XrayDumpableConfig):
        outbounds: List[str]
        default: str
        interrupt_exist_connections: bool

    class URLTest(XrayDumpableConfig):
        outbounds: List[str]
        url: str
        interval: str
        tolerance: int
        idle_timeout: str
        interrupt_exist_connections: bool

    _type_mapping = {
        "direct": Direct,
        "block": Block,
        "socks": SOCKS,
        "http": HTTP,
        "shadowsocks": Shadowsocks,
        "vmess": VMESS,
        "trojan": Trojan,
        "wireguard": Wireguard,
        "hysteria": Hysteria,
        "vless": VLESS,
        "shadowtls": ShadowTLS,
        "tuic": TUIC,
        "hysteria2": Hysteria2,
        "anytls": AnyTLS,
        "tor": Tor,
        "ssh": SSH,
        "selector": Selector,
        "urltest": URLTest,
    }

    type_union = Union[_type_mapping.values()]

    def __new__(cls, outbound_type: str, outbound_tag: str, *args, **kwargs):
        new_class = cls._type_mapping[outbound_type]()
        new_class.type = outbound_type
        new_class.tag = outbound_tag
        return new_class
