# Solitaire in Ebiten

Using [Sawayama rules](https://www.watsonbrosgames.com/solitaire/) by [Zachtronics](https://www.zachtronics.com/solitaire-collection/)

Posts about this on my blog, [grislygrotto.nz](https://grislygrotto.nz):

- [Artisanal, 28th June 2025](https://grislygrotto.nz/post/artisanal)
- [Messy Conclusions, 31st June 2025](https://grislygrotto.nz/post/messy-conclusions)

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

## Learnings (for me)

Harder than expected, due to a 20/80, 80/20 sort of thing; getting the architecture right took most of the work. Progress on this can be seen in my stream of consciousness [design.md](./design.md)

No AI was used for code generation, though I did use a little DeepSeek to figure out perculiarities in Go. Also in what might have been the most frustrating part, where I had a card-duplication bug that took me *hours* and *hours* to track down, as it wasn't consisted and required just continuously playing until it happened. In frustration I finally submitted the chunk of code where I had identified it happening to DeepSeek for help, and it told me about slices being reference types which was the problem :) Don't know why that was so hard for me, but thanks Chinese Supermind!
