from typing import List

from .parser import XrayDumpableConfig
from .shared import DialFields, ListenFields


class WireguardPeer(XrayDumpableConfig):
    address: str
    port: int
    public_key: str
    pre_shared_key: str
    allowed_ips: List[str]
    persistent_keepalive_interval: int
    reserved: List[int]


class Endpoint(XrayDumpableConfig):
    type: str
    tag: str

    class Wireguard(DialFields, XrayDumpableConfig):
        system: bool = False
        name: str
        mtu: int
        address: List[str]
        private_key: str
        listen_port: int
        peers: WireguardPeer
        udp_timeout: str
        workers: int

    class Tailscale(DialFields, XrayDumpableConfig):
        state_directory: str
        auth_key: str
        control_url: str
        ephemeral: bool = False
        hostname: str
        accept_routes: bool = False
        exit_node: str
        exit_node_allow_lan_access: bool = False
        advertise_routes: List[str]
        advertise_exit_node: bool = True
        udp_timeout: str

    _type_mapping = {
        "wireguard": Wireguard,
        "tailscale": Tailscale,
    }

    def __new__(cls, endpoint_type: str, endpoint_tag: str, *args, **kwargs):
        new_class = cls._type_mapping[endpoint_type]()
        new_class.type = endpoint_type
        new_class.tag = endpoint_tag
        return new_class
