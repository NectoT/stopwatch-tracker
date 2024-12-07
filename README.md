# An app to track your activities

*Currently this is just a proof of concept*

This is an app for those who were wondering where did all of their time go, but did not want to 
install apps that check your every action on the PC. Create a stopwatch for every 
time-consuming activity and start it whenever you start doing that activity, and then stop it afterwards. There is also a simple stats page

## Tech stack

This app was made using *Wails* (v3 alpha) and *Svelte*. For the charts *chart.js* is used.

## Building it yourself

Since this is a wails app, just follow the [instructions](https://v3alpha.wails.io/getting-started/installation/) for its installation. At the moment of this being written, `wails3 build` command isn't working, so you'll have to also download [Task](https://taskfile.dev/), which Wails v3 is using, and use `task build` command instead.