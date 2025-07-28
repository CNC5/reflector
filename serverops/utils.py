import os.path
from os import PathLike

import psutil

def check_existence_on_disk(*path: PathLike[str]):
    for p in path:
        if os.path.isdir(p):
            continue
        if os.path.isfile(p):
            continue
        raise FileNotFoundError(p)

_previously_allocated_ports = set()
def find_free_port(first: int = 42000, last: int = 65535) -> int:
    """
    Find a port that is not used by any program and hasn't been returned previously this run
    """
    desired_ports = set(range(first, last))
    busy_ports = set([conn.laddr.port for conn in psutil.net_connections()]).union(_previously_allocated_ports)
    usable_ports = desired_ports.difference(busy_ports)
    selected_port = usable_ports.pop()
    _previously_allocated_ports.add(selected_port)
    return selected_port

def reset_allocated_ports():
    _previously_allocated_ports.clear()

def bind(instance, func, as_name=None):
    """
    Bind the function *func* to *instance*, with either provided name *as_name*
    or the existing name of *func*. The provided *func* should accept the
    instance as the first argument, i.e. "self".
    """
    if as_name is None:
        as_name = func.__name__
    bound_method = func.__get__(instance, instance.__class__)
    setattr(instance, as_name, bound_method)
    return bound_method
