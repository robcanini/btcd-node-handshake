# btcd-node-handshake

This project implements a [BTCd](https://github.com/btcsuite/btcd) node handshake in Golang. It allows you to establish a network handshake with a BTCd node.

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

The configuration file `config.yml` allows you to specify various settings for the BTCd node handshake, including network parameters, server settings, and logging options. Make sure to customize the configuration file according to your needs.

## Dependencies

This project relies on the following dependencies:

- [github.com/spf13/viper](https://github.com/spf13/viper) for configuration management.
- [github.com/rs/zerolog](https://github.com/rs/zerolog) for logging.

These dependencies will be automatically downloaded and managed using Go modules.

## Contributing

Contributions are welcome! If you find any issues or have suggestions for improvements, please feel free to open an issue or submit a pull request.

## License

This project is licensed under the [MIT License](LICENSE).
