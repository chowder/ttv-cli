# `ttv-cli`

A collection of programs to help the average chatter use Twitch from the command line

## `ttv-live`

Checks whether your favourite Twitch streamers are currently live. 

<img src="https://user-images.githubusercontent.com/16789070/182046526-6cc16a6d-32e1-4902-aef9-af8a7905b72b.png" width="600" />

### Configuration

Populate your `~/.config/ttv-cli/ttv-cli.config` file with the list of streamers to track.

```json
{
  "streamers": [
    "xqc",
    "skillspecs",
    "dismellion",
    "dino_xx",
    "usteppin"
  ]
}
```

## `ttv-points`

View and redeem Twitch channel point rewards from the command line. 

[ttv-rewards.webm](https://user-images.githubusercontent.com/16789070/184022312-a947118a-c777-4fea-b71c-7485efded5b8.webm)

Also displays countdowns for rewards with cooldowns - common for games with Twitch integration. 
