# Simple IRC Server

A limited implementation of the IRC (Internet Relay Chat) server protocol written in Go. The implementation follows specifications laid out in the RFC documents.

## Features

This is an implementation of a standalone IRC server. It doesn't support Server-Server communication but provides the core functionality needed for clients to connect and communicate.

### Supported Commands

The server implements a limited set of commands:

- `NICK`: Set or change nickname
- `USER`: Set username, mode, and realname
- `JOIN`: Join one or more channels
- `PRIVMSG`: Send messages to channels or users
- `PART`: Leave one or more channels
- `QUIT`: Disconnect from the server
- `WHO`: Query information about users

### RFC Compliance

This implementation follows the RFC specifications for IRC, particularly:
- RFC 1459: Internet Relay Chat Protocol
- RFC 2810: Internet Relay Chat: Architecture
- RFC 2811: Internet Relay Chat: Channel Management
- RFC 2812: Internet Relay Chat: Client Protocol

## Usage

### Building and Running

```bash
go build
./main
```

By default, the server listens on port 6667. You can connect to it using any IRC client.

### Testing

This server has been tested using a GUI client to ensure compatibility and proper functioning.

## Implementation Notes

This is my first semi-serious project in Golang, so the architecture choices reflect my learning process.

I opted for a "fat client" design where most of the logic is contained within the Client struct. This decision was made because splitting components into their own packages introduced cyclical dependencies. Solving that would have required introducing abstractions that don't really address any domain problems.

The main components include:

- `Server`: Main coordinator that handles connections and maintains state
- `Client`: Represents a connected user and handles command processing
- `Chatroom`: Manages channel membership and messaging
- `Message`: Parses and represents IRC protocol messages

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

```
MIT License

Copyright (c) 2025 Ged Dackys

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```
