# gomali

Just a simple, fast and pretty biased **Ma**rkdown **li**nter written in **Go**.

It's mainly meant for integration into CIs and editors.

#### Installation

```sh
$ go get -u github.com/mhinz/gomali
```

#### Usage

```sh
$ gomali <files>
```

#### Vim

Put this in your vimrc:

```vim
autocmd FileType markdown setlocal makeprg=gomali\ %
```

Afterwards, in a Markdown file, use `:make` to populate the [quickfix
list](https://github.com/mhinz/vim-galore#quickfix-and-location-lists) with all
formatting issues, if there are any. Have a look at them with `:copen`.
