# Screen Manager

So this project was made because I use a VPS that cannot launch docker containers.
Meaning that it was pretty annoying to launch multiple apps in the background, playing with **screen** has become my 
main way of dealing with that.

BUT, playing those scripts everytime my computer started, checking which screens were up
or not if they had crashed was really annoying too.

That's why I made this manager.

## Features

- Starting a screen
- Attach to a screen to monitor it
- Edit a screen config in app (While not running)
- Stop a screen

### Upcoming
- Launching a group of screens

## Usage

To use that app, you need to create a YAML file for every screen you want to run.

The yaml files follow this template : 
```yaml
name: string  # Name of the container
id: string    # Name of the screen session
description: string   # Description of the container
user: string  # User that can run that script
commands:     # Commands that the screen will play in order
  - string
  - string
  - ...
```

You can see an example file in ``config/``.

## Dependencies
- Go: 19^
- GNU/Linux
- Screen

## Special Thanks
I need to shout out a special thanks to the [Charm.sh](https://charm.sh)
community for the amazing terminal UI/UX packages they have made.

## Licence
MIT.
