annyong
=======

![Hello](annyong.jpg)

annyong is a program which annoys a given cmd with random interrupts

it can be useful to test services which need graceful interrupt handling

usage
=====

    annyong <flags> <cmd> <args>

will run the `cmd` with `args` n times and interrupt each one after
a random time

to control the times and duration before interrupt, use the following flags (defaults indicated)

flags
=====

    -times=10 // number of times to re-launch the cmd
    -minwait=1 // minimum seconds after before interrupting
    -maxwait=10 // maximum seconds after before interrupting
    -parallel=false // when true runs the cmd in parallel using goroutines
    -quit=true // when true stop after receiving the first non-zero return code (unused if -parallel=true)
