annyong
=======

![Hello](annyong.jpg)

annyong is a program which annoys a given cmd with random interrupts

it can be useful to test services which need graceful interrupt handling

usage
=====

    annyong <flags> <cmd> <args>

will run the `cmd` with `args` n times in parallel and interrupt each one after
a random time

to control the times and duration before interrupt, use the following flags

flags
=====

    -times=10
    -minwait=1
    -maxwait=30
