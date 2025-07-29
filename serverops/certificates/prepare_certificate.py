import dataclasses
import logging
import os.path
import subprocess
from os import PathLike
from cryptography import x509
from cryptography.hazmat.backends import default_backend
from datetime import datetime, timezone


@dataclasses.dataclass
class Certificate:
    certificate_path: PathLike | str
    key_path: PathLike | str


class IssueError(Exception):
    pass


def days_until_certificate_expiry(cert_path) -> int:
    with open(cert_path, "rb") as f:
        cert_data = f.read()
    cert = x509.load_pem_x509_certificate(cert_data, default_backend())
    expires = cert.not_valid_after_utc
    now = datetime.now(timezone.utc)
    days_left = (expires - now).days
    return days_left


def prepare_certificate(certificate_type: str,
                        domain: str,
                        email: str = None) -> Certificate:
    match certificate_type:
        case "selfsigned":
            return _prepare_selfsigned_certificate(domain)
        case "letsencrypt":
            if email is None:
                raise IssueError(
                    "email is required to issue a letsencrypt certificate")
            return _prepare_letsencrypt_certificate(domain, email)
    raise Exception("Unknown certificate issuer type")


def _prepare_selfsigned_certificate(domain: str) -> Certificate:
    directory = f"/tmp/test-cert/{domain}/"
    directory = os.path.realpath(directory)
    os.makedirs(directory, exist_ok=True)
    certificate = os.path.join(directory, "cert.pem")
    key = os.path.join(directory, "key.pem")
    if (
            os.path.isfile(key)) and (
            os.path.isfile(certificate)) and (
            days_until_certificate_expiry(certificate) > 15):
        logging.getLogger(__name__).debug(
            "certificate exists and is up to date, not reissuing")
        return Certificate(certificate_path=certificate, key_path=key)
    country_name = "NE"
    state_name = "StateName"
    city_name = "CityName"
    company_name = "selfsigner.org"
    company_section_name = "sso"
    domain_name = domain
    logging.getLogger(__name__).debug(
        f"issuing new selfsigned certificate for {domain_name}")
    issuing = subprocess.Popen([
        "openssl",
        "req", "-x509", "-newkey", "rsa:4096",
        "-keyout", key, "-out", certificate,
        "-sha256", "-days", "3650", "-nodes",
        "-subj", f"/C={country_name}"
                 f"/ST={state_name}"
                 f"/L={city_name}"
                 f"/O={company_name}"
                 f"/OU={company_section_name}"
                 f"/CN={domain_name}"])
    issuing.wait()
    return Certificate(certificate_path=certificate, key_path=key)


def _prepare_letsencrypt_certificate(domain: str, email: str) -> Certificate:
    directory = f"/etc/letsencrypt/live/{domain}/"
    directory = os.path.realpath(directory)
    certificate = os.path.join(directory, "fullchain.pem")
    key = os.path.join(directory, "privkey.pem")
    if (
            os.path.isfile(key)) and (
            os.path.isfile(certificate)) and (
            days_until_certificate_expiry(certificate) > 15):
        logging.getLogger(__name__).debug(
            "certificate exists and is up to date, not reissuing")
        return Certificate(certificate_path=certificate, key_path=key)
    logging.getLogger(__name__).debug(
        f"issuing new letsencrypt certificate for {domain}")
    certbot_process = subprocess.Popen([
        "certbot", "certonly",
        "--standalone",
        "--preferred-challenges", "http",
        "-d", domain,
        "-m", email,
        "--agree-tos",
        "--non-interactive",
        "--keep-until-expiring"
    ])
    exit_code = certbot_process.wait()
    if exit_code != 0:
        raise IssueError(
            f"certbot failed, code: {exit_code}, "
            f"{certbot_process.stdout.read()}")
    return Certificate(certificate_path=certificate, key_path=key)
