<p align="center">
  <img src=".github/assets/maersk.jpeg" width="500" alt="logo" /><br />
</p>

<br />

<p align="center">
    <img src="https://img.shields.io/badge/Go-1.19-00ADD8?style=for-the-badge&logo=go" alt="go version" />
    <img src="https://img.shields.io/badge/Version-0.1.1-00ADD8?style=for-the-badge&logo=github" alt="version" /><br />
    Concurrent and Secure file downloader implemented in Golang
</p>

## Table of contents

- [How to use Maersk](#how-to-use)
- [Params](#params)
- [Example](#example)
- [Idea](#idea)
- [Contribute](#contribute)

## What is Maersk?

Inspired from [Maersk](https://www.maersk.com/) shipping company, I implemented a file downloader that downloads files
in concurrent. Golang workers allow us to have concurrency when we are downloading a file. Also this can helps us with
error handling and storing the downloaded file parts, so we don't have to download the whole file if an error occurs.

## About Maersk

**Maersk**, is a Danish shipping company, active in ocean and inland freight transportation and associated services, 
such as supply chain management and port operation. 
Maersk was the largest container shipping line and vessel operator in the world from 1996 until 2022.

## How to use?

```shell
git clone https://github.com/amirhnajafiz/maersk.git
cd maersk
make build
```

In order to use ```maersk``` in every place on your system, make sure to the followings to either ~/.zshrc, ~/.bash_profile, or ~/.bashrc.

```shell
export PATH="<path-to-cloned-repository>:$PATH"
```

## Params

Parameters of ```maersk``` struct are as follows:
 
|  Field  | Description                                                     |  Value   | Example                         |
|:-------:|-----------------------------------------------------------------|:--------:|---------------------------------|
| output  | Output file name to store the downloaded information in it      |  string  | ```"file.zip"```                |
|   url   | Address of the file that you want to download (http, https url) |  string  | ```"example.com/file.tar.gz"``` |
| workers | The number of workers to download the file in concurrent        |   int    | ```5```                         |
| chunks  | The max number of chunks to download the file in concurrent     |   int    | ```20```                        |
| timeout | Timeout for downloading each chunck of file from server         | duration | ```10 * time.Second```          |
|  mode   | Set the error modes of cargo (debug or info or off)             |  string  | ```DEBUG, INFO, OFF```          |

## Example

```shell
maersk -output "5MB.zip" -url "http://212.183.159.230/5MB.zip" -workers 10
```

## Idea

The idea behind this project is inspired from [Cheikh Seck](https://blog.devgenius.io/concurrent-file-download-with-go-495d7b946492) 
medium story about concurrent file download with Go. Special thanks to **Cheikh** for giving me this idea to build **Maersk**.

## Contribute

Feel free to submit **Issues** about project or give your ideas.
