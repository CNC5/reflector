from os import PathLike

from pydantic import BaseModel, StrictStr, StrictInt, Field
from typing import List, Optional, Literal, Annotated, Union
import yaml


class User(BaseModel):
    name: StrictStr
    uuid: StrictStr
    flow: Optional[StrictStr]
    short_id: StrictStr


class CamoIssuerLetsencrypt(BaseModel):
    type: Literal["letsencrypt"]
    email: StrictStr


class CamoIssuerSelfsigned(BaseModel):
    type: Literal["selfsigned"]


CamoIssuer = Annotated[Union[
    CamoIssuerSelfsigned,
    CamoIssuerLetsencrypt
], Field(discriminator="type")]


class CamoLocal(BaseModel):
    type: Literal["local"]
    template: StrictStr
    fqdn: StrictStr
    issuer: CamoIssuer


Camo = Annotated[Union[CamoLocal], Field(discriminator="type")]


class InboundRawVless(BaseModel):
    name: StrictStr
    type: Literal["vless"]
    listen: StrictStr
    listen_port: StrictInt
    users: List[User]
    private_key: StrictStr
    camo: Camo


Inbound = Annotated[Union[InboundRawVless], Field(discriminator="type")]


class OutboundLink(BaseModel):
    name: StrictStr
    type: Literal["link"]
    link: Optional[StrictStr] = None


class OutboundRawVless(BaseModel):
    name: StrictStr
    type: Literal["vless"]
    server: StrictStr
    server_port: StrictInt
    server_name: Optional[StrictStr]
    fingerprint: Optional[StrictStr]
    users: List[User]
    public_key: StrictStr


class OutboundRawDirect(BaseModel):
    name: StrictStr
    type: Literal["direct"]


Outbound = Annotated[Union[
    OutboundLink,
    OutboundRawVless,
    OutboundRawDirect
], Field(discriminator="type")]


class Route(BaseModel):
    user: StrictStr
    outbound: StrictStr


class Metrics(BaseModel):
    port: StrictInt
    listen: StrictStr


class Spec(BaseModel):
    inbounds: List[Inbound]
    outbounds: List[Outbound]
    routes: List[Route]
    metrics: Metrics = None

class ConfigV1(BaseModel):
    apiVersion: StrictStr
    kind: StrictStr
    spec: Spec


def load_config(file_path: PathLike[str]) -> ConfigV1:
    with open(file_path, "r") as f:
        raw = yaml.safe_load(f)
    return ConfigV1(**raw)
