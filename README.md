# GreenScout-CLI
The Command-line interface for FRC team 1816's scouting app, GreenScout!

# Getting Started

It's expected that you at least know basic programming principles and practices. Additionally, knowing the [Go Programming Language](https://go.dev/learn/) will be neccessary when contributing to this CLI.

To get started, you'll first need the [Go Programming Language](https://go.dev/dl/), [Git](https://git-scm.com/downloads), and [VS Code (Optional)](https://code.visualstudio.com/Download).

Once you have all that installed, open up your terminal and enter this command
```bash
git clone https://github.com/TheGreenMachine/GreenScout-CLI.git
```

This will download the repository onto your computer and to move into it type this
```bash
cd GreenScout-CLI
```

Then, download all of the dependencies of this project with
```bash
go get
```

Then, compile the CLI by running
```bash
go build
```

Finally, to see the CLI's options, type
```bash
./GreenScoutCLI --help
```

If you are so inclined, you can move the executable into your $PATH, though I do not recommend this, as the CLI is under constant development and constantly replacing such an executable may get tiresome. 

If you are using VS Code, I highly recommend installing the [Golang Extension](https://marketplace.visualstudio.com/items?itemName=golang.Go) and [Dart Extension]as they provide code highlighting, code suggestions, an integrated build system, and debug console into the editor.

Additional documentation, such as a list of commands available can be found [here](./docs/).

# Roadmap (Tasks for future devs)
- Add proper error handling

# Contributers

- [Tag C](https://github.com/TagCiccone)