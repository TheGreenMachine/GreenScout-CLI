# Commands
### A Guide to all the included commands

# General Utilities

## login
- ### Alias: L
- ### Usage: Tries to log in to the app, storing credentials locally if successful
- ### Flags
    - username
        - Alias: u
    - password
        - Alias: p
- ### Example:
```bash
./GreenScoutCLI login --username "user" --password "password"
```

## update-address
- ### Usage: Sets/updates the web address the CLI will attempt to use as the server
- ### Flags
    - address
        - Alias: a
- ### Example:
```bash
./GreenScoutCLI update-adress --address "example.com"
```

## getAddress
- ### Usage: Gets the configured server address
- ### Example:
```bash
./GreenScoutCLI getAddress
```


## validate
- ### Alias: v
- ### Usage: Sends a request to the configured server address to confirm if it is active
- ### Example:
```bash
./GreenScoutCLI validate
```

## getSchedule
- ### Usage: Gets the schedule.json file representing the match schedule of the current event from the server
- ### Example:
```bash
./GreenScoutCLI getSchedule
```

## getLeaderboard
- ### Usage: Gets the current leaderboard from the server
- ### Example:
```bash
./GreenScoutCLI getLeaderboard
```
- ### Note: The leaderboard's scope has expanded greatly since this command was created. Future devs may want to add options to sort through the returned JSON.

## getUsers
- ### Usage: Gets the list of all users from the server
- ### Example:
```bash
./GreenScoutCLI getUsers
```
- ### Note: Users now have many more properties than when this command was created. Future devs may want to add options to sort through the returned JSON.

## genPassword
- ### Usage: Generates a password hashed through the bcrypt algorithm for insertion into the backend
- ### Flags
    - password
        -   Alias: p
- ### Example:
```bash
./GreenScoutCLI genPassword --password "example"
```

# Server configuration tools

## setKey
- ### Alias: sk
- ### Usage: Sends a request to change the Blue Alliance event key currently configured on the server
- ### Flags
    - Key
        - Alias: k
- ### Example:
```bash
./GreenScoutCLI setKey -k "2024mnst" // The 2024 Minnesota State Championship
```

## update-sheet
- ### Usage: updates the google sheets spreadsheet ID the server will attempt to use 
- ### Flags
    - sheet
        - Alias: s
- ### Example:
```bash
./GreenScoutCLI update-sheet --sheet "15OFX-FCFd2GKtGbzR8ozeuFGOkAetrHXBWyHiQvAW4E" // The spreadsheet ID of the 2024 Granite city spreadsheet
```
- ### How to Get:
    - The sheet id is between the /d/ and the /edit in the URL of a google sheets link
    - docs.google.com/spreadsheets/d/**SheetID**/edit

# Admin Tools

## getScouterSchedule
- ### Usage: Gets the assigned matches of a given scouter from the server
- ### Flags:
    - scouter
- ### Example:
```bash
./GreenScoutCLI getScouterSchedule --scouter "somescouter"
```

## modify-leaderboard
- ### Usage: A tool for modifying user scores
- ### Flags:
    - name
        - Alias: n
    - Modification
        - Alias: m
        - Options:
            - Increase
            - Set
            - Decrease
    - By
        - Alias: b
- ### Example:
```bash
./GreenScoutCLI modify-leaderboard -n "example" -m "Increase" -b 1
```

## addBadge
- ### Usage: Adds a badge to a user on the server
- ### Flags:
    - name
        - Alias: n
    - badge
        - Alias: b
    - description
        - Alias: d
- ### Example:
```bash
./GreenScoutCLI addBadge -n "example" -b "Example badge" -d "An example description"
```
