---
description: This section describes the configuration parameters and their types for INX-POI.
keywords:
- IOTA Node 
- Hornet Node
- POI
- Proof-Of-Inclusion
- Proof
- Inclusion
- Configuration
- JSON
- Customize
- Config
- reference
---


# Core Configuration

INX-POI uses a JSON standard format as a config file. If you are unsure about JSON syntax, you can find more information in the [official JSON specs](https://www.json.org).

You can change the path of the config file by using the `-c` or `--config` argument while executing `inx-poi` executable.

For example:
```bash
inx-poi -c config_defaults.json
```

You can always get the most up-to-date description of the config parameters by running:

```bash
inx-poi -h --full
```

## <a id="app"></a> 1. Application

| Name            | Description                                                                                            | Type    | Default value |
| --------------- | ------------------------------------------------------------------------------------------------------ | ------- | ------------- |
| checkForUpdates | Whether to check for updates of the application or not                                                 | boolean | true          |
| stopGracePeriod | The maximum time to wait for background processes to finish during shutdown before terminating the app | string  | "5m"          |

Example:

```json
  {
    "app": {
      "checkForUpdates": true,
      "stopGracePeriod": "5m"
    }
  }
```

## <a id="inx"></a> 2. INX

| Name    | Description                            | Type   | Default value    |
| ------- | -------------------------------------- | ------ | ---------------- |
| address | The INX address to which to connect to | string | "localhost:9029" |

Example:

```json
  {
    "inx": {
      "address": "localhost:9029"
    }
  }
```

## <a id="poi"></a> 3. Proof-Of-Inclusion

| Name        | Description                                       | Type   | Default value    |
| ----------- | ------------------------------------------------- | ------ | ---------------- |
| bindAddress | Bind address on which the POI HTTP server listens | string | "localhost:9687" |

Example:

```json
  {
    "poi": {
      "bindAddress": "localhost:9687"
    }
  }
```

## <a id="profiling"></a> 4. Profiling

| Name        | Description                                       | Type    | Default value    |
| ----------- | ------------------------------------------------- | ------- | ---------------- |
| enabled     | Whether the profiling plugin is enabled           | boolean | false            |
| bindAddress | The bind address on which the profiler listens on | string  | "localhost:6060" |

Example:

```json
  {
    "profiling": {
      "enabled": false,
      "bindAddress": "localhost:6060"
    }
  }
```

