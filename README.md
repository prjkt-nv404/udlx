# UDLX - Yet Another Udemy Course Downloader!
> *A cross-platform command-line based Udemy course video downloader.*
----

<p align="left">
  <a href="" rel="noopener">
 <img width=550px height=550px src="https://raw.githubusercontent.com/prjkt-nv404/udlx/main/preview/udlx-ss.png?raw=true?raw=true" alt="logo"></a>
</p>

### 🔥 Download from releases
- **[Click here ⏬](https://github.com/prjkt-nv404/udlx/releases)**
## Note: This Project created just for Education purpose.
### 🔥 Features
- _`Choose video quality.`_
- _`Set Download Start and Download End.`_
 
**❗By default the courses will be downloaded to the program's main folder. The structure of course content will be preserved.**
### ❔ Todo
- _`Proper download directory.`_
- _`Subtitles download`_

### 💻 Platform 
- _`Linux`_
- _`Mac`_
- _`Windows`_
----
### Build:
  - use Go 1.19.3 _(with prebuilt `go.mod` & `go.sum`)_
   ```
   go build -o ./bin/udlx-debug udlx.go

   ./bin/udlx-debug
   ```
---
### How to use
  - You should have udemy account and purchased course.
  - Login to udemy -> Right click -> Go to Inspect Element
  - Go to Aplication -> Cookies -> Type access_token in filter, then copy it.

  - Now download the Executables binary from _`bin`_ folder according your platform
  - After download execute(double click) the binary file
  - You can also run through terminal _`./<executable-file-name>`_
  - In case permission got change then apply the executable permission _`chmod 755 <executable-file-name>`_
  - Mac doesn't allow to execute binary file directly, so you first need to allow it from security.
---
### ⚠️ Disclaimer:

This software is intended to help you download Udemy courses for personal use only. Sharing the content of your subscribed courses is strictly prohibited under Udemy Terms of Use. Each and every course on Udemy is subjected to copyright infringement.

This software does not magically download any paid course available on Udemy, you need to provide your Udemy `access_token` to download the courses you have enrolled in.

---
**❗No part of this tool is owned by me!**

**All credit to [maabiddevra](https://github.com/maabiddevra/udemy-dl.git)**

_**with love ❤️ from space!**_
