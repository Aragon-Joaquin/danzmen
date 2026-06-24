more like a tui daily agenda using a .toml config file.

the idea is:
- when you open a terminal instead of opening your `riced fastfest with a pokemon`, or whatever bloat you use, it reminds you what daily tasks you need to do.
- after few days you probably wont even care or remove this, but at that point i expect you to KNOW your routine, making this program accomplish its mission.
- yes, im using this cuz pen and paper is too hard

```toml
schedule.start = "sunday"

# each Monday
[monday]
tasks = [ "run 500km", "read the man pages" ]

# starting from Monday, every 2 days until Friday
[loop.2d]
start = "monday"
end = "friday"
tasks = ["goto gym"] # Monday, Wednesday, Friday

[loop.1w] # 1 week
tasks = ["shower"] # every week, starting from Sunday


#ERROR: cannot loop if the loop is greater than the difference between days/weeks
# from Monday to Tuesday there's a 1 day and its declare to repeat each 4 days
[loop.4d]
start = "monday"
tasks = ["this an invalid format"]
end = "tuesday" 
```

for now, the .config file is subject to changes since im horrible at naming things.

## todo: 
- [ ] sqlite
- [ ] everything else
