from typing import List

from .parser import XrayDumpableConfig


class RouteRule(XrayDumpableConfig):
    inbound: List[str]
    ip_version: int
    network: List[str]
    auth_user: List[str]
    protocol: List[str]
    client: List[str]
    domain: List[str]
    domain_suffix: List[str]
    domain_keyword: List[str]
    domain_regex: List[str]
    geosite: List[str]
    source_geoip: List[str]
    geoip: List[str]
    source_ip_cidr: List[str]
    source_ip_is_private: bool
    ip_cidr: List[str]
    ip_is_private: bool
    source_port: List[int]
    source_port_range: List[str]
    port: List[int]
    port_range: List[str]
    process_name: List[str]
    process_path: List[str]
    process_path_regex: List[str]
    package_name: List[str]
    user: List[str]
    user_id: List[int]
    clash_mode: str
    network_type: List[str]
    network_is_expensive: bool
    network_is_constrained: bool
    wifi_ssid: List[str]
    wifi_bssid: List[str]
    rule_set: List[str]
    # DEPRECATED
    outbound: str


class RouteRuleLogical(XrayDumpableConfig):
    type: str
    mode: str
    rules: List[RouteRule]
    invert: bool
    action: str
    outbound: str


class RuleSet(XrayDumpableConfig):
    type: str
    tag: str
    rules: List[RouteRule | RouteRuleLogical]


class Route(XrayDumpableConfig):
    rules: List[RouteRule | RouteRuleLogical]
    rule_set: List[RuleSet]
    final: str
    auto_detect_interface: bool
    override_android_vpn: bool
    default_interface: str
    default_mark: int
    default_domain_resolver: str | dict
    default_network_strategy: str
    default_network_type: List
    default_fallback_network_type: List
    default_fallback_delay: str
