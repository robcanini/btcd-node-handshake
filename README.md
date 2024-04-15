# btcd-node-handshake

This project implements a [btcd](https://github.com/btcsuite/btcd) node handshake in Golang. It allows you to establish a network handshake with a btcd node.

## Installation

To build the project, simply run the following command:

```bash
$ make build
```

This will compile the project and generate an executable file.

## Usage

Before launching the project, you need to create a configuration file named `config.yml`. You can find a sample configuration file inside the `config` folder. Copy this sample configuration file and customize it according to your requirements.

To launch the project, run the generated executable and pass the path to your configuration file using the `config` argument:

```bash
$ ./btcd-node-handshake --config=config/config.yml
```


Replace `config/config.yml` with the path to your actual configuration file.

## Configuration

The configuration file `config.yml` allows you to specify various settings for the btcd node handshake, including network parameters, server settings, and logging options. Make sure to customize the configuration file according to your needs.

## Testing

To test the handshake with a btcd node locally, you need to launch a btcd node with debugging enabled (please follow the related [installation instructions](https://github.com/btcsuite/btcd?tab=readme-ov-file#installation)). Run the following command:

```bash
$ btcd --debuglevel=debug
```

After launching the btcd node, launch the btcd-node-handshake program as described above. Upon successful handshake completion, you should see
a bunch of logs in the btcd-node-handshake detailing the exchanged comm messages and a "Handshake completed" log eventually.
<br/><br/>
In the btcd node instead, you should read logs as described below:

```bash
2024-04-15 20:12:44.425 [DBG] PEER: Received version (agent /btcwire:0.5.0/, pver 70016, block 212672) from 127.0.0.1:64590 (inbound)
2024-04-15 20:12:44.425 [DBG] PEER: Negotiated protocol version 70016 for peer 127.0.0.1:64590 (inbound)
2024-04-15 20:12:44.425 [DBG] CHAN: Added time sample of 0s (total: 7)
2024-04-15 20:12:44.425 [DBG] CHAN: New time offset: 0s
2024-04-15 20:12:44.425 [DBG] PEER: Sending version (agent /btcwire:0.5.0/btcd:0.24.2/, pver 70016, block 11111) to 127.0.0.1:64590 (inbound)
2024-04-15 20:12:44.426 [DBG] PEER: Sending sendaddrv2 to 127.0.0.1:64590 (inbound)
2024-04-15 20:12:44.426 [DBG] PEER: Sending verack to 127.0.0.1:64590 (inbound)
2024-04-15 20:12:44.426 [DBG] PEER: Received verack from 127.0.0.1:64590 (inbound)
2024-04-15 20:12:44.426 [DBG] PEER: Connected to 127.0.0.1:64590
2024-04-15 20:12:44.426 [DBG] SRVR: New peer 127.0.0.1:64590 (inbound)
2024-04-15 20:12:44.426 [INF] SYNC: New valid peer 127.0.0.1:64590 (inbound) (/btcwire:0.5.0/)
```
which means that the btcd node has accepted the handshake program as a new peer, and it's ready for post-handshake communication (out of scope of this implementation).
You should also read these logs:

```bash
2024-04-15 20:12:44.426 [INF] SYNC: Lost peer 127.0.0.1:64590 (inbound)
2024-04-15 20:12:44.426 [DBG] SRVR: Removed peer 127.0.0.1:64590 (inbound)
```
because the btcd-node-handshake does not keep the connection alive.

## Dependencies

This project relies on the following dependencies:

- [github.com/spf13/viper](https://github.com/spf13/viper) for configuration management.
- [github.com/rs/zerolog](https://github.com/rs/zerolog) for logging.

These dependencies will be automatically downloaded and managed using Go modules.

## Contributing

Contributions are welcome! If you find any issues or have suggestions for improvements, please feel free to open an issue or submit a pull request.

## License

This project is licensed under the [MIT License](LICENSE).
