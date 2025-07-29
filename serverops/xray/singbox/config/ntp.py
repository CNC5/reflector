from .shared import DialFields


class NTP(DialFields):
    enabled: bool = False
    server: str
    server_port: int
    interval: str
