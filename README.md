# :white_square_button: mQR • My QR Codes in Terminal

Generate QR codes using terminal

<img src="docs/demo.gif" alt="mQR demo example">

## :gear: Install

```shell
go install github.com/mymmrac/mqr@latest
```

## :zap: Usage

As CLI

```shell
mqr "Hello Gopher!"
```

> Will print QR code containing data `Hello Gopher!`

As TUI

```shell
mqr
```

> Will launch TUI where you can enter (or paste) your data

Keybindings

- `Up` / `Down` - move focus
- `Space` - select setting
- `Esq` - clear input or exit
- `Enter` - print result and exit

## :lock: License

mQR is distributed under [MIT license](LICENSE)
