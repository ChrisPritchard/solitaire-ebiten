# Technical design

There are two overall approaches to card management:

- every card has an x,y and maybe z, and are dumbly rendered accordingly. game rules with stacks, tableaus etc is not represented in the state directly but instead in allowable actions and consequences. some additional state data might be required, like origin for moved cards
- cards are not objects at all, but instead each section of the board is represented as an array or nilable property. the only dynamic thing is anything being moved, or floated which might require additional data be kept like origin.

i think cards will either need to track what is stacked on them, or perhaps moveable stacks will need to exist as an object. this is to facilitate not just moving the bottom or free card of any stack, but in the case of the tableau piles, the ability to move the top most card where all stacked cards follow stacking rules.

to make the project more useful for other projects, i think embedding the rules directly into the structure is NOT appropriate. instead:

- track all cards by position and z-order
- each card also has a content which could be a go interface, with a draw function
- for solitaire, the content contains suit, value and whether its face up / shown. draw can then select the appropriate image

rules can be described as an object with a few methods maybe. board setup, can drag, drop action, has won etc. maybe a mvu system where events include game init, drag start, stop, click etc.

- init setups up the board, perhaps by first creating the stock and then emitting a series of move events
- when player clicks on a card, if the card can be moved its marked as moving in internal game state
- as the player drags, the moving card and all valid stacked are also moved
- if the player releases, check if it can be stacked in place. if not issue move orders back to position
- check victory conditions as appropriate

## dragging/mouse clicks

I want to support both mouse and touch inputs, which work differently in ebiten. To do this I need to track:

- whether these are pressed, or released
- the current state of a drag action
- if just pressed, is there a card that can be moved? if so, memorise it
- on updates, if the drag position has moved, then move the moving card
- on release, determine if the card can be placed there, if so place it, else snap back to its original position

suggestion: start with mouse and abstract to just whats needed: position, state (justpressed, dragging, released)
pass info to game rules, which can manage if something should occur

## initial setup

- shuffle deck
- render in position as a stack
- render spaces for foundation, waste etc
- deal out piles
  - move action moving one card at a time...
  -

## coordinate system

cards can only be in certain positions, when placed. they can be anyway while being dragged, but the placement positions on the board are fixed:

```
[]  []
[]  [] [] []
[]     [] []
[]        []
```

etc.

establishing this 'grid' would likely be a help

the grid is basically defined in card 'units', with whole units for things like the foundations, and partial units for the piles which overlay their contents

rather than using pixel coordinates, then, it might be advisable to store card positions in terms of units. however, how does this apply to dragging?

> is the separation between 'game' and the 'rules' getting in the way? or does it help keep things unique. i think i should drop it, for now, while still being mindful of separation. the main difference would be removing 'card content' from cards, and just embedding that data straight. also not relying entirely on a separate rules object for where things are placed. that being said, i like the moving of pure game stuff, like reading ebitengine's touch state, into its own thing. perhaps those could be removed into a helper class instead?
> there is also the separation of game objects from their image assets, which allows the 'game' to be more portable (if i do what i originally intended to do, and also build this in godot)
> and the fact that certain rendering considerations, like resolution, don't need to be a factor in the game. as a counter though, pixel positions are a mix, and sounds will also be a mix (though that might be solved with some sort of mvu eventing system)

a few modifications both to help with these grid mechanics, and dragging.

the idea is that a card tracks its board 'position', defined in card units. this could be, in the y direction, four units per card while in the x direction just two. the renderer part of things would handle converting this to pixel units
dragging becomes a little complex - while the game stores the drag state, and this could be modified to not use pixel coordinates on the card, how would that be surfaced to the renderer? presently cards are provided as a slice, so it would need to be (in the current structure) part of the card slice elements.

couple of options here; actually they're sort of the same option: use a different model for the presentation to the renderer. this would just be pixel positions and image asset - almost like a view model. this does take asset and size conversion out of the renderer and puts them in this view model transform, wherever that is, but that could be fine - it further simplifies the core game loop object.

## MVVM in go

how would this look?

we have the game mechanics, in one file
we have a view transform, that could just be a function
we have the renderer interfaces with ebitengine, including picking up the touch state on update

maybe keep Main, holding a GameLoop (instead of just Game, making somewhat clear it isnt the whole thing)
on update, the GameLoop would call a Transform function, passing in the Sawayama model object. This would no longer be a 'RuleSet' interface, maybe, but just a model or struct.

> in this way, the transform function doesnt get silly with its signature. it can be defined directly on sawayama. the close coupling might suggest they could be one object, but this distinctness allows for sawayama to be taken out and converted into another engine like the original concept.

the transform returns images, and maybe sounds? which the gameloop then draws.
