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
