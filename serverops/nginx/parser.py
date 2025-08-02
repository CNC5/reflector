from typing import List


class NginxDumpableConfig:
    params: dict

    def __init__(self):
        self.params = dict()

    def dump(self) -> dict:
        dumpdict = dict()
        for key, value in vars(self).items():
            if key == "params":
                dumpdict.update(value)
                continue
            if isinstance(value, list):
                dumpdict[key] = [i.dump() for i in value]
                continue
        return dumpdict

    def dump_formatted(self, indent: int = 4) -> List[str]:
        lines = []
        for k, v in self.params.items():
            lines.append(f"{k} {v};")
        for key, value in vars(self).items():
            if key == "params":
                continue
            if hasattr(value, "dump_formatted"):
                lines.append(key + " {")
                for line in value.dump_formatted():
                    lines.append(" " * indent + line)
                lines.append("}")
            if isinstance(value, list):
                for item in value:
                    lines.append(key + " {")
                    for line in item.dump_formatted(indent):
                        lines.append(" " * indent + line)
                    lines.append("}")
                continue
            if isinstance(value, dict):
                for k, v in value.items():
                    lines.append(key + " " + k + " {")
                    for line in v.dump_formatted(indent):
                        lines.append(" " * indent + line)
                    lines.append("}")
                continue

        return lines

    def set_param(self, key, value):
        self.params[key] = value
        return self
