![markclient](https://github.com/markusnetwork/.github/assets/96031819/fe5d9023-d98d-4fc3-8038-f93d54a72c69)
![donuts-are-good's followers](https://img.shields.io/github/followers/donuts-are-good?&color=555&style=for-the-badge&label=followers) ![donuts-are-good's stars](https://img.shields.io/github/stars/donuts-are-good?affiliations=OWNER%2CCOLLABORATOR&color=555&style=for-the-badge) ![donuts-are-good's visitors](https://komarev.com/ghpvc/?username=donuts-are-good&color=555555&style=for-the-badge&label=visitors)

# MarkClient

MarkClient is a basic terminal-based client for browsing the Markus document network. It retrieves and displays Markdown documents served by a compatible server, like MarkServ.

## Installation

To install MarkClient, clone the repository and build the project:

```shell
git clone https://github.com/donuts-are-good/markclient.git
cd markclient
go build
```
## Usage

MarkClient is a command-line tool and accepts a URL as its single argument:

```shell
./markclient http://localhost:88/myfile.md
```

The client will connect to the specified URL, retrieve the Markdown document, and display it in the terminal. Text formatting (such as bold, italics, and underlined text) is preserved using ANSI escape codes.

If the document contains links, they are numbered and listed at the bottom of the document. To follow a link, type its number and press Enter. To quit the client, type q and press Enter.

## Offline Mode

In addition to displaying documents, MarkClient also saves them for offline reading. Each document is saved in the offline directory, under a subdirectory that matches the domain name from the URL. For example, a document from http://localhost:88/myfile.md would be saved as `offline/localhost/myfile.md`.

The offline mode allows you to access previously viewed documents even when you're not connected to the network.

## Compatibility

MarkClient is designed to work with any server that serves Markdown documents over HTTP, including MarkServ. However, it can also display any Markdown file that is accessible via an HTTP GET request.
Contribution

Contributions are always appreciated. If you see something, say something, and if you're motivated, please open an issue or a pull request on the GitHub repository.

Enjoy your journey through the Markus document network!

## license
MIT License 2023 donuts-are-good, for more info see license.md
