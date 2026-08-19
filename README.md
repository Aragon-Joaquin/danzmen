> [!WARNING]  
> Pretty much on beta. not even for a v0.0.1 and reiterating over and over with api changes

## danzmen
more like a tui agenda using a .toml config file.


[example_image](./.github/example.png)
(the status logo comes from [ufetch](https://github.com/jschx/ufetch/blob/master/ufetch-void))

the idea is:
- when you open a terminal instead of opening your `riced fastfest with a pokemon`, or whatever bloat you use, it reminds you what monthly/longterm tasks you need to do.
- the better solution would be using a list.todo.md with [checkmate.nvim](https://github.com/bngarren/checkmate.nvim) (not sponsored btw) & zk.org/neorg/obsidian but whatever fits your bloat
- yes, im using this cuz pen and paper is too hard

```toml

[longterm]
tasks = [
  {ends = "12/30/2027", name = "Get marry idk", priority = "low"}, # date format is MM/DD/YYYY
  {ends = "500d", name = "Get one follower in github lol", priority = "high"} # you can also do this but it gets registered once you run the program. its also kinda buggy
]

[month.july]
tasks = [
    {name = "DONT bath"}, #one time only!
    {name = "watch 300 hours of vtubers", times = 300 } # keep track of your progress being 300 the objective!
]

# every month it gets refreshed
[month.every]
tasks = [ 
        "talk to a wom*n", # you can combine both!!! same as {name = ""}
        {name = "farm *any gacha currency*", times = 1337.55},
        {name = "consume 120 petabytes of reels", times = 120, metric = "PB"} # and specify the metric of what we're talking about! (it just to be reminded of, provides NO functionality)
    ]
```

### Installation (+ more commands lol)!

**Requirements:**
- GO `1.26.5 linux/amd64` (only tested in this version, expect to work in multiples)
- SQLite3
- Have the $HOME env variable set lol (like every distro)

```sh
sudo make install
echo 'eval "$(danzmen list)"' >> ${HOME}/.bashrc #it executes each time a shell opens!

# you can skip these:
make config # setups + makes the default config file (skips if already one exists)
make symlink # makes a symbolic link to your desktop so its easy to change this constantly
```

And more useful things:
```sh
make del-db # deletes the db in case of an error
sudo make clean # removes the program from /usr/local
```
