# CIDR Calculator

A simple, interactive CIDR (Classless Inter-Domain Routing) calculator that runs in your terminal.

![Screenshot](./screenshot.png)

## Installation

You can download the latest pre-compiled binaries for your system from the [GitHub Releases](https://github.com/rzippert/cidr/releases) page.

Alternatively, you can install it using `go install`:

```sh
go install github.com/rzippert/cidr@latest
```

## Usage

Run the `cidr` command to start the interactive calculator:

```sh
cidr
```

You will be prompted to enter an IP address and mask bits. The calculator will then display the following information:

*   CIDR Netmask
*   Wildcard Mask
*   Total Addresses
*   Maximum Addresses
*   CIDR Network (Route)
*   Net: CIDR Notation
*   CIDR Address Range

To quit the application, press `q` or `ctrl+c`.

You can also get the version information by running:

```sh
cidr -v
```
or
```sh
cidr --version
```

## Building from Source

To build the `cidr` calculator from source, you will need to have Go installed.

1.  Clone the repository:

    ```sh
    git clone https://github.com/rzippert/cidr.git
    ```

2.  Change into the directory:

    ```sh
    cd cidr
    ```

3.  Build the binary:

    ```sh
    go build
    ```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.