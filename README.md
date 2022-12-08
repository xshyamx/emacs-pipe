# Emacs Pipe #

Solves the problem of sending text from stdin to emacs. `emacsclient`
can handle files but, not `stdin`. The common solution is to redirect to a temporary file and open it via emacsclient (See [Piping things into Emacs][emacs-wiki]).

The problem with a simple script like [`to-emacs`](to-emacs) is when tailing something. For example the following command will block till the container is stopped to show the logs. Another disadvantage is that since it is opening a file emacs will immediately switch to that buffer interrupting whatever you where doing.

```sh
docker logs -f <cid> 2>&1 | to-emacs
```

Alternatively, we can leverage the eval argument of emacsclient to create a buffer in emacs without creating a temporary file (See [Reading stdin with Emacs Client][emacs-stdin])

Taking the above idea this program attempts to stream the stdin as it becomes available to a temporary emacs buffer called `*stdin*` (and `*stdin<n>*` when there are more than one). The program exits if either the `stdin` is read completely or the `*stdin*` buffer is closed.

## Build & Run ##

Build and move to somewhere in the path

```sh
# build executable
make build
# copy to path
cp build/ep $HOME/.local/bin
```

Start tailing to emacs

``` sh
./simple.sh | ep &
```

Alternatively, you can specify the mode you want eg.

``` sh
kubectl get deploy/nginx -o yaml | ep yaml
```

### Eshell ###

Another alternative is to use [Eshell][eshell] where you can redirect to emacs buffers directly eg.

```sh
jq '.items[0]' pods.json > #<first.json>
```

While this works well since Eshell is not a terminal several commands have problems notably git & pagers

[emacs-wiki]: https://www.emacswiki.org/emacs/EmacsPipe
[emacs-stdin]: https://mina86.com/2021/emacs-stdin/
[eshell]: https://www.emacswiki.org/emacs/CategoryEshell
