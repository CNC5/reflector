import logging


class XrayDumpableConfig:
    def dump(self) -> dict:
        dumpdict = dict()
        for key, value in vars(self).items():
            if hasattr(value, "dump"):
                dumpdict[key] = value.dump()
                continue
            if isinstance(value, list):
                try:
                    dumpdict[key] = []
                    for i in value:
                        if hasattr(i, "dump"):
                            dumpdict[key].append(i.dump())
                        else:
                            dumpdict[key].append(i)
                    continue
                except TypeError as e:
                    logging.error(f"parsing failed on: {value}")
                    raise e
            if isinstance(value, dict):
                valuedict = dict([[k, v.dump()] for k, v in value.items()])
                dumpdict[key] = valuedict
                continue
            if value is None:
                continue
            dumpdict[key] = value
        return dumpdict
