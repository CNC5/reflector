import argparse
import logging
import sys

from serverops import Operator


def main():
    argparser = argparse.ArgumentParser()
    argparser.add_argument(
        "--tmp", help="directory for tmp storage",
        default="/tmp/reflector/")
    argparser.add_argument(
        "-c", "--config", help="config file, in cwd",
        default="config.yaml")
    argparser.add_argument(
        "--pid-file", help="pid file, in the tmp directory",
        default="ops.pid")
    argparser.add_argument(
        "--nginx-bin", help="nginx binary",
        default="serverops/bin/nginx")
    argparser.add_argument(
        "--xray-bin", help="xray binary",
        default="serverops/bin/sing-box")
    argparser.add_argument(
        "--camo-dir", help="camo templates dir",
        default="serverops/camo/templates")
    argparser.add_argument(
        "-d", "--debug", default=False, action="store_true")
    argparser.add_argument(
        "-s", "--signal", help="send a signal to the operator",
        default=None)

    args = argparser.parse_args()
    logging.basicConfig(stream=sys.stdout,
                        level=logging.DEBUG if args.debug else logging.INFO,
                        format='%(asctime)s %(levelname)-8s '
                               '%(name)-16s: %(message)s',
                        datefmt='%Y-%m-%d %H:%M:%S')
    ops = Operator(
        config_location=args.config,
        tmp_dir=args.tmp,
        nginx_bin=args.nginx_bin,
        xray_bin=args.xray_bin,
        camo_dir=args.camo_dir,
        pid_file=args.pid_file
    )
    if args.signal is None:
        ops.run()
        return
    match args.signal:
        case "reload":
            ops.send_reload_signal()
            logging.getLogger(__name__).info("reload signal sent")
        case _:
            logging.getLogger(__name__).error("unrecognized signal")
