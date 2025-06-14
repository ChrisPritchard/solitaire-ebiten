# Solitaire in Ebiten

Using [Sawayama rules](https://www.watsonbrosgames.com/solitaire/) by [Zachtronics](https://www.zachtronics.com/solitaire-collection/)

## Asset credits

Cards: [Pixel Art Cards by Glenn Dittman](https://opengameart.org/content/pixel-art-cards)
Sounds: [54 Casino sound effects (cards, dice, chips) by Kenney](https://opengameart.org/content/54-casino-sound-effects-cards-dice-chips)
Backgrounds: [Felt Backgrounds](https://opengameart.org/content/felt-backgrounds)

The game cards are provided as a folder full of aseprite files. I have aseprite on steam for windows, but also WSL (which can run windows binaries from linux), so was able to convert them all to png with the following bash:

```bash
alias aseprite='/mnt/c/Program\ Files\ \(x86\)/Steam/steamapps/common/Aseprite/Aseprite.exe'
for file in *.aseprite; do
    aseprite -b "$file" --save-as "${file%.*}.png"
done
```

Though ultimately only 'card-suites.png' was needed.
