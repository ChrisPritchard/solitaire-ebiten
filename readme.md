# Solitaire in Ebiten

Using [Sawayama rules](https://www.watsonbrosgames.com/solitaire/) by [Zachtronics](https://www.zachtronics.com/solitaire-collection/)

## Rule Mechanics

Four areas:

- the 'tableau' which is made up of seven piles, all face up, each with one more card than the previous
- the 'stock', consisting of all remaining cards from which three cards can be drawn at a time into...
- the 'waste', where a new horizontal pile is made from three-card-at-a-time draws from the stock
- four 'foundations', non-staggered stacks built from Aces up (Ace, 2, 3 etc through to King)

The goal is to build up and complete all four foundations. Cards can be moved from anywhere except the foundations (which can only be moved to) and placed on a 1-higher card with an opposing colour, or in a free spot. A stack of cards can be moved by selecting the highest card that follows these rules for all lower cards. Finally, when the stock has been exhausted (all its cards used or in the waste) its spot can be used as a free cell for a single card.

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
